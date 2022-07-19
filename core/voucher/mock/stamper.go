package mock

import (
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/voucher"
)

type mockStamper struct{}

// NewStamper returns anew new mock stamper.
func NewStamper() voucher.Stamper {
	return &mockStamper{}
}

// Stamp implements the Stamper interface. It returns an empty voucher stamp.
func (mockStamper) Stamp(_ swarm.Address) (*voucher.Stamp, error) {
	return &voucher.Stamp{}, nil
}
