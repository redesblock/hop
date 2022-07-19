// Package debugapi exposes the debug API used to
// control and analyze low-level and runtime
// features and functionalities of hop.
package debugapi

import (
	"crypto/ecdsa"
	"math/big"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redesblock/hop/core/account"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/pingpong"
	"github.com/redesblock/hop/core/postage"
	"github.com/redesblock/hop/core/postage/postagecontract"
	"github.com/redesblock/hop/core/settlement"
	"github.com/redesblock/hop/core/settlement/swap"
	"github.com/redesblock/hop/core/settlement/swap/chequebook"
	"github.com/redesblock/hop/core/settlement/swap/erc20"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
	"github.com/redesblock/hop/core/topology"
	"github.com/redesblock/hop/core/topology/lightnode"
	"github.com/redesblock/hop/core/tracing"
	"github.com/redesblock/hop/core/transaction"
	"github.com/redesblock/hop/core/traversal"
	"golang.org/x/sync/semaphore"
)

type authenticator interface {
	Authorize(string) bool
	GenerateKey(string, int) (string, error)
	Enforce(string, string, string) (bool, error)
}

// Service implements http.Handler interface to be used in HTTP server.
type Service struct {
	restricted         bool
	auth               authenticator
	overlay            *swarm.Address
	publicKey          ecdsa.PublicKey
	pssPublicKey       ecdsa.PublicKey
	ethereumAddress    common.Address
	p2p                p2p.DebugService
	pingpong           pingpong.Interface
	topologyDriver     topology.Driver
	storer             storage.Storer
	tracer             *tracing.Tracer
	tags               *tags.Tags
	accounting         account.Interface
	pseudosettle       settlement.Interface
	chequebookEnabled  bool
	swapEnabled        bool
	chequebook         chequebook.Service
	swap               swap.Interface
	batchStore         postage.Storer
	transaction        transaction.Service
	chainBackend       transaction.Backend
	post               postage.Service
	postageContract    postagecontract.Interface
	logger             logging.Logger
	corsAllowedOrigins []string
	metricsRegistry    *prometheus.Registry
	lightNodes         *lightnode.Container
	blockTime          *big.Int
	traverser          traversal.Traverser
	hopMode            HopNodeMode
	gatewayMode        bool
	erc20Service       erc20.Service
	chainID            int64

	// handler is changed in the Configure method
	handler   http.Handler
	handlerMu sync.RWMutex

	// The following are semaphores which exists to limit concurrent access
	// to some parts of the resources in order to avoid undefined behaviour.
	postageSem       *semaphore.Weighted
	cashOutChequeSem *semaphore.Weighted
}

// New creates a new Debug API Service with only basic routers enabled in order
// to expose /addresses, /health endpoints, Go metrics and pprof. It is useful to expose
// these endpoints before all dependencies are configured and injected to have
// access to basic debugging tools and /health endpoint.
func New(publicKey, pssPublicKey ecdsa.PublicKey, ethereumAddress common.Address, logger logging.Logger, tracer *tracing.Tracer, corsAllowedOrigins []string, blockTime *big.Int, transaction transaction.Service, chainBackend transaction.Backend, restrict bool, auth authenticator, gatewayMode bool, hopMode HopNodeMode, chainID int64) *Service {
	s := new(Service)
	s.auth = auth
	s.restricted = restrict
	s.publicKey = publicKey
	s.pssPublicKey = pssPublicKey
	s.ethereumAddress = ethereumAddress
	s.logger = logger
	s.tracer = tracer
	s.corsAllowedOrigins = corsAllowedOrigins
	s.blockTime = blockTime
	s.metricsRegistry = newMetricsRegistry()
	s.transaction = transaction
	s.chainBackend = chainBackend
	s.postageSem = semaphore.NewWeighted(1)
	s.cashOutChequeSem = semaphore.NewWeighted(1)
	s.hopMode = hopMode
	s.gatewayMode = gatewayMode
	s.chainID = chainID

	s.setRouter(s.newBasicRouter())

	return s
}

// Configure injects required dependencies and configuration parameters and
// constructs HTTP routes that depend on them. It is intended and safe to call
// this method only once.
func (s *Service) Configure(overlay swarm.Address, p2p p2p.DebugService, pingpong pingpong.Interface, topologyDriver topology.Driver, lightNodes *lightnode.Container, storer storage.Storer, tags *tags.Tags, accounting account.Interface, pseudosettle settlement.Interface, swapEnabled bool, chequebookEnabled bool, swap swap.Interface, chequebook chequebook.Service, batchStore postage.Storer, post postage.Service, postageContract postagecontract.Interface, traverser traversal.Traverser, erc20Service erc20.Service) {
	s.p2p = p2p
	s.pingpong = pingpong
	s.topologyDriver = topologyDriver
	s.storer = storer
	s.tags = tags
	s.accounting = accounting
	s.chequebookEnabled = chequebookEnabled
	s.chequebook = chequebook
	s.swapEnabled = swapEnabled
	s.swap = swap
	s.lightNodes = lightNodes
	s.batchStore = batchStore
	s.pseudosettle = pseudosettle
	s.overlay = &overlay
	s.post = post
	s.postageContract = postageContract
	s.traverser = traverser
	s.erc20Service = erc20Service

	s.setRouter(s.newRouter())
}

// ServeHTTP implements http.Handler interface.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// protect handler as it is changed by the Configure method
	s.handlerMu.RLock()
	h := s.handler
	s.handlerMu.RUnlock()

	h.ServeHTTP(w, r)
}
