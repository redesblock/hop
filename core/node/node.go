package node

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/redesblock/hop/core/accounting"
	"github.com/redesblock/hop/core/addressbook"
	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/content"
	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/hive"
	"github.com/redesblock/hop/core/kademlia"
	"github.com/redesblock/hop/core/keystore"
	filekeystore "github.com/redesblock/hop/core/keystore/file"
	memkeystore "github.com/redesblock/hop/core/keystore/mem"
	"github.com/redesblock/hop/core/localstore"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/metrics"
	"github.com/redesblock/hop/core/netstore"
	"github.com/redesblock/hop/core/p2p/libp2p"
	"github.com/redesblock/hop/core/pingpong"
	"github.com/redesblock/hop/core/puller"
	"github.com/redesblock/hop/core/pullsync"
	"github.com/redesblock/hop/core/pullsync/pullstorage"
	"github.com/redesblock/hop/core/pusher"
	"github.com/redesblock/hop/core/pushsync"
	"github.com/redesblock/hop/core/retrieval"
	"github.com/redesblock/hop/core/settlement/pseudosettle"
	"github.com/redesblock/hop/core/soc"
	"github.com/redesblock/hop/core/statestore/leveldb"
	mockinmem "github.com/redesblock/hop/core/statestore/mock"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
	"github.com/redesblock/hop/core/tracing"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Node struct {
	p2pService       io.Closer
	p2pCancel        context.CancelFunc
	apiServer        *http.Server
	debugAPIServer   *http.Server
	errorLogWriter   *io.PipeWriter
	tracerCloser     io.Closer
	stateStoreCloser io.Closer
	localstoreCloser io.Closer
	topologyCloser   io.Closer
	pusherCloser     io.Closer
	pullerCloser     io.Closer
	pullSyncCloser   io.Closer
}

type Options struct {
	DataDir             string
	DBCapacity          uint64
	Password            string
	APIAddr             string
	DebugAPIAddr        string
	NATAddr             string
	EnableWS            bool
	EnableQUIC          bool
	NetworkID           uint64
	WelcomeMessage      string
	Bootnodes           []string
	CORSAllowedOrigins  []string
	TracingEnabled      bool
	TracingEndpoint     string
	TracingServiceName  string
	DisconnectThreshold uint64
	PaymentThreshold    uint64
	PaymentTolerance    uint64
}

