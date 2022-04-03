package file

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/redesblock/hop/core/crypto"
)

type Service struct {
	dir string
}

func New(dir string) *Service {
	return &Service{dir: dir}
}

func (s *Service) Key(name, password string) (pk *ecdsa.PrivateKey, created bool, err error) {
	filename := s.keyFilename(name)

	data, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, fmt.Errorf("read private key: %w", err)
	}
	if len(data) == 0 {
		var err error
		pk, err = crypto.GenerateSecp256k1Key()
		if err != nil {
			return nil, false, fmt.Errorf("generate secp256k1 key: %w", err)
		}

		d, err := encryptKey(pk, password)
		if err != nil {
			return nil, false, err
		}

		if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
			return nil, false, err
		}
		if err := ioutil.WriteFile(filename, d, 0600); err != nil {
			return nil, false, err
		}
		return pk, true, nil
	}

	pk, err = decryptKey(data, password)
	if err != nil {
		return nil, false, err
	}
	return pk, false, nil
}

func (s *Service) keyFilename(name string) string {
	return filepath.Join(s.dir, fmt.Sprintf("%s.key", name))
}
