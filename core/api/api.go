package api

import (
	"net/http"

	"github.com/redesblock/hop/core/logging"
	m "github.com/redesblock/hop/core/metrics"
	"github.com/redesblock/hop/core/pingpong"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/tracing"
)

type Service interface {
	http.Handler
	m.Collector
}

type server struct {
	Options
	http.Handler
	metrics metrics
}

type Options struct {
	Pingpong pingpong.Interface
	Storer   storage.Storer
	Logger   logging.Logger
	Tracer   *tracing.Tracer
}

func New(o Options) Service {
	s := &server{
		Options: o,
		metrics: newMetrics(),
	}

	s.setupRouting()

	return s
}