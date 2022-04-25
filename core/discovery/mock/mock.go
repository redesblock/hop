package mock

import (
	"context"
	"sync"

	"github.com/redesblock/hop/core/swarm"
)

type Discovery struct {
	mtx     sync.Mutex
	ctr     int //how many ops
	records map[string][]swarm.Address
}

func NewDiscovery() *Discovery {
	return &Discovery{
		records: make(map[string][]swarm.Address),
	}
}

func (d *Discovery) BroadcastPeers(ctx context.Context, addressee swarm.Address, peers ...swarm.Address) error {
	for _, peer := range peers {
		d.mtx.Lock()
		d.records[addressee.String()] = append(d.records[addressee.String()], peer)
		d.mtx.Unlock()
	}

	d.mtx.Lock()
	d.ctr++
	d.mtx.Unlock()
	return nil
}

func (d *Discovery) Broadcasts() int {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	return d.ctr
}

func (d *Discovery) AddresseeRecords(addressee swarm.Address) (peers []swarm.Address, exists bool) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	peers, exists = d.records[addressee.String()]
	return
}