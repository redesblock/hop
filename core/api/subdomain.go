package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/tracing"
)

func (s *Service) subdomainHandler(w http.ResponseWriter, r *http.Request) {
	logger := tracing.NewLoggerWithTraceID(r.Context(), s.logger)

	nameOrHex := mux.Vars(r)["subdomain"]
	pathVar := mux.Vars(r)["path"]
	if strings.HasSuffix(pathVar, "/") {
		pathVar = strings.TrimRight(pathVar, "/")
		// NOTE: leave one slash if there was some
		pathVar += "/"
	}

	address, err := s.resolveNameOrAddress(nameOrHex)
	if err != nil {
		logger.Debugf("subdomain handler: parse address %s: %v", nameOrHex, err)
		logger.Error("subdomain handler: parse address")
		jsonhttp.NotFound(w, nil)
		return
	}

	s.serveReference(address, pathVar, w, r)
}
