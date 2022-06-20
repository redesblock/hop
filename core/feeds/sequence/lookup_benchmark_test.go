package sequence_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/feeds"
	"github.com/redesblock/hop/core/feeds/sequence"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
)

type timeout struct {
	storage.Storer
}

var searchTimeout = 30 * time.Millisecond

// Get overrides the mock storer and introduces latency
func (t *timeout) Get(ctx context.Context, mode storage.ModeGet, addr swarm.Address) (swarm.Chunk, error) {
	ch, err := t.Storer.Get(ctx, mode, addr)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			time.Sleep(searchTimeout)
		}
		return ch, err
	}
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return ch, nil
}

func BenchmarkFinder(b *testing.B) {
	for _, prefill := range []int64{1, 100, 1000, 5000} {
		storer := &timeout{mock.NewStorer()}
		topicStr := "testtopic"
		topic, err := crypto.LegacyKeccak256([]byte(topicStr))
		if err != nil {
			b.Fatal(err)
		}

		pk, _ := crypto.GenerateSecp256k1Key()
		signer := crypto.NewDefaultSigner(pk)

		updater, err := sequence.NewUpdater(storer, signer, topic)
		if err != nil {
			b.Fatal(err)
		}
		payload := []byte("payload")

		ctx := context.Background()

		for at := int64(0); at < prefill; at++ {
			err = updater.Update(ctx, at, payload)
			if err != nil {
				b.Fatal(err)
			}
		}
		latest := prefill
		err = updater.Update(ctx, latest, payload)
		if err != nil {
			b.Fatal(err)
		}
		now := prefill
		for k, finder := range []feeds.Lookup{
			sequence.NewFinder(storer, updater.Feed()),
			sequence.NewAsyncFinder(storer, updater.Feed()),
		} {
			names := []string{"sync", "async"}
			b.Run(fmt.Sprintf("%s:prefill=%d, latest/now=%d", names[k], prefill, now), func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					_, _, _, err := finder.At(ctx, now, 0)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
