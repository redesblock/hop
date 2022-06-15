package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/file/pipeline/builder"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/sctx"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tracing"
)

type bytesPostResponse struct {
	Reference swarm.Address `json:"reference"`
}

// bytesUploadHandler handles upload of raw binary data of arbitrary length.
func (s *server) bytesUploadHandler(w http.ResponseWriter, r *http.Request) {
	logger := tracing.NewLoggerWithTraceID(r.Context(), s.Logger)

	tag, created, err := s.getOrCreateTag(r.Header.Get(SwarmTagUidHeader))
	if err != nil {
		logger.Debugf("bytes upload: get or create tag: %v", err)
		logger.Error("bytes upload: get or create tag")
		jsonhttp.InternalServerError(w, "cannot get or create tag")
		return
	}

	// Add the tag to the context
	ctx := sctx.SetTag(r.Context(), tag)

	pipe := builder.NewPipelineBuilder(ctx, s.Storer, requestModePut(r), requestEncrypt(r))
	address, err := builder.FeedPipeline(ctx, pipe, r.Body, r.ContentLength)
	if err != nil {
		logger.Debugf("bytes upload: split write all: %v", err)
		logger.Error("bytes upload: split write all")
		jsonhttp.InternalServerError(w, nil)
		return
	}
	if created {
		_, err = tag.DoneSplit(address)
		if err != nil {
			logger.Debugf("bytes upload: done split: %v", err)
			logger.Error("bytes upload: done split failed")
			jsonhttp.InternalServerError(w, nil)
			return
		}
	}
	w.Header().Set(SwarmTagUidHeader, fmt.Sprint(tag.Uid))
	w.Header().Set("Access-Control-Expose-Headers", SwarmTagUidHeader)
	jsonhttp.OK(w, bytesPostResponse{
		Reference: address,
	})
}

// bytesGetHandler handles retrieval of raw binary data of arbitrary length.
func (s *server) bytesGetHandler(w http.ResponseWriter, r *http.Request) {
	logger := tracing.NewLoggerWithTraceID(r.Context(), s.Logger).Logger
	nameOrHex := mux.Vars(r)["address"]

	address, err := s.resolveNameOrAddress(nameOrHex)
	if err != nil {
		logger.Debugf("bytes: parse address %s: %v", nameOrHex, err)
		logger.Error("bytes: parse address error")
		jsonhttp.BadRequest(w, "invalid address")
		return
	}

	additionalHeaders := http.Header{
		"Content-Type": {"application/octet-stream"},
	}

	s.downloadHandler(w, r, address, additionalHeaders)
}
