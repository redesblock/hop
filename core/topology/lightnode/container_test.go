package lightnode_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/topology"
	"github.com/redesblock/hop/core/topology/lightnode"
)

func TestContainer(t *testing.T) {
	t.Run("new container is empty container", func(t *testing.T) {
		c := lightnode.NewContainer()

		var empty topology.BinInfo

		if !reflect.DeepEqual(empty, c.PeerInfo()) {
			t.Errorf("expected %v, got %v", empty, c.PeerInfo())
		}
	})

	t.Run("can add peers to container", func(t *testing.T) {
		c := lightnode.NewContainer()

		c.Connected(context.Background(), p2p.Peer{Address: swarm.NewAddress([]byte("123"))})
		c.Connected(context.Background(), p2p.Peer{Address: swarm.NewAddress([]byte("456"))})

		peerCount := len(c.PeerInfo().ConnectedPeers)

		if peerCount != 2 {
			t.Errorf("expected %d connected peer, got %d", 2, peerCount)
		}
	})
	t.Run("empty container after peer disconnect", func(t *testing.T) {
		c := lightnode.NewContainer()

		peer := p2p.Peer{Address: swarm.NewAddress([]byte("123"))}

		c.Connected(context.Background(), peer)
		c.Disconnected(peer)

		discPeerCount := len(c.PeerInfo().DisconnectedPeers)
		if discPeerCount != 1 {
			t.Errorf("expected %d connected peer, got %d", 1, discPeerCount)
		}

		connPeerCount := len(c.PeerInfo().ConnectedPeers)
		if connPeerCount != 0 {
			t.Errorf("expected %d connected peer, got %d", 0, connPeerCount)
		}
	})
}
