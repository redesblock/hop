package mock

import (
	"github.com/redesblock/hop/core/swarm"
)

type MockPricer struct {
	peerPrice uint64
	price     uint64
}

func NewMockService(price, peerPrice uint64) *MockPricer {
	return &MockPricer{
		peerPrice: peerPrice,
		price:     price,
	}
}

func (pricer *MockPricer) PeerPrice(peer, chunk swarm.Address) uint64 {
	return pricer.peerPrice
}

func (pricer *MockPricer) Price(chunk swarm.Address) uint64 {
	return pricer.price
}
