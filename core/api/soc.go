package api

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/bmtpool"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/soc"
	"github.com/redesblock/hop/core/swarm"
)

var (
	errBadRequestParams = errors.New("owner, id or span is not well formed")
)

type socPostResponse struct {
	Reference swarm.Address `json:"reference"`
}

func (s *server) socUploadHandler(w http.ResponseWriter, r *http.Request) {
	owner, err := hex.DecodeString(mux.Vars(r)["owner"])
	if err != nil {
		s.Logger.Debugf("soc upload: bad owner: %v", err)
		s.Logger.Error("soc upload: %v", errBadRequestParams)
		jsonhttp.BadRequest(w, "bad owner")
		return
	}
	id, err := hex.DecodeString(mux.Vars(r)["id"])
	if err != nil {
		s.Logger.Debugf("soc upload: bad id: %v", err)
		s.Logger.Error("soc upload: %v", errBadRequestParams)
		jsonhttp.BadRequest(w, "bad id")
		return
	}

	sigStr := r.URL.Query().Get("sig")
	if sigStr == "" {
		s.Logger.Debugf("soc upload: empty signature")
		s.Logger.Error("soc upload: empty signature")
		jsonhttp.BadRequest(w, "empty signature")
		return
	}

	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		s.Logger.Debugf("soc upload: bad signature: %v", err)
		s.Logger.Error("soc upload: bad signature")
		jsonhttp.BadRequest(w, "bad signature")
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if jsonhttp.HandleBodyReadError(err, w) {
			return
		}
		s.Logger.Debugf("soc upload: read chunk data error: %v", err)
		s.Logger.Error("soc upload: read chunk data error")
		jsonhttp.InternalServerError(w, "cannot read chunk data")
		return
	}

	if len(data) < swarm.SpanSize {
		s.Logger.Debugf("soc upload: chunk data too short")
		s.Logger.Error("soc upload: %v", errBadRequestParams)
		jsonhttp.BadRequest(w, "short chunk data")
		return
	}

	if len(data) > swarm.ChunkSize+swarm.SpanSize {
		s.Logger.Debugf("soc upload: chunk data exceeds %d bytes", swarm.ChunkSize+swarm.SpanSize)
		s.Logger.Error("soc upload: chunk data error")
		jsonhttp.RequestEntityTooLarge(w, "payload too large")
		return
	}

	ch, err := chunk(data)
	if err != nil {
		s.Logger.Debugf("soc upload: create content addressed chunk: %v", err)
		s.Logger.Error("soc upload: chunk data error")
		jsonhttp.BadRequest(w, "chunk data error")
		return
	}

	chunk, err := soc.NewSignedChunk(id, ch, owner, sig)
	if err != nil {
		s.Logger.Debugf("soc upload: read chunk data error: %v", err)
		s.Logger.Error("soc upload: read chunk data error")
		jsonhttp.InternalServerError(w, "cannot read chunk data")
		return
	}

	if !soc.Valid(chunk) {
		s.Logger.Debugf("soc upload: invalid chunk: %v", err)
		s.Logger.Error("soc upload: invalid chunk")
		jsonhttp.Unauthorized(w, "invalid chunk")
		return

	}
	ctx := r.Context()

	_, err = s.Storer.Put(ctx, requestModePut(r), chunk)
	if err != nil {
		s.Logger.Debugf("soc upload: chunk write error: %v", err)
		s.Logger.Error("soc upload: chunk write error")
		jsonhttp.BadRequest(w, "chunk write error")
		return
	}

	jsonhttp.Created(w, chunkAddressResponse{Reference: chunk.Address()})
}

func chunk(data []byte) (swarm.Chunk, error) {
	hasher := bmtpool.Get()
	defer bmtpool.Put(hasher)
	err := hasher.SetSpanBytes(data[:swarm.SpanSize])
	if err != nil {
		return nil, err
	}
	_, err = hasher.Write(data[swarm.SpanSize:])
	if err != nil {
		return nil, err
	}

	return swarm.NewChunk(swarm.NewAddress(hasher.Sum(nil)), data), nil
}
