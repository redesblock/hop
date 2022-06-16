package localstore

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/redesblock/hop/core/logging"
	statestore "github.com/redesblock/hop/core/statestore/mock"

	"github.com/redesblock/hop/core/shed"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/tags"
	tagtesting "github.com/redesblock/hop/core/tags/testing"
	"github.com/syndtr/goleveldb/leveldb"
)

// TestModeSetAccess validates ModeSetAccess index values on the provided DB.
func TestModeSetAccess(t *testing.T) {
	for _, tc := range multiChunkTestCases {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t, nil)

			chunks := generateTestRandomChunks(tc.count)

			wantTimestamp := time.Now().UTC().UnixNano()
			defer setNow(func() (t int64) {
				return wantTimestamp
			})()

			err := db.Set(context.Background(), storage.ModeSetAccess, chunkAddresses(chunks)...)
			if err != nil {
				t.Fatal(err)
			}

			binIDs := make(map[uint8]uint64)

			for _, ch := range chunks {
				po := db.po(ch.Address())
				binIDs[po]++

				newPullIndexTest(db, ch, binIDs[po], nil)(t)
				newGCIndexTest(db, ch, wantTimestamp, wantTimestamp, binIDs[po], nil)(t)
			}

			t.Run("gc index count", newItemsCountTest(db.gcIndex, tc.count))

			t.Run("pull index count", newItemsCountTest(db.pullIndex, tc.count))

			t.Run("gc size", newIndexGCSizeTest(db))
		})
	}
}

// here we try to set a normal tag (that should be handled by pushsync)
// as a result we should expect the tag value to remain in the pull index
// and we expect that the tag should not be incremented by pull sync set
func TestModeSetSyncNormalTag(t *testing.T) {
	mockStatestore := statestore.NewStateStore()
	logger := logging.New(ioutil.Discard, 0)
	db := newTestDB(t, &Options{Tags: tags.NewTags(mockStatestore, logger)})

	tag, err := db.tags.Create("test", 1)
	if err != nil {
		t.Fatal(err)
	}

	ch := generateTestRandomChunk().WithTagID(tag.Uid)
	_, err = db.Put(context.Background(), storage.ModePutUpload, ch)
	if err != nil {
		t.Fatal(err)
	}

	err = tag.Inc(tags.StateStored) // so we don't get an error on tag.Status later on
	if err != nil {
		t.Fatal(err)
	}

	item, err := db.pullIndex.Get(shed.Item{
		Address: ch.Address().Bytes(),
		BinID:   1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if item.Tag != tag.Uid {
		t.Fatalf("unexpected tag id value got %d want %d", item.Tag, tag.Uid)
	}

	err = db.Set(context.Background(), storage.ModeSetSync, ch.Address())
	if err != nil {
		t.Fatal(err)
	}

	item, err = db.pullIndex.Get(shed.Item{
		Address: ch.Address().Bytes(),
		BinID:   1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// expect the same tag Uid because when we set pull sync on a normal tag
	// the tag Uid should remain untouched in pull index
	if item.Tag != tag.Uid {
		t.Fatalf("unexpected tag id value got %d want %d", item.Tag, tag.Uid)
	}

	// 1 stored (because incremented manually in test), 1 sent, 1 synced, 1 total
	tagtesting.CheckTag(t, tag, 0, 1, 0, 1, 1, 1)
}

// TestModeSetRemove validates ModeSetRemove index values on the provided DB.
func TestModeSetRemove(t *testing.T) {
	for _, tc := range multiChunkTestCases {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t, nil)

			chunks := generateTestRandomChunks(tc.count)

			_, err := db.Put(context.Background(), storage.ModePutUpload, chunks...)
			if err != nil {
				t.Fatal(err)
			}

			err = db.Set(context.Background(), storage.ModeSetRemove, chunkAddresses(chunks)...)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("retrieve indexes", func(t *testing.T) {
				for _, ch := range chunks {
					wantErr := leveldb.ErrNotFound
					_, err := db.retrievalDataIndex.Get(addressToItem(ch.Address()))
					if !errors.Is(err, wantErr) {
						t.Errorf("got error %v, want %v", err, wantErr)
					}

					// access index should not be set
					_, err = db.retrievalAccessIndex.Get(addressToItem(ch.Address()))
					if !errors.Is(err, wantErr) {
						t.Errorf("got error %v, want %v", err, wantErr)
					}
				}

				t.Run("retrieve data index count", newItemsCountTest(db.retrievalDataIndex, 0))

				t.Run("retrieve access index count", newItemsCountTest(db.retrievalAccessIndex, 0))
			})

			for _, ch := range chunks {
				newPullIndexTest(db, ch, 0, leveldb.ErrNotFound)(t)
			}

			t.Run("pull index count", newItemsCountTest(db.pullIndex, 0))

			t.Run("gc index count", newItemsCountTest(db.gcIndex, 0))

			t.Run("gc size", newIndexGCSizeTest(db))
		})
	}
}
