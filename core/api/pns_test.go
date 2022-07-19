package api_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/file/loadsave"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/manifest"
	"github.com/redesblock/hop/core/pns"
	testingsoc "github.com/redesblock/hop/core/soc/testing"
	statestore "github.com/redesblock/hop/core/statestore/mock"
	"github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
	"github.com/redesblock/hop/core/util/logging"
	"github.com/redesblock/hop/core/voucher"
	mockpost "github.com/redesblock/hop/core/voucher/mock"
)

const ownerString = "8d3766440f0d7b949a5e32995d09619a7f86e632"

var expReference = swarm.MustParseHexAddress("891a1d1c8436c792d02fc2e8883fef7ab387eaeaacd25aa9f518be7be7856d54")

func TestFeed_Get(t *testing.T) {
	var (
		feedResource = func(owner, topic, at string) string {
			if at != "" {
				return fmt.Sprintf("/pns/%s/%s?at=%s", owner, topic, at)
			}
			return fmt.Sprintf("/pns/%s/%s", owner, topic)
		}
		mockStatestore  = statestore.NewStateStore()
		logger          = logging.New(io.Discard, 0)
		tag             = tags.NewTags(mockStatestore, logger)
		mockStorer      = mock.NewStorer()
		client, _, _, _ = newTestServer(t, testServerOptions{
			Storer: mockStorer,
			Tags:   tag,
		})
	)

	t.Run("malformed owner", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodGet, feedResource("xyz", "cc", ""), http.StatusBadRequest,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "bad owner",
				Code:    http.StatusBadRequest,
			}),
		)
	})

	t.Run("malformed topic", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodGet, feedResource("8d3766440f0d7b949a5e32995d09619a7f86e632", "xxzzyy", ""), http.StatusBadRequest,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "bad topic",
				Code:    http.StatusBadRequest,
			}),
		)
	})

	t.Run("at malformed", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodGet, feedResource("8d3766440f0d7b949a5e32995d09619a7f86e632", "aabbcc", "unbekannt"), http.StatusBadRequest,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "bad at",
				Code:    http.StatusBadRequest,
			}),
		)
	})

	t.Run("with at", func(t *testing.T) {
		var (
			timestamp       = int64(12121212)
			ch              = toChunk(t, uint64(timestamp), expReference.Bytes())
			look            = newMockLookup(12, 0, ch, nil, &id{}, &id{})
			factory         = newMockFactory(look)
			idBytes, _      = (&id{}).MarshalBinary()
			client, _, _, _ = newTestServer(t, testServerOptions{
				Storer: mockStorer,
				Tags:   tag,
				Feeds:  factory,
			})
		)

		respHeaders := jsonhttptest.Request(t, client, http.MethodGet, feedResource(ownerString, "aabbcc", "12"), http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(api.FeedReferenceResponse{Reference: expReference}),
		)

		h := respHeaders[api.SwarmFeedIndexHeader]
		if len(h) == 0 {
			t.Fatal("expected swarm feed index header to be set")
		}
		b, err := hex.DecodeString(h[0])
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b, idBytes) {
			t.Fatalf("feed index header mismatch. got %v want %v", b, idBytes)
		}
	})

	t.Run("latest", func(t *testing.T) {
		var (
			timestamp  = int64(12121212)
			ch         = toChunk(t, uint64(timestamp), expReference.Bytes())
			look       = newMockLookup(-1, 2, ch, nil, &id{}, &id{})
			factory    = newMockFactory(look)
			idBytes, _ = (&id{}).MarshalBinary()

			client, _, _, _ = newTestServer(t, testServerOptions{
				Storer: mockStorer,
				Tags:   tag,
				Feeds:  factory,
			})
		)

		respHeaders := jsonhttptest.Request(t, client, http.MethodGet, feedResource(ownerString, "aabbcc", ""), http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(api.FeedReferenceResponse{Reference: expReference}),
		)

		if h := respHeaders[api.SwarmFeedIndexHeader]; len(h) > 0 {
			b, err := hex.DecodeString(h[0])
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, idBytes) {
				t.Fatalf("feed index header mismatch. got %v want %v", b, idBytes)
			}
		} else {
			t.Fatal("expected swarm feed index header to be set")
		}
	})
}

