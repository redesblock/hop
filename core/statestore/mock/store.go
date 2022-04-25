package mock

import (
	"encoding"
	"encoding/json"
	"strings"
	"sync"

	"github.com/redesblock/hop/core/storage"
)

var _ storage.StateStorer = (*store)(nil)

type store struct {
	store map[string][]byte
	mtx   sync.Mutex
}

func NewStateStore() storage.StateStorer {
	return &store{
		store: make(map[string][]byte),
	}
}

func (s *store) Get(key string, i interface{}) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	data, ok := s.store[key]
	if !ok {
		return storage.ErrNotFound
	}

	if unmarshaler, ok := i.(encoding.BinaryUnmarshaler); ok {
		return unmarshaler.UnmarshalBinary(data)
	}

	return json.Unmarshal(data, i)
}

func (s *store) Put(key string, i interface{}) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var bytes []byte
	if marshaler, ok := i.(encoding.BinaryMarshaler); ok {
		if bytes, err = marshaler.MarshalBinary(); err != nil {
			return err
		}
	} else if bytes, err = json.Marshal(i); err != nil {
		return err
	}

	s.store[key] = bytes
	return nil
}

func (s *store) Delete(key string) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.store, key)
	return nil
}

func (s *store) Iterate(prefix string, iterFunc storage.StateIterFunc) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for k, v := range s.store {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		val := make([]byte, len(v))
		copy(val, v)
		stop, err := iterFunc([]byte(k), val)
		if err != nil {
			return err
		}

		if stop {
			return nil
		}
	}
	return nil
}

func (s *store) Close() (err error) {
	return nil
}