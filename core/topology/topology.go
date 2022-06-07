package topology

import (
	"context"
	"errors"
	"io"

	"github.com/redesblock/hop/core/swarm"
)

var (
	ErrNotFound = errors.New("no peer found")
	ErrWantSelf = errors.New("node wants self")
)

type Driver interface {
	PeerAdder
	ClosestPeerer
	EachPeerer
	Notifier
	NeighborhoodDepth() uint8
	SubscribePeersChange() (c <-chan struct{}, unsubscribe func())
	io.Closer
}

type Notifier interface {
	Connecter
	Disconnecter
}

type PeerAdder interface {
	// AddPeers is called when peers are added to the topology backlog
	AddPeers(ctx context.Context, addr ...swarm.Address) error
}

type Connecter interface {
	// Connected is called when a peer dials in, or in case explicit
	// notification to kademlia on dial out is requested.
	Connected(context.Context, swarm.Address) error
}

type Disconnecter interface {
	// Disconnected is called when a peer disconnects.
	// The disconnect event can be initiated on the local
	// node or on the remote node, this handle does not make
	// any distinctions between either of them.
	Disconnected(swarm.Address)
}

type ClosestPeerer interface {
	ClosestPeer(addr swarm.Address) (peerAddr swarm.Address, err error)
}

type EachPeerer interface {
	// EachPeer iterates from closest bin to farthest
	EachPeer(EachPeerFunc) error
	// EachPeerRev iterates from farthest bin to closest
	EachPeerRev(EachPeerFunc) error
}

// EachPeerFunc is a callback that is called with a peer and its PO
type EachPeerFunc func(swarm.Address, uint8) (stop, jumpToNext bool, err error)
