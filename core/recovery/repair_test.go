package recovery_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	accountingmock "github.com/redesblock/hop/core/accounting/mock"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/netstore"
	"github.com/redesblock/hop/core/p2p/streamtest"
	"github.com/redesblock/hop/core/pss"
	"github.com/redesblock/hop/core/pushsync"
	pushsyncmock "github.com/redesblock/hop/core/pushsync/mock"
	"github.com/redesblock/hop/core/recovery"
	"github.com/redesblock/hop/core/retrieval"
	"github.com/redesblock/hop/core/sctx"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/storage/mock"
	storemock "github.com/redesblock/hop/core/storage/mock"
	chunktesting "github.com/redesblock/hop/core/storage/testing"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/topology"
)

// TestRecoveryHook tests that a recovery hook can be created and called.
func TestRecoveryHook(t *testing.T) {
	// test variables needed to be correctly set for any recovery hook to reach the sender func
	chunkAddr := chunktesting.GenerateTestRandomChunk().Address()
	targets := pss.Targets{[]byte{0xED}}

	//setup the sender
	hookWasCalled := make(chan bool, 1) // channel to check if hook is called
	pssSender := &mockPssSender{
		hookC: hookWasCalled,
	}

	// create recovery hook and call it
	recoveryHook := recovery.NewRecoveryHook(pssSender)
	if err := recoveryHook(chunkAddr, targets); err != nil {
		t.Fatal(err)
	}
	select {
	case <-hookWasCalled:
		break
	case <-time.After(100 * time.Millisecond):
		t.Fatal("recovery hook was not called")
	}
}

// RecoveryHookTestCase is a struct used as test cases for the TestRecoveryHookCalls func.
type recoveryHookTestCase struct {
	name           string
	ctx            context.Context
	expectsFailure bool
}

// TestRecoveryHookCalls verifies that recovery hooks are being called as expected when net store attempts to get a chunk.
func TestRecoveryHookCalls(t *testing.T) {
	// generate test chunk, store and publisher
	c := chunktesting.GenerateTestRandomChunk()
	ref := c.Address()
	target := "BE"

	// test cases variables
	targetContext := sctx.SetTargets(context.Background(), target)

	for _, tc := range []recoveryHookTestCase{
		{
			name:           "targets set in context",
			ctx:            targetContext,
			expectsFailure: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hookWasCalled := make(chan bool, 1) // channel to check if hook is called

			// setup the sender
			pssSender := &mockPssSender{
				hookC: hookWasCalled,
			}
			recoverFunc := recovery.NewRecoveryHook(pssSender)
			ns := newTestNetStore(t, recoverFunc)

			// fetch test chunk
			_, err := ns.Get(tc.ctx, storage.ModeGetRequest, ref)
			if err != nil && !errors.Is(err, netstore.ErrRecoveryAttempt) && err.Error() != "error decoding prefix string" {
				t.Fatal(err)
			}

			// checks whether the callback is invoked or the test case times out
			select {
			case <-hookWasCalled:
				if !tc.expectsFailure {
					return
				}
				t.Fatal("recovery hook was unexpectedly called")
			case <-time.After(1000 * time.Millisecond):
				if tc.expectsFailure {
					return
				}
				t.Fatal("recovery hook was not called when expected")
			}
		})
	}
}

