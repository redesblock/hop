package mock

import (
	"context"
	"errors"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/swarm"
)

type Service struct {
	addProtocolFunc         func(p2p.ProtocolSpec) error
	connectFunc             func(ctx context.Context, addr ma.Multiaddr) (overlay swarm.Address, err error)
	disconnectFunc          func(overlay swarm.Address) error
	peersFunc               func() []p2p.Peer
	setPeerAddedHandlerFunc func(func(context.Context, swarm.Address) error)
	addressesFunc           func() ([]ma.Multiaddr, error)
}

func WithAddProtocolFunc(f func(p2p.ProtocolSpec) error) Option {
	return optionFunc(func(s *Service) {
		s.addProtocolFunc = f
	})
}

func WithConnectFunc(f func(ctx context.Context, addr ma.Multiaddr) (overlay swarm.Address, err error)) Option {
	return optionFunc(func(s *Service) {
		s.connectFunc = f
	})
}

func WithDisconnectFunc(f func(overlay swarm.Address) error) Option {
	return optionFunc(func(s *Service) {
		s.disconnectFunc = f
	})
}

func WithPeersFunc(f func() []p2p.Peer) Option {
	return optionFunc(func(s *Service) {
		s.peersFunc = f
	})
}

func WithSetPeerAddedHandlerFunc(f func(func(context.Context, swarm.Address) error)) Option {
	return optionFunc(func(s *Service) {
		s.setPeerAddedHandlerFunc = f
	})
}

func WithAddressesFunc(f func() ([]ma.Multiaddr, error)) Option {
	return optionFunc(func(s *Service) {
		s.addressesFunc = f
	})
}

func New(opts ...Option) *Service {
	s := new(Service)
	for _, o := range opts {
		o.apply(s)
	}
	return s
}

func (s *Service) AddProtocol(spec p2p.ProtocolSpec) error {
	if s.addProtocolFunc == nil {
		return errors.New("function AddProtocol not configured")
	}
	return s.addProtocolFunc(spec)
}

func (s *Service) Connect(ctx context.Context, addr ma.Multiaddr) (overlay swarm.Address, err error) {
	if s.connectFunc == nil {
		return swarm.Address{}, errors.New("function Connect not configured")
	}
	return s.connectFunc(ctx, addr)
}

func (s *Service) Disconnect(overlay swarm.Address) error {
	if s.disconnectFunc == nil {
		return errors.New("function Disconnect not configured")
	}
	return s.disconnectFunc(overlay)
}

func (s *Service) SetPeerAddedHandler(f func(context.Context, swarm.Address) error) {
	if s.setPeerAddedHandlerFunc == nil {
		return
	}

	s.setPeerAddedHandlerFunc(f)
}

func (s *Service) Addresses() ([]ma.Multiaddr, error) {
	if s.addressesFunc == nil {
		return nil, errors.New("function Addresses not configured")
	}
	return s.addressesFunc()
}

func (s *Service) Peers() []p2p.Peer {
	if s.peersFunc == nil {
		return nil
	}
	return s.peersFunc()
}

type Option interface {
	apply(*Service)
}
type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }