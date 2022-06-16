package netstore_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"testing"
	"time"

	validatormock "github.com/redesblock/hop/core/content/mock"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/netstore"
	"github.com/redesblock/hop/core/pss"
	"github.com/redesblock/hop/core/sctx"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
)

var chunkData = []byte("mockdata")

// TestNetstoreRetrieval verifies that a chunk is asked from the network whenever
// it is not found locally
func TestNetstoreRetrieval(t *testing.T) {
	retrieve, store, nstore := newRetrievingNetstore(nil)
	addr := swarm.MustParseHexAddress("000001")
	_, err := nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}
	if !retrieve.called {
		t.Fatal("retrieve request not issued")
	}
	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !retrieve.addr.Equal(addr) {
		t.Fatalf("addresses not equal. got %s want %s", retrieve.addr, addr)
	}

	// store should have the chunk now
	d, err := store.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}

	// check that the second call does not result in another retrieve request
	d, err = nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}

}

// TestNetstoreNoRetrieval verifies that a chunk is not requested from the network
// whenever it is found locally.
func TestNetstoreNoRetrieval(t *testing.T) {
	retrieve, store, nstore := newRetrievingNetstore(nil)
	addr := swarm.MustParseHexAddress("000001")

	// store should have the chunk in advance
	_, err := store.Put(context.Background(), storage.ModePutUpload, swarm.NewChunk(addr, chunkData))
	if err != nil {
		t.Fatal(err)
	}

	c, err := nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}
	if retrieve.called {
		t.Fatal("retrieve request issued but shouldn't")
	}
	if retrieve.callCount != 0 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !bytes.Equal(c.Data(), chunkData) {
		t.Fatal("chunk data mismatch")
	}
}

func TestRecovery(t *testing.T) {
	hookWasCalled := make(chan bool, 1)
	rec := &mockRecovery{
		hookC: hookWasCalled,
	}

	retrieve, _, nstore := newRetrievingNetstore(rec)
	addr := swarm.MustParseHexAddress("deadbeef")
	retrieve.failure = true
	ctx := context.Background()
	ctx = sctx.SetTargets(ctx, "be, cd")

	_, err := nstore.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil && !errors.Is(err, netstore.ErrRecoveryAttempt) {
		t.Fatal(err)
	}

	select {
	case <-hookWasCalled:
		break
	case <-time.After(100 * time.Millisecond):
		t.Fatal("recovery hook was not called")
	}
}

func TestInvalidRecoveryFunction(t *testing.T) {
	retrieve, _, nstore := newRetrievingNetstore(nil)
	addr := swarm.MustParseHexAddress("deadbeef")
	retrieve.failure = true
	ctx := context.Background()
	ctx = sctx.SetTargets(ctx, "be, cd")

	_, err := nstore.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil && err.Error() != "chunk not found" {
		t.Fatal(err)
	}
}

// returns a mock retrieval protocol, a mock local storage and a netstore
func newRetrievingNetstore(rec *mockRecovery) (ret *retrievalMock, mockStore, ns storage.Storer) {
	retrieve := &retrievalMock{}
	store := mock.NewStorer()
	logger := logging.New(ioutil.Discard, 0)
	validator := swarm.NewChunkValidator(validatormock.NewValidator(true))

	var nstore storage.Storer
	if rec != nil {
		nstore = netstore.New(store, rec.recovery, retrieve, logger, validator)
	} else {
		nstore = netstore.New(store, nil, retrieve, logger, validator)
	}
	return retrieve, store, nstore
}

type retrievalMock struct {
	called    bool
	callCount int32
	failure   bool
	addr      swarm.Address
}

func (r *retrievalMock) RetrieveChunk(ctx context.Context, addr swarm.Address) (chunk swarm.Chunk, err error) {
	if r.failure {
		return nil, fmt.Errorf("chunk not found")
	}
	r.called = true
	atomic.AddInt32(&r.callCount, 1)
	r.addr = addr
	return swarm.NewChunk(addr, chunkData), nil
}

type mockRecovery struct {
	hookC chan bool
}

// Send mocks the pss Send function
func (mr *mockRecovery) recovery(chunkAddress swarm.Address, targets pss.Targets) error {
	mr.hookC <- true
	return nil
}

func (r *mockRecovery) RetrieveChunk(ctx context.Context, addr swarm.Address) (chunk swarm.Chunk, err error) {
	return nil, fmt.Errorf("chunk not found")
}
