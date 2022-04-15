package pushsync_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/redesblock/hop/core/pushsync"
	"github.com/redesblock/hop/core/pushsync/pb"
	"github.com/redesblock/hop/core/topology"
	"github.com/redesblock/hop/core/topology/mock"

	"github.com/redesblock/hop/core/localstore"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/p2p/protobuf"
	"github.com/redesblock/hop/core/p2p/streamtest"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
)

// TestSendToClosest tests that a chunk that is uploaded to localstore is sent to the appropriate peer.
func TestSendToClosest(t *testing.T) {
	logger := logging.New(ioutil.Discard, 0)

	// chunk data to upload
	chunkAddress := swarm.MustParseHexAddress("7000000000000000000000000000000000000000000000000000000000000000")
	chunkData := []byte("1234")

	// create a pivot node and a mocked closest node
	pivotNode := swarm.MustParseHexAddress("0000000000000000000000000000000000000000000000000000000000000000")   // base is 0000
	closestPeer := swarm.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000") // binary 0110 -> po 1

	// Create a mock connectivity between the peers
	mockTopology := mock.NewTopologyDriver(mock.WithClosestPeer(closestPeer))

	storer, err := localstore.New("", pivotNode.Bytes(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	// setup the stream recorder to record stream data
	recorder := streamtest.New(
		streamtest.WithMiddlewares(func(f p2p.HandlerFunc) p2p.HandlerFunc {
			return func(context.Context, p2p.Peer, p2p.Stream) error {
				// dont call any handlers
				return nil
			}
		}),
	)

	// instantiate a pushsync instance
	ps := pushsync.New(pushsync.Options{
		Streamer:      recorder,
		Logger:        logger,
		ClosestPeerer: mockTopology,
		Storer:        storer,
	})
	defer ps.Close()
	recorder.SetProtocols(ps.Protocol())

	// upload the chunk to the pivot node
	_, err = storer.Put(context.Background(), storage.ModePutUpload, swarm.NewChunk(chunkAddress, chunkData))
	if err != nil {
		t.Fatal(err)
	}

	records := recorder.WaitRecords(t, closestPeer, pushsync.ProtocolName, pushsync.ProtocolVersion, pushsync.StreamName, 1, 5)
	messages, err := protobuf.ReadMessages(
		bytes.NewReader(records[0].In()),
		func() protobuf.Message { return new(pb.Delivery) },
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) > 1 {
		t.Fatal("too many messages")
	}
	delivery := messages[0].(*pb.Delivery)
	chunk := swarm.NewChunk(swarm.NewAddress(delivery.Address), delivery.Data)

	if !bytes.Equal(chunk.Address().Bytes(), chunkAddress.Bytes()) {
		t.Fatalf("chunk address mismatch")
	}

	if !bytes.Equal(chunk.Data(), chunkData) {
		t.Fatalf("chunk data mismatch")
	}
}

// TestForwardChunk tests that when a closer node exists within the topology, we forward a received
// chunk to it.
func TestForwardChunk(t *testing.T) {
	logger := logging.New(ioutil.Discard, 0)

	// chunk data to upload
	chunkAddress := swarm.MustParseHexAddress("7000000000000000000000000000000000000000000000000000000000000000")
	chunkData := []byte("1234")

	// create a pivot node and a closest mocked closer node address
	pivotNode := swarm.MustParseHexAddress("0000000000000000000000000000000000000000000000000000000000000000")   // pivot is 0000
	closestPeer := swarm.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000") // binary 0110

	// Create a mock connectivity driver
	mockTopology := mock.NewTopologyDriver(mock.WithClosestPeer(closestPeer))
	storer, err := localstore.New("", pivotNode.Bytes(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	targetCalled := false
	var mtx sync.Mutex

	// setup the stream recorder to record stream data
	recorder := streamtest.New(
		streamtest.WithMiddlewares(func(f p2p.HandlerFunc) p2p.HandlerFunc {

			// this is a custom middleware that is needed because of the design of
			// the recorder. since we want to test only one unit, but the first message
			// is supposedly coming from another node, we don't want to execute the handler
			// when the peer address is the peer of `closestPeer`, since this will create an
			// unnecessary entry in the recorder
			return func(ctx context.Context, p p2p.Peer, s p2p.Stream) error {
				if p.Address.Equal(closestPeer) {
					mtx.Lock()
					defer mtx.Unlock()
					if targetCalled {
						t.Fatal("target called more than once")
					}
					targetCalled = true
					return nil
				}
				return f(ctx, p, s)
			}
		}),
	)

	ps := pushsync.New(pushsync.Options{
		Streamer:      recorder,
		Logger:        logger,
		ClosestPeerer: mockTopology,
		Storer:        storer,
	})
	defer ps.Close()

	recorder.SetProtocols(ps.Protocol())

	stream, err := recorder.NewStream(context.Background(), pivotNode, nil, pushsync.ProtocolName, pushsync.ProtocolVersion, pushsync.StreamName)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	w := protobuf.NewWriter(stream)

	// this triggers the handler of the pivot with a delivery stream
	err = w.WriteMsg(&pb.Delivery{
		Address: chunkAddress.Bytes(),
		Data:    chunkData,
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = recorder.WaitRecords(t, closestPeer, pushsync.ProtocolName, pushsync.ProtocolVersion, pushsync.StreamName, 1, 5)
	mtx.Lock()
	defer mtx.Unlock()
	if !targetCalled {
		t.Fatal("target not called")
	}
}

// TestNoForwardChunk tests that the closest node to a chunk doesn't forward it to other nodes.
func TestNoForwardChunk(t *testing.T) {
	logger := logging.New(ioutil.Discard, 0)

	// chunk data to upload
	chunkAddress := swarm.MustParseHexAddress("7000000000000000000000000000000000000000000000000000000000000000") // binary 0111
	chunkData := []byte("1234")

	// create a pivot node and a cluster of nodes
	pivotNode := swarm.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000") // pivot is 0110

	// Create a mock connectivity
	mockTopology := mock.NewTopologyDriver(mock.WithClosestPeerErr(topology.ErrWantSelf))

	storer, err := localstore.New("", pivotNode.Bytes(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	recorder := streamtest.New(
		streamtest.WithMiddlewares(func(f p2p.HandlerFunc) p2p.HandlerFunc {
			return f
		}),
	)

	ps := pushsync.New(pushsync.Options{
		Streamer:      recorder,
		Logger:        logger,
		ClosestPeerer: mockTopology,
		Storer:        storer,
	})
	defer ps.Close()

	recorder.SetProtocols(ps.Protocol())

	stream, err := recorder.NewStream(context.Background(), pivotNode, nil, pushsync.ProtocolName, pushsync.ProtocolVersion, pushsync.StreamName)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	w := protobuf.NewWriter(stream)

	// this triggers the handler of the pivot with a delivery stream
	err = w.WriteMsg(&pb.Delivery{
		Address: chunkAddress.Bytes(),
		Data:    chunkData,
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = recorder.WaitRecords(t, pivotNode, pushsync.ProtocolName, pushsync.ProtocolVersion, pushsync.StreamName, 1, 5)
}
