package feeds

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/redesblock/hop/core/soc"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
)

// Lookup is the interface for time based feed lookup
type Lookup interface {
	At(ctx context.Context, at, after int64) (swarm.Chunk, error)
}

// Getter encapsulates a chunk Getter getter and a feed and provides
//  non-concurrent lookup methods
type Getter struct {
	getter storage.Getter
	*Feed
}

// NewGetter constructs a feed Getter
func NewGetter(getter storage.Getter, feed *Feed) *Getter {
	return &Getter{getter, feed}
}

// Latest looks up the latest update of the feed
// after is a unix time hint of the latest known update
func Latest(ctx context.Context, l Lookup, after int64) (swarm.Chunk, error) {
	return l.At(ctx, time.Now().Unix(), after)
}

// Get creates an update of the underlying feed at the given epoch
// and looks it up in the chunk Getter based on its address
func (f *Getter) Get(ctx context.Context, i Index) (swarm.Chunk, error) {
	addr, err := f.Feed.Update(i).Address()
	if err != nil {
		return nil, err
	}
	return f.getter.Get(ctx, storage.ModeGetRequest, addr)
}

// FromChunk parses out the timestamp and the payload
func FromChunk(ch swarm.Chunk) (uint64, []byte, error) {
	s, err := soc.FromChunk(ch)
	if err != nil {
		return 0, nil, err
	}
	cac := s.Chunk
	if len(cac.Data()) < 16 {
		return 0, nil, fmt.Errorf("feed update payload too short")
	}
	payload := cac.Data()[16:]
	at := binary.BigEndian.Uint64(cac.Data()[8:16])
	return at, payload, nil
}

// UpdatedAt extracts the time of feed other than update
func UpdatedAt(ch swarm.Chunk) (uint64, error) {
	d := ch.Data()
	if len(d) < 113 {
		return 0, fmt.Errorf("too short: %d", len(d))
	}
	// a soc chunk with time information in the wrapped content addressed chunk
	// 0-32    index,
	// 65-97   signature,
	// 98-105  span of wrapped chunk
	// 105-113 timestamp
	return binary.BigEndian.Uint64(d[105:113]), nil
}
