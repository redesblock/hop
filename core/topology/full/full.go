package full

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/redesblock/hop/core/addressbook"
	"github.com/redesblock/hop/core/discovery"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/topology"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ topology.Driver = (*driver)(nil)

// driver drives the connectivity between nodes. It is a basic implementation of a connectivity driver.
// that enabled full connectivity in the sense that:
// - Every peer which is added to the driver gets broadcasted to every other peer regardless of its address.
// - A random peer is picked when asking for a peer to retrieve an arbitrary chunk (Peerer interface).
type driver struct {
	base swarm.Address // the base address for this node

	discovery     discovery.Driver
	addressBook   addressbook.GetPutter
	p2pService    p2p.Service
	receivedPeers map[string]struct{} // track already received peers. Note: implement cleanup or expiration if needed to stop infinite grow
	mtx           sync.Mutex          // guards received peers
	logger        logging.Logger
}

func New(disc discovery.Driver, addressBook addressbook.GetPutter, p2pService p2p.Service, logger logging.Logger, baseAddress swarm.Address) topology.Driver {
	return &driver{
		base:          baseAddress,
		discovery:     disc,
		addressBook:   addressBook,
		p2pService:    p2pService,
		receivedPeers: make(map[string]struct{}),
		logger:        logger,
	}
}

// AddPeer adds a new peer to the topology driver.
// The peer would be subsequently broadcasted to all connected peers.
// All conneceted peers are also broadcasted to the new peer.
func (d *driver) AddPeer(ctx context.Context, addr swarm.Address) error {
	d.mtx.Lock()
	if _, ok := d.receivedPeers[addr.ByteString()]; ok {
		d.mtx.Unlock()
		return nil
	}

	d.receivedPeers[addr.ByteString()] = struct{}{}
	d.mtx.Unlock()

	connectedPeers := d.p2pService.Peers()
	ma, err := d.addressBook.Get(addr)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return topology.ErrNotFound
		}
		return err
	}

	if !isConnected(addr, connectedPeers) {
		peerAddr, err := d.p2pService.Connect(ctx, ma)
		if err != nil {
			return err
		}

		// update addr if it is wrong or it has been changed
		if !addr.Equal(peerAddr) {
			addr = peerAddr
			err := d.addressBook.Put(peerAddr, ma)
			if err != nil {
				return err
			}
		}
	}

	connectedAddrs := []swarm.Address{}
	for _, addressee := range connectedPeers {
		// skip newly added peer
		if addressee.Address.Equal(addr) {
			continue
		}

		connectedAddrs = append(connectedAddrs, addressee.Address)
		if err := d.discovery.BroadcastPeers(context.Background(), addressee.Address, addr); err != nil {
			return err
		}
	}

	if len(connectedAddrs) == 0 {
		return nil
	}

	if err := d.discovery.BroadcastPeers(context.Background(), addr, connectedAddrs...); err != nil {
		return err
	}

	return nil
}

// ChunkPeer is used to suggest a peer to ask a certain chunk from.
func (d *driver) ChunkPeer(addr swarm.Address) (peerAddr swarm.Address, err error) {
	connectedPeers := d.p2pService.Peers()
	if len(connectedPeers) == 0 {
		return swarm.Address{}, topology.ErrNotFound
	}

	itemIdx := rand.Intn(len(connectedPeers))
	i := 0
	for _, v := range connectedPeers {
		if i == itemIdx {
			return v.Address, nil
		}
		i++
	}

	return swarm.Address{}, topology.ErrNotFound
}

// SyncPeer returns a peer to which we would like to sync an arbitrary
// chunk address. Returns the closest peer in relation to the chunk.
func (d *driver) SyncPeer(addr swarm.Address) (swarm.Address, error) {
	connectedPeers := d.p2pService.Peers()
	if len(connectedPeers) == 0 {
		return swarm.Address{}, topology.ErrNotFound
	}

	overlays := make([]swarm.Address, len(connectedPeers))
	for i, v := range connectedPeers {
		overlays[i] = v.Address
	}

	return closestPeer(addr, d.base, overlays)
}

// closestPeer returns the closest peer from the supplied peers slice.
// returns topology.ErrWantSelf if the base address is the closest
func closestPeer(addr, self swarm.Address, peers []swarm.Address) (swarm.Address, error) {
	// start checking closest from _self_
	closest := self
	for _, peer := range peers {
		dcmp, err := swarm.DistanceCmp(addr.Bytes(), closest.Bytes(), peer.Bytes())
		if err != nil {
			return swarm.Address{}, err
		}
		switch dcmp {
		case 0:
			// do nothing
		case -1:
			// current peer is closer
			closest = peer
		case 1:
			// closest is already closer to chunk
			// do nothing
		}
	}

	// check if self
	if closest.Equal(self) {
		return swarm.Address{}, topology.ErrWantSelf
	}

	return closest, nil
}

func isConnected(addr swarm.Address, connectedPeers []p2p.Peer) bool {
	for _, p := range connectedPeers {
		if p.Address.Equal(addr) {
			return true
		}
	}

	return false
}
