package mock

import (
	"context"
	"math/big"

	"github.com/redesblock/hop/core/postage/postagecontract"
)

type contractMock struct {
	createBatch func(ctx context.Context, initialBalance *big.Int, depth uint8, label string) ([]byte, error)
}

func (c *contractMock) CreateBatch(ctx context.Context, initialBalance *big.Int, depth uint8, label string) ([]byte, error) {
	return c.createBatch(ctx, initialBalance, depth, label)
}

// Option is a an option passed to New
type Option func(*contractMock)

// New creates a new mock BatchStore
func New(opts ...Option) postagecontract.Interface {
	bs := &contractMock{}

	for _, o := range opts {
		o(bs)
	}

	return bs
}

func WithCreateBatchFunc(f func(ctx context.Context, initialBalance *big.Int, depth uint8, label string) ([]byte, error)) Option {
	return func(m *contractMock) {
		m.createBatch = f
	}
}
