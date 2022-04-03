package mem

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/keystore"
)

var _ keystore.Service = (*Service)(nil)

type Service struct {
	m  map[string]key
	mu sync.Mutex
}

func New() *Service {
	return &Service{
		m: make(map[string]key),
	}
}

func (s *Service) Key(name, password string) (pk *ecdsa.PrivateKey, created bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k, ok := s.m[name]
	if !ok {
		pk, err := crypto.GenerateSecp256k1Key()
		if err != nil {
			return nil, false, fmt.Errorf("generate secp256k1 key: %w", err)
		}
		s.m[name] = key{
			pk:       pk,
			password: password,
		}
		return pk, true, nil
	}

	if k.password != password {
		return nil, false, keystore.ErrInvalidPassword
	}

	return k.pk, created, nil
}

type key struct {
	pk       *ecdsa.PrivateKey
	password string
}
