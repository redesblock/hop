package localstore

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
)

func TestPinning(t *testing.T) {
	chunks := generateTestRandomChunks(21)
	addresses := chunksToSortedStrings(chunks)

	db := newTestDB(t, nil)
	_, err := db.PinnedChunks(context.Background(), 0, 10)

	// error should be nil
	if err != nil {
		t.Fatal(err)
	}

	err = db.Set(context.Background(), storage.ModeSetPin, chunkAddresses(chunks)...)
	if err != nil {
		t.Fatal(err)
	}

	pinnedChunks, err := db.PinnedChunks(context.Background(), 0, 30)
	if err != nil {
		t.Fatal(err)
	}

	if len(pinnedChunks) != len(chunks) {
		t.Fatalf("want %d pins but got %d", len(chunks), len(pinnedChunks))
	}

	// Check if they are sorted
	for i, addr := range pinnedChunks {
		if addresses[i] != addr.Address.String() {
			t.Fatal("error in getting sorted address")
		}
	}
}

func TestPinCounter(t *testing.T) {
	chunk := generateTestRandomChunk()
	db := newTestDB(t, nil)

	// pin once
	err := db.Set(context.Background(), storage.ModeSetPin, swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	pinCounter, err := db.PinCounter(swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if pinCounter != 1 {
		t.Fatalf("want pin counter %d but got %d", 1, pinCounter)
	}

	// pin twice
	err = db.Set(context.Background(), storage.ModeSetPin, swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	pinCounter, err = db.PinCounter(swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if pinCounter != 2 {
		t.Fatalf("want pin counter %d but got %d", 2, pinCounter)
	}

	err = db.Set(context.Background(), storage.ModeSetUnpin, swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.PinCounter(swarm.NewAddress(chunk.Address().Bytes()))
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			t.Fatal(err)
		}
	}
}

func TestPaging(t *testing.T) {
	chunks := generateTestRandomChunks(10)
	addresses := chunksToSortedStrings(chunks)
	db := newTestDB(t, nil)

	// pin once
	err := db.Set(context.Background(), storage.ModeSetPin, chunkAddresses(chunks)...)
	if err != nil {
		t.Fatal(err)
	}

	pinnedChunks, err := db.PinnedChunks(context.Background(), 0, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(pinnedChunks) != 5 {
		t.Fatalf("want %d pins but got %d", 5, len(pinnedChunks))
	}

	// Check if they are sorted
	for i, addr := range pinnedChunks {
		if addresses[i] != addr.Address.String() {
			t.Fatal("error in getting sorted address")
		}
	}
	pinnedChunks, err = db.PinnedChunks(context.Background(), 5, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(pinnedChunks) != 5 {
		t.Fatalf("want %d pins but got %d", 5, len(pinnedChunks))
	}

	// Check if they are sorted
	for i, addr := range pinnedChunks {
		if addresses[5+i] != addr.Address.String() {
			t.Fatal("error in getting sorted address")
		}
	}
}

func chunksToSortedStrings(chunks []swarm.Chunk) []string {
	var addresses []string
	for _, c := range chunks {
		addresses = append(addresses, c.Address().String())
	}
	sort.Strings(addresses)
	return addresses
}