// TestNewRepairHandler tests the function of repairing a chunk when a request for chunk repair is received.
func TestNewRepairHandler(t *testing.T) {
	logger := logging.New(ioutil.Discard, 0)

	t.Run("repair-chunk", func(t *testing.T) {
		// generate test chunk, store and publisher
		c1 := chunktesting.GenerateTestRandomChunk()

		// create a mock storer and put a chunk that will be repaired
		mockStorer := storemock.NewStorer()
		defer mockStorer.Close()
		_, err := mockStorer.Put(context.Background(), storage.ModePutRequest, c1)
		if err != nil {
			t.Fatal(err)
		}

		// create a mock pushsync service to push the chunk to its destination
		var receipt *pushsync.Receipt
		pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk swarm.Chunk) (*pushsync.Receipt, error) {
			receipt = &pushsync.Receipt{
				Address: swarm.NewAddress(chunk.Address().Bytes()),
			}
			return receipt, nil
		})

		// create the chunk repair handler
		repairHandler := recovery.NewRepairHandler(mockStorer, logger, pushSyncService)

		// invoke the chunk repair handler
		repairHandler(context.Background(), c1.Address().Bytes())

		// check if receipt is received
		if receipt == nil {
			t.Fatal("receipt not received")
		}

		if !receipt.Address.Equal(c1.Address()) {
			t.Fatalf("invalid address in receipt: expected %s received %s", c1.Address(), receipt.Address)
		}

	})

	t.Run("repair-chunk-not-present", func(t *testing.T) {
		// generate test chunk, store and publisher
		c2 := chunktesting.GenerateTestRandomChunk()

		// create a mock storer
		mockStorer := storemock.NewStorer()
		defer mockStorer.Close()

		// create a mock pushsync service
		pushServiceCalled := false
		pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk swarm.Chunk) (*pushsync.Receipt, error) {
			pushServiceCalled = true
			return nil, nil
		})

		// create the chunk repair handler
		repairHandler := recovery.NewRepairHandler(mockStorer, logger, pushSyncService)

		// invoke the chunk repair handler
		repairHandler(context.Background(), c2.Address().Bytes())

		if pushServiceCalled {
			t.Fatal("push service called even if the chunk is not present")
		}
	})

	t.Run("repair-chunk-closest-peer-not-present", func(t *testing.T) {
		// generate test chunk, store and publisher
		c3 := chunktesting.GenerateTestRandomChunk()

		// create a mock storer
		mockStorer := storemock.NewStorer()
		defer mockStorer.Close()
		_, err := mockStorer.Put(context.Background(), storage.ModePutRequest, c3)
		if err != nil {
			t.Fatal(err)
		}

		// create a mock pushsync service
		var receiptError error
		pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk swarm.Chunk) (*pushsync.Receipt, error) {
			receiptError = errors.New("invalid receipt")
			return nil, receiptError
		})

		// create the chunk repair handler
		repairHandler := recovery.NewRepairHandler(mockStorer, logger, pushSyncService)

		// invoke the chunk repair handler
		repairHandler(context.Background(), c3.Address().Bytes())

		if receiptError == nil {
			t.Fatal("pushsync did not generate a receipt error")
		}
	})
}

// newTestNetStore creates a test store with a set RemoteGet func.
func newTestNetStore(t *testing.T, recoveryFunc recovery.RecoveryHook) storage.Storer {
	t.Helper()
	storer := mock.NewStorer()
	logger := logging.New(ioutil.Discard, 5)

	mockStorer := storemock.NewStorer()
	serverMockAccounting := accountingmock.NewAccounting()
	price := uint64(12345)
	pricerMock := accountingmock.NewPricer(price, price)
	peerID := swarm.MustParseHexAddress("deadbeef")
	ps := mockPeerSuggester{eachPeerRevFunc: func(f topology.EachPeerFunc) error {
		_, _, _ = f(peerID, 0)
		return nil
	}}
	server := retrieval.New(swarm.ZeroAddress, mockStorer, nil, ps, logger, serverMockAccounting, nil, nil, nil)
	recorder := streamtest.NewRecorderDisconnecter(streamtest.New(
		streamtest.WithProtocols(server.Protocol()),
	))
	retrieve := retrieval.New(swarm.ZeroAddress, mockStorer, recorder, ps, logger, serverMockAccounting, pricerMock, nil, nil)
	ns := netstore.New(storer, recoveryFunc, retrieve, logger, nil)
	return ns
}

type mockPeerSuggester struct {
	eachPeerRevFunc func(f topology.EachPeerFunc) error
}

func (s mockPeerSuggester) EachPeer(topology.EachPeerFunc) error {
	return errors.New("not implemented")
}
func (s mockPeerSuggester) EachPeerRev(f topology.EachPeerFunc) error {
	return s.eachPeerRevFunc(f)
}

type mockPssSender struct {
	hookC chan bool
}

// Send mocks the pss Send function
func (mp *mockPssSender) Send(ctx context.Context, topic pss.Topic, payload []byte, recipient *ecdsa.PublicKey, targets pss.Targets) error {
	mp.hookC <- true
	return nil
}
