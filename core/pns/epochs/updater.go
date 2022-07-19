package epochs

import (
	"context"

	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/pns"
	"github.com/redesblock/hop/core/storage"
)

var _ pns.Updater = (*updater)(nil)

// Updater encapsulates a pns putter to generate successive updates for epoch based pns
// it persists the last update
type updater struct {
	*pns.Putter
	last  int64
	epoch pns.Index
}

// NewUpdater constructs a feed updater
func NewUpdater(putter storage.Putter, signer crypto.Signer, topic []byte) (pns.Updater, error) {
	p, err := pns.NewPutter(putter, signer, topic)
	if err != nil {
		return nil, err
	}
	return &updater{Putter: p}, nil
}

// Update pushes an update to the feed through the chunk stores
func (u *updater) Update(ctx context.Context, at int64, payload []byte) error {
	e := next(u.epoch, u.last, uint64(at))
	err := u.Put(ctx, e, at, payload)
	if err != nil {
		return err
	}
	u.last = at
	u.epoch = e
	return nil
}

func (u *updater) Feed() *pns.Feed {
	return u.Putter.Feed
}
