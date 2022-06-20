package retrieval

import (
	"context"

	"github.com/redesblock/hop/core/p2p"
)

func (s *Service) Handler(ctx context.Context, p p2p.Peer, stream p2p.Stream) error {
	return s.handler(ctx, p, stream)
}
