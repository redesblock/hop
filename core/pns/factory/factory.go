package factory

import (
	"github.com/redesblock/hop/core/pns"
	"github.com/redesblock/hop/core/pns/epochs"
	"github.com/redesblock/hop/core/pns/sequence"
	"github.com/redesblock/hop/core/storage"
)

type factory struct {
	storage.Getter
}

func New(getter storage.Getter) pns.Factory {
	return &factory{getter}
}

func (f *factory) NewLookup(t pns.Type, feed *pns.Feed) (pns.Lookup, error) {
	switch t {
	case pns.Sequence:
		return sequence.NewAsyncFinder(f.Getter, feed), nil
	case pns.Epoch:
		return epochs.NewAsyncFinder(f.Getter, feed), nil
	}

	return nil, pns.ErrFeedTypeNotFound
}
