package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
)

func (s *Service) hasChunkHandler(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.logger.Debugf("parse chunk address: %v", err)
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.storer.Has(r.Context(), addr)
	if err != nil {
		s.logger.Debugf("chunk address: localstore has: %v", err)
		jsonhttp.BadRequest(w, err)
		return
	}

	if !has {
		jsonhttp.NotFound(w, nil)
		return
	}
	jsonhttp.OK(w, nil)
}

func (s *Service) removeChunk(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.logger.Debugf("parse chunk address: %v", err)
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.storer.Has(r.Context(), addr)
	if err != nil {
		s.logger.Debugf("chunk address: localstore remove: %v", err)
		jsonhttp.BadRequest(w, err)
		return
	}

	if !has {
		jsonhttp.OK(w, nil)
		return
	}

	err = s.storer.Set(r.Context(), storage.ModeSetRemove, addr)
	if err != nil {
		s.logger.Debugf("chunk address: localstore remove: %v", err)
		jsonhttp.InternalServerError(w, err)
		return
	}
	jsonhttp.OK(w, nil)
}
