package lightnode

import (
	"context"
	"sync"

	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/topology"
	"github.com/redesblock/hop/core/topology/pslice"
)

type Container struct {
	connectedPeers    *pslice.PSlice
	disconnectedPeers *pslice.PSlice
	peerMu            sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		connectedPeers:    pslice.New(1),
		disconnectedPeers: pslice.New(1),
	}
}

const defaultBin = uint8(0)

func (c *Container) Connected(ctx context.Context, peer p2p.Peer) {
	c.peerMu.Lock()
	defer c.peerMu.Unlock()

	addr := peer.Address
	c.connectedPeers.Add(addr, defaultBin)
	c.disconnectedPeers.Remove(addr, defaultBin)
}

func (c *Container) Disconnected(peer p2p.Peer) {
	c.peerMu.Lock()
	defer c.peerMu.Unlock()

	addr := peer.Address
	if found := c.connectedPeers.Exists(addr); found {
		c.connectedPeers.Remove(addr, defaultBin)
		c.disconnectedPeers.Add(addr, defaultBin)
	}
}

func (c *Container) PeerInfo() topology.BinInfo {
	return topology.BinInfo{
		BinPopulation:     uint(c.connectedPeers.Length()),
		BinConnected:      uint(c.connectedPeers.Length()),
		DisconnectedPeers: toAddrs(c.disconnectedPeers),
		ConnectedPeers:    toAddrs(c.connectedPeers),
	}
}

func toAddrs(s *pslice.PSlice) (addrs []string) {
	_ = s.EachBin(func(addr swarm.Address, po uint8) (bool, bool, error) {
		addrs = append(addrs, addr.String())
		return false, false, nil
	})

	return
}
