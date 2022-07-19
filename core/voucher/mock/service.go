package mock

import (
	"errors"
	"math/big"
	"sync"

	"github.com/redesblock/hop/core/voucher"
)

type optionFunc func(*mockPostage)

// Option is an option passed to a mock voucher Service.
type Option interface {
	apply(*mockPostage)
}

func (f optionFunc) apply(r *mockPostage) { f(r) }

// New creates a new mock voucher service.
func New(o ...Option) voucher.Service {
	m := &mockPostage{
		issuersMap: make(map[string]*voucher.StampIssuer),
	}
	for _, v := range o {
		v.apply(m)
	}

	return m
}

// WithAcceptAll sets the mock to return a new BatchIssuer on every
// call to GetStampIssuer.
func WithAcceptAll() Option {
	return optionFunc(func(m *mockPostage) { m.acceptAll = true })
}

func WithIssuer(s *voucher.StampIssuer) Option {
	return optionFunc(func(m *mockPostage) {
		m.issuersMap = map[string]*voucher.StampIssuer{string(s.ID()): s}
	})
}

type mockPostage struct {
	issuersMap map[string]*voucher.StampIssuer
	issuerLock sync.Mutex
	acceptAll  bool
}

func (m *mockPostage) Add(s *voucher.StampIssuer) error {
	m.issuerLock.Lock()
	defer m.issuerLock.Unlock()

	m.issuersMap[string(s.ID())] = s
	return nil
}

func (m *mockPostage) StampIssuers() []*voucher.StampIssuer {
	m.issuerLock.Lock()
	defer m.issuerLock.Unlock()

	issuers := []*voucher.StampIssuer{}
	for _, v := range m.issuersMap {
		issuers = append(issuers, v)
	}
	return issuers
}

func (m *mockPostage) GetStampIssuer(id []byte) (*voucher.StampIssuer, error) {
	if m.acceptAll {
		return voucher.NewStampIssuer("test fallback", "test identity", id, big.NewInt(3), 24, 6, 1000, true), nil
	}

	m.issuerLock.Lock()
	defer m.issuerLock.Unlock()

	i, exists := m.issuersMap[string(id)]
	if !exists {
		return nil, errors.New("stampissuer not found")
	}
	return i, nil
}

func (m *mockPostage) IssuerUsable(_ *voucher.StampIssuer) bool {
	return true
}

func (m *mockPostage) HandleCreate(_ *voucher.Batch) error { return nil }

func (m *mockPostage) HandleTopUp(_ []byte, _ *big.Int) {}

func (m *mockPostage) HandleDepthIncrease(_ []byte, _ uint8, _ *big.Int) {}

func (m *mockPostage) Close() error {
	return nil
}