func TestFeed_Post(t *testing.T) {
	// post to owner, tpoic, then expect a reference
	// get the reference from the store, unmarshal to a
	// manifest entry and make sure all metadata correct
	var (
		mockStatestore  = statestore.NewStateStore()
		logger          = logging.New(io.Discard, 0)
		tag             = tags.NewTags(mockStatestore, logger)
		topic           = "aabbcc"
		mp              = mockpost.New(mockpost.WithIssuer(voucher.NewStampIssuer("", "", batchOk, big.NewInt(3), 11, 10, 1000, true)))
		mockStorer      = mock.NewStorer()
		client, _, _, _ = newTestServer(t, testServerOptions{
			Storer: mockStorer,
			Tags:   tag,
			Logger: logger,
			Post:   mp,
		})
		url = fmt.Sprintf("/pns/%s/%s?type=%s", ownerString, topic, "sequence")
	)

	t.Run("ok", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodPost, url, http.StatusCreated,
			jsonhttptest.WithRequestHeader(api.SwarmDeferredUploadHeader, "true"),
			jsonhttptest.WithRequestHeader(api.SwarmPostageBatchIdHeader, batchOkStr),
			jsonhttptest.WithExpectedJSONResponse(api.FeedReferenceResponse{
				Reference: expReference,
			}),
		)

		ls := loadsave.NewReadonly(mockStorer)
		i, err := manifest.NewMantarayManifestReference(expReference, ls)
		if err != nil {
			t.Fatal(err)
		}
		e, err := i.Lookup(context.Background(), "/")
		if err != nil {
			t.Fatal(err)
		}

		meta := e.Metadata()
		if e := meta[api.FeedMetadataEntryOwner]; e != ownerString {
			t.Fatalf("owner mismatch. got %s want %s", e, ownerString)
		}
		if e := meta[api.FeedMetadataEntryTopic]; e != topic {
			t.Fatalf("topic mismatch. got %s want %s", e, topic)
		}
		if e := meta[api.FeedMetadataEntryType]; e != "Sequence" {
			t.Fatalf("type mismatch. got %s want %s", e, "Sequence")
		}
	})
	t.Run("voucher", func(t *testing.T) {
		t.Run("err - bad batch", func(t *testing.T) {
			hexbatch := hex.EncodeToString(batchInvalid)
			jsonhttptest.Request(t, client, http.MethodPost, url, http.StatusBadRequest,
				jsonhttptest.WithRequestHeader(api.SwarmPostageBatchIdHeader, hexbatch),
				jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
					Message: "invalid voucher batch id",
					Code:    http.StatusBadRequest,
				}))
		})

		t.Run("ok - batch zeros", func(t *testing.T) {
			hexbatch := hex.EncodeToString(batchOk)
			jsonhttptest.Request(t, client, http.MethodPost, url, http.StatusCreated,
				jsonhttptest.WithRequestHeader(api.SwarmDeferredUploadHeader, "true"),
				jsonhttptest.WithRequestHeader(api.SwarmPostageBatchIdHeader, hexbatch),
			)
		})
		t.Run("bad request - batch empty", func(t *testing.T) {
			hexbatch := hex.EncodeToString(batchEmpty)
			jsonhttptest.Request(t, client, http.MethodPost, url, http.StatusBadRequest,
				jsonhttptest.WithRequestHeader(api.SwarmPostageBatchIdHeader, hexbatch),
			)
		})
	})

}

type factoryMock struct {
	sequenceCalled bool
	epochCalled    bool
	feed           *pns.Feed
	lookup         pns.Lookup
}

func newMockFactory(mockLookup pns.Lookup) *factoryMock {
	return &factoryMock{lookup: mockLookup}
}

func (f *factoryMock) NewLookup(t pns.Type, feed *pns.Feed) (pns.Lookup, error) {
	switch t {
	case pns.Sequence:
		f.sequenceCalled = true
	case pns.Epoch:
		f.epochCalled = true
	}
	f.feed = feed
	return f.lookup, nil
}

type mockLookup struct {
	at, after int64
	chunk     swarm.Chunk
	err       error
	cur, next pns.Index
}

func newMockLookup(at, after int64, ch swarm.Chunk, err error, cur, next pns.Index) *mockLookup {
	return &mockLookup{at: at, after: after, chunk: ch, err: err, cur: cur, next: next}
}

func (l *mockLookup) At(_ context.Context, at, after int64) (swarm.Chunk, pns.Index, pns.Index, error) {
	if l.at == -1 {
		// shortcut to ignore the value in the call since time.Now() is a moving target
		return l.chunk, l.cur, l.next, nil
	}
	if at == l.at && after == l.after {
		return l.chunk, l.cur, l.next, nil
	}
	return nil, nil, nil, errors.New("no feed update found")
}

func toChunk(t *testing.T, at uint64, payload []byte) swarm.Chunk {
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, at)
	content := append(ts, payload...)

	s := testingsoc.GenerateMockSOC(t, content)
	return s.Chunk()
}

type id struct{}

func (i *id) MarshalBinary() ([]byte, error) {
	return []byte("accd"), nil
}

func (i *id) String() string {
	return "44237"
}

func (*id) Next(last int64, at uint64) pns.Index {
	return &id{}
}