func New(addr string, logger logging.Logger, o Options) (*Node, error) {
	tracer, tracerCloser, err := tracing.NewTracer(&tracing.Options{
		Enabled:     o.TracingEnabled,
		Endpoint:    o.TracingEndpoint,
		ServiceName: o.TracingServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("tracer: %w", err)
	}

	p2pCtx, p2pCancel := context.WithCancel(context.Background())

	b := &Node{
		p2pCancel:      p2pCancel,
		errorLogWriter: logger.WriterLevel(logrus.ErrorLevel),
		tracerCloser:   tracerCloser,
	}

	var keyStore keystore.Service
	if o.DataDir == "" {
		keyStore = memkeystore.New()
		logger.Warning("data directory not provided, keys are not persisted")
	} else {
		keyStore = filekeystore.New(filepath.Join(o.DataDir, "keys"))
	}

	swarmPrivateKey, created, err := keyStore.Key("swarm", o.Password)
	if err != nil {
		return nil, fmt.Errorf("swarm key: %w", err)
	}
	address, err := crypto.NewOverlayAddress(swarmPrivateKey.PublicKey, o.NetworkID)
	if err != nil {
		return nil, err
	}
	if created {
		logger.Infof("new swarm network address created: %s", address)
	} else {
		logger.Infof("using existing swarm network address: %s", address)
	}

	// Construct P2P service.
	libp2pPrivateKey, created, err := keyStore.Key("libp2p", o.Password)
	if err != nil {
		return nil, fmt.Errorf("libp2p key: %w", err)
	}
	if created {
		logger.Debugf("new libp2p key created")
	} else {
		logger.Debugf("using existing libp2p key")
	}

	var stateStore storage.StateStorer
	if o.DataDir == "" {
		stateStore = mockinmem.NewStateStore()
		logger.Warning("using in-mem state store. no node state will be persisted")
	} else {
		stateStore, err = leveldb.NewStateStore(filepath.Join(o.DataDir, "statestore"))
		if err != nil {
			return nil, fmt.Errorf("statestore: %w", err)
		}
	}
	b.stateStoreCloser = stateStore
	addressbook := addressbook.New(stateStore)
	signer := crypto.NewDefaultSigner(swarmPrivateKey)

	p2ps, err := libp2p.New(p2pCtx, signer, o.NetworkID, address, addr, addressbook, logger, tracer, libp2p.Options{
		PrivateKey:     libp2pPrivateKey,
		NATAddr:        o.NATAddr,
		EnableWS:       o.EnableWS,
		EnableQUIC:     o.EnableQUIC,
		WelcomeMessage: o.WelcomeMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("p2p service: %w", err)
	}
	b.p2pService = p2ps

	if natManager := p2ps.NATManager(); natManager != nil {
		// wait for nat manager to init
		logger.Debug("initializing NAT manager")
		select {
		case <-natManager.Ready():
			// this is magic sleep to give NAT time to sync the mappings
			// this is a hack, kind of alchemy and should be improved
			time.Sleep(3 * time.Second)
			logger.Debug("NAT manager initialized")
		case <-time.After(10 * time.Second):
			logger.Warning("NAT manager init timeout")
		}
	}

	// Construct protocols.
	pingPong := pingpong.New(pingpong.Options{
		Streamer: p2ps,
		Logger:   logger,
		Tracer:   tracer,
	})

	if err = p2ps.AddProtocol(pingPong.Protocol()); err != nil {
		return nil, fmt.Errorf("pingpong service: %w", err)
	}

	hive := hive.New(hive.Options{
		Streamer:    p2ps,
		AddressBook: addressbook,
		NetworkID:   o.NetworkID,
		Logger:      logger,
	})

	if err = p2ps.AddProtocol(hive.Protocol()); err != nil {
		return nil, fmt.Errorf("hive service: %w", err)
	}

	var bootnodes []ma.Multiaddr
	for _, a := range o.Bootnodes {
		addr, err := ma.NewMultiaddr(a)
		if err != nil {
			logger.Debugf("multiaddress fail %s: %v", a, err)
			logger.Warningf("invalid bootnode address %s", a)
			continue
		}

		bootnodes = append(bootnodes, addr)
	}

	kad := kademlia.New(kademlia.Options{Base: address, Discovery: hive, AddressBook: addressbook, P2P: p2ps, Bootnodes: bootnodes, Logger: logger})
	b.topologyCloser = kad
	hive.SetAddPeersHandler(kad.AddPeers)
	p2ps.AddNotifier(kad)
	addrs, err := p2ps.Addresses()
	if err != nil {
		return nil, fmt.Errorf("get server addresses: %w", err)
	}

	for _, addr := range addrs {
		logger.Debugf("p2p address: %s", addr)
	}

	var path string

	if o.DataDir != "" {
		path = filepath.Join(o.DataDir, "localstore")
	}
	lo := &localstore.Options{
		Capacity: o.DBCapacity,
	}
	storer, err := localstore.New(path, address.Bytes(), lo, logger)
	if err != nil {
		return nil, fmt.Errorf("localstore: %w", err)
	}
	b.localstoreCloser = storer

	settlement := pseudosettle.New(pseudosettle.Options{
		Streamer: p2ps,
		Logger:   logger,
	})

	if err = p2ps.AddProtocol(settlement.Protocol()); err != nil {
		return nil, fmt.Errorf("pseudosettle service: %w", err)
	}

	acc, err := accounting.NewAccounting(accounting.Options{
		Logger:           logger,
		Store:            stateStore,
		PaymentThreshold: o.PaymentThreshold,
		PaymentTolerance: o.PaymentTolerance,
		Settlement:       settlement,
	})
	if err != nil {
		return nil, fmt.Errorf("accounting: %w", err)
	}

	settlement.SetPaymentObserver(acc)

	chunkvalidator := swarm.NewChunkValidator(soc.NewValidator(), content.NewValidator())

	retrieve := retrieval.New(retrieval.Options{
		Streamer:    p2ps,
		ChunkPeerer: kad,
		Logger:      logger,
		Accounting:  acc,
		Pricer:      accounting.NewFixedPricer(address, 10),
		Validator:   chunkvalidator,
	})
	tagg := tags.NewTags()

	if err = p2ps.AddProtocol(retrieve.Protocol()); err != nil {
		return nil, fmt.Errorf("retrieval service: %w", err)
	}

	ns := netstore.New(storer, retrieve, logger, chunkvalidator)

	retrieve.SetStorer(ns)

	pushSyncProtocol := pushsync.New(pushsync.Options{
		Streamer:      p2ps,
		Storer:        storer,
		ClosestPeerer: kad,
		Tagger:        tagg,
		Logger:        logger,
	})

	if err = p2ps.AddProtocol(pushSyncProtocol.Protocol()); err != nil {
		return nil, fmt.Errorf("pushsync service: %w", err)
	}

	pushSyncPusher := pusher.New(pusher.Options{
		Storer:        storer,
		PeerSuggester: kad,
		PushSyncer:    pushSyncProtocol,
		Tagger:        tagg,
		Logger:        logger,
	})
	b.pusherCloser = pushSyncPusher

	pullStorage := pullstorage.New(storer)

	pullSync := pullsync.New(pullsync.Options{
		Streamer: p2ps,
		Storage:  pullStorage,
		Logger:   logger,
	})
	b.pullSyncCloser = pullSync

	if err = p2ps.AddProtocol(pullSync.Protocol()); err != nil {
		return nil, fmt.Errorf("pullsync protocol: %w", err)
	}

	puller := puller.New(puller.Options{
		StateStore: stateStore,
		Topology:   kad,
		PullSync:   pullSync,
		Logger:     logger,
	})

	b.pullerCloser = puller

	var apiService api.Service
	if o.APIAddr != "" {
		// API server
		apiService = api.New(tagg, ns, o.CORSAllowedOrigins, logger, tracer)
		apiListener, err := net.Listen("tcp", o.APIAddr)
		if err != nil {
			return nil, fmt.Errorf("api listener: %w", err)
		}

		apiServer := &http.Server{
			Handler:  apiService,
			ErrorLog: log.New(b.errorLogWriter, "", 0),
		}

		go func() {
			logger.Infof("api address: %s", apiListener.Addr())

			if err := apiServer.Serve(apiListener); err != nil && err != http.ErrServerClosed {
				logger.Debugf("api server: %v", err)
				logger.Error("unable to serve api")
			}
		}()

		b.apiServer = apiServer
	}

	if o.DebugAPIAddr != "" {
		// Debug API server
		debugAPIService := debugapi.New(debugapi.Options{
			Overlay:        address,
			P2P:            p2ps,
			Pingpong:       pingPong,
			Logger:         logger,
			Tracer:         tracer,
			TopologyDriver: kad,
			Storer:         storer,
			Tags:           tagg,
			Accounting:     acc,
		})
		// register metrics from components
		debugAPIService.MustRegisterMetrics(p2ps.Metrics()...)
		debugAPIService.MustRegisterMetrics(pingPong.Metrics()...)
		debugAPIService.MustRegisterMetrics(acc.Metrics()...)

		if apiService != nil {
			debugAPIService.MustRegisterMetrics(apiService.Metrics()...)
		}
		if l, ok := logger.(metrics.Collector); ok {
			debugAPIService.MustRegisterMetrics(l.Metrics()...)
		}

		debugAPIListener, err := net.Listen("tcp", o.DebugAPIAddr)
		if err != nil {
			return nil, fmt.Errorf("debug api listener: %w", err)
		}

		debugAPIServer := &http.Server{
			Handler:  debugAPIService,
			ErrorLog: log.New(b.errorLogWriter, "", 0),
		}

		go func() {
			logger.Infof("debug api address: %s", debugAPIListener.Addr())

			if err := debugAPIServer.Serve(debugAPIListener); err != nil && err != http.ErrServerClosed {
				logger.Debugf("debug api server: %v", err)
				logger.Error("unable to serve debug api")
			}
		}()

		b.debugAPIServer = debugAPIServer
	}

	if err := kad.Start(p2pCtx); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Node) Shutdown(ctx context.Context) error {
	errs := new(multiError)

	var eg errgroup.Group
	if b.apiServer != nil {
		eg.Go(func() error {
			if err := b.apiServer.Shutdown(ctx); err != nil {
				return fmt.Errorf("api server: %w", err)
			}
			return nil
		})
	}
	if b.debugAPIServer != nil {
		eg.Go(func() error {
			if err := b.debugAPIServer.Shutdown(ctx); err != nil {
				return fmt.Errorf("debug api server: %w", err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		errs.add(err)
	}

	if err := b.pusherCloser.Close(); err != nil {
		errs.add(fmt.Errorf("pusher: %w", err))
	}

	if err := b.pullerCloser.Close(); err != nil {
		errs.add(fmt.Errorf("puller: %w", err))
	}

	if err := b.pullSyncCloser.Close(); err != nil {
		errs.add(fmt.Errorf("pull sync: %w", err))
	}

	b.p2pCancel()
	if err := b.p2pService.Close(); err != nil {
		errs.add(fmt.Errorf("p2p server: %w", err))
	}

	if err := b.tracerCloser.Close(); err != nil {
		errs.add(fmt.Errorf("tracer: %w", err))
	}

	if err := b.stateStoreCloser.Close(); err != nil {
		errs.add(fmt.Errorf("statestore: %w", err))
	}

	if err := b.localstoreCloser.Close(); err != nil {
		errs.add(fmt.Errorf("localstore: %w", err))
	}

	if err := b.topologyCloser.Close(); err != nil {
		errs.add(fmt.Errorf("topology driver: %w", err))
	}

	if err := b.errorLogWriter.Close(); err != nil {
		errs.add(fmt.Errorf("error log writer: %w", err))
	}

	if errs.hasErrors() {
		return errs
	}

	return nil
}

type multiError struct {
	errors []error
}

func (e *multiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	s := e.errors[0].Error()
	for _, err := range e.errors[1:] {
		s += "; " + err.Error()
	}
	return s
}

func (e *multiError) add(err error) {
	e.errors = append(e.errors, err)
}

func (e *multiError) hasErrors() bool {
	return len(e.errors) > 0
}
