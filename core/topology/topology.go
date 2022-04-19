package topology

import (
	"context"
	"errors"
	"io"

	"github.com/redesblock/hop/core/swarm"
)

var ErrNotFound = errors.New("no peer found")
var ErrWantSelf = errors.New("node wants self")

type Driver interface {
	PeerAdder
	ClosestPeerer
	io.Closer
}

type PeerAdder interface {
	AddPeer(ctx context.Context, addr swarm.Address) error
}

type ClosestPeerer interface {
	ClosestPeer(addr swarm.Address) (peerAddr swarm.Address, err error)
}

// EachPeerFunc is a callback that is called with a peer and its PO
type EachPeerFunc func(swarm.Address, uint8) (stop, jumpToNext bool, err error)
