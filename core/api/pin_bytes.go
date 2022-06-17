package api

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/traversal"
)

// pinBytes is used to pin an already uploaded content.
func (s *server) pinBytes(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.Logger.Debugf("pin bytes: parse address: %v", err)
		s.Logger.Error("pin bytes: parse address")
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.Storer.Has(r.Context(), addr)
	if err != nil {
		s.Logger.Debugf("pin bytes: localstore has: %v", err)
		s.Logger.Error("pin bytes: store")
		jsonhttp.InternalServerError(w, err)
		return
	}

	if !has {
		jsonhttp.NotFound(w, nil)
		return
	}

	ctx := r.Context()

	chunkAddressFn := s.pinChunkAddressFn(ctx, addr)

	err = s.Traversal.TraverseBytesAddresses(ctx, addr, chunkAddressFn)
	if err != nil {
		s.Logger.Debugf("pin bytes: traverse chunks: %v, addr %s", err, addr)

		if errors.Is(err, traversal.ErrInvalidType) {
			s.Logger.Error("pin bytes: invalid type")
			jsonhttp.BadRequest(w, "invalid type")
			return
		}

		s.Logger.Error("pin bytes: cannot pin")
		jsonhttp.InternalServerError(w, "cannot pin")
		return
	}

	jsonhttp.OK(w, nil)
}

// unpinBytes removes pinning from content.
func (s *server) unpinBytes(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.Logger.Debugf("pin bytes: parse address: %v", err)
		s.Logger.Error("pin bytes: parse address")
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.Storer.Has(r.Context(), addr)
	if err != nil {
		s.Logger.Debugf("pin bytes: localstore has: %v", err)
		s.Logger.Error("pin bytes: store")
		jsonhttp.InternalServerError(w, err)
		return
	}

	if !has {
		jsonhttp.NotFound(w, nil)
		return
	}

	ctx := r.Context()

	chunkAddressFn := s.unpinChunkAddressFn(ctx, addr)

	err = s.Traversal.TraverseBytesAddresses(ctx, addr, chunkAddressFn)
	if err != nil {
		s.Logger.Debugf("pin bytes: traverse chunks: %v, addr %s", err, addr)

		if errors.Is(err, traversal.ErrInvalidType) {
			s.Logger.Error("pin bytes: invalid type")
			jsonhttp.BadRequest(w, "invalid type")
			return
		}

		s.Logger.Error("pin bytes: cannot unpin")
		jsonhttp.InternalServerError(w, "cannot unpin")
		return
	}

	jsonhttp.OK(w, nil)
}
