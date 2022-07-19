package debugapi_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"github.com/redesblock/hop/cmd/version"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/multiformats/go-multiaddr"
	accountingmock "github.com/redesblock/hop/core/account/mock"
	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	p2pmock "github.com/redesblock/hop/core/p2p/mock"
	"github.com/redesblock/hop/core/pingpong"
	"github.com/redesblock/hop/core/resolver"
	chequebookmock "github.com/redesblock/hop/core/settle/swap/chequebook/mock"
	swapmock "github.com/redesblock/hop/core/settle/swap/mock"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
	"github.com/redesblock/hop/core/topology/lightnode"
	topologymock "github.com/redesblock/hop/core/topology/mock"
	"github.com/redesblock/hop/core/transaction/backendmock"
	transactionmock "github.com/redesblock/hop/core/transaction/mock"
	"github.com/redesblock/hop/core/traversal"
	"github.com/redesblock/hop/core/util/logging"
	"github.com/redesblock/hop/core/voucher"
	mockpost "github.com/redesblock/hop/core/voucher/mock"
	"github.com/redesblock/hop/core/voucher/vouchercontract"
	"resenje.org/web"

	erc20mock "github.com/redesblock/hop/core/settle/swap/erc20/mock"
)

var (
	batchOk    = make([]byte, 32)
	batchOkStr string
)

func init() {
	_, _ = rand.Read(batchOk)

	batchOkStr = hex.EncodeToString(batchOk)
}

type testServerOptions struct {
	Overlay            swarm.Address
	PublicKey          ecdsa.PublicKey
	PSSPublicKey       ecdsa.PublicKey
	EthereumAddress    common.Address
	CORSAllowedOrigins []string
	P2P                *p2pmock.Service
	Pingpong           pingpong.Interface
	Storer             storage.Storer
	Resolver           resolver.Interface
	TopologyOpts       []topologymock.Option
	Tags               *tags.Tags
	AccountingOpts     []accountingmock.Option
	SettlementOpts     []swapmock.Option
	ChequebookOpts     []chequebookmock.Option
	SwapOpts           []swapmock.Option
	BatchStore         voucher.Storer
	TransactionOpts    []transactionmock.Option
	BackendOpts        []backendmock.Option
	PostageContract    vouchercontract.Interface
	Post               voucher.Service
	Traverser          traversal.Traverser
	Erc20Opts          []erc20mock.Option
	ChainID            int64
}

type testServer struct {
	Client  *http.Client
	P2PMock *p2pmock.Service
}

func newTestServer(t *testing.T, o testServerOptions) *testServer {
	topologyDriver := topologymock.NewTopologyDriver(o.TopologyOpts...)
	acc := accountingmock.NewAccounting(o.AccountingOpts...)
	settlement := swapmock.New(o.SettlementOpts...)
	chequebook := chequebookmock.NewChequebook(o.ChequebookOpts...)
	swapserv := swapmock.New(o.SwapOpts...)
	transaction := transactionmock.New(o.TransactionOpts...)
	backend := backendmock.New(o.BackendOpts...)
	erc20 := erc20mock.New(o.Erc20Opts...)
	ln := lightnode.NewContainer(o.Overlay)
	s := debugapi.New(o.PublicKey, o.PSSPublicKey, o.EthereumAddress, logging.New(io.Discard, 0), nil, o.CORSAllowedOrigins, big.NewInt(2), transaction, backend, false, nil, false, debugapi.FullMode, o.ChainID)
	s.Configure(o.Overlay, o.P2P, o.Pingpong, topologyDriver, ln, o.Storer, o.Tags, acc, settlement, true, true, swapserv, chequebook, o.BatchStore, o.Post, o.PostageContract, o.Traverser, erc20)
	ts := httptest.NewServer(s)
	t.Cleanup(ts.Close)

	client := &http.Client{
		Transport: web.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			u, err := url.Parse(ts.URL + r.URL.String())
			if err != nil {
				return nil, err
			}
			r.URL = u
			return ts.Client().Transport.RoundTrip(r)
		}),
	}
	return &testServer{
		Client:  client,
		P2PMock: o.P2P,
	}
}

func mustMultiaddr(t *testing.T, s string) multiaddr.Multiaddr {
	t.Helper()

	a, err := multiaddr.NewMultiaddr(s)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// TestServer_Configure validates that http routes are correct when server is
// constructed with only basic routes and after it is configured with
// dependencies.
func TestServer_Configure(t *testing.T) {
	privateKey, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	pssPrivateKey, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	overlay := swarm.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c")
	addresses := []multiaddr.Multiaddr{
		mustMultiaddr(t, "/ip4/127.0.0.1/tcp/7071/p2p/16Uiu2HAmTBuJT9LvNmBiQiNoTsxE5mtNy6YG3paw79m94CRa9sRb"),
		mustMultiaddr(t, "/ip4/192.168.0.101/tcp/7071/p2p/16Uiu2HAmTBuJT9LvNmBiQiNoTsxE5mtNy6YG3paw79m94CRa9sRb"),
	}

	ethereumAddress := common.HexToAddress("abcd")

	o := testServerOptions{
		PublicKey:       privateKey.PublicKey,
		PSSPublicKey:    pssPrivateKey.PublicKey,
		Overlay:         overlay,
		EthereumAddress: ethereumAddress,
		P2P: p2pmock.New(p2pmock.WithAddressesFunc(func() ([]multiaddr.Multiaddr, error) {
			return addresses, nil
		})),
	}
	topologyDriver := topologymock.NewTopologyDriver(o.TopologyOpts...)
	acc := accountingmock.NewAccounting(o.AccountingOpts...)
	settlement := swapmock.New(o.SettlementOpts...)
	chequebook := chequebookmock.NewChequebook(o.ChequebookOpts...)
	swapserv := swapmock.New(o.SwapOpts...)
	ln := lightnode.NewContainer(o.Overlay)
	transaction := transactionmock.New(o.TransactionOpts...)
	gatewayMode := false
	hopMode := debugapi.FullMode
	s := debugapi.New(o.PublicKey, o.PSSPublicKey, o.EthereumAddress, logging.New(io.Discard, 0), nil, nil, big.NewInt(2), transaction, nil, false, nil, gatewayMode, hopMode, 1)
	ts := httptest.NewServer(s)
	t.Cleanup(ts.Close)

	client := &http.Client{
		Transport: web.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			u, err := url.Parse(ts.URL + r.URL.String())
			if err != nil {
				return nil, err
			}
			r.URL = u
			return ts.Client().Transport.RoundTrip(r)
		}),
	}

	testBasicRouter(t, client)
	jsonhttptest.Request(t, client, http.MethodGet, "/readiness", http.StatusNotFound,
		jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
			Message: http.StatusText(http.StatusNotFound),
			Code:    http.StatusNotFound,
		}),
	)
	jsonhttptest.Request(t, client, http.MethodGet, "/addresses", http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(debugapi.AddressesResponse{
			Underlay:     make([]multiaddr.Multiaddr, 0),
			Ethereum:     o.EthereumAddress,
			PublicKey:    hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&o.PublicKey)),
			PSSPublicKey: hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&o.PSSPublicKey)),
		}),
	)

	s.Configure(o.Overlay, o.P2P, o.Pingpong, topologyDriver, ln, o.Storer, o.Tags, acc, settlement, true, true, swapserv, chequebook, nil, mockpost.New(), nil, nil, nil)

	testBasicRouter(t, client)
	jsonhttptest.Request(t, client, http.MethodGet, "/readiness", http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(debugapi.StatusResponse{
			Status:          "ok",
			Version:         version.Version,
			APIVersion:      api.Version,
			DebugAPIVersion: debugapi.Version,
		}),
	)
	jsonhttptest.Request(t, client, http.MethodGet, "/addresses", http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(debugapi.AddressesResponse{
			Overlay:      &o.Overlay,
			Underlay:     addresses,
			Ethereum:     o.EthereumAddress,
			PublicKey:    hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&o.PublicKey)),
			PSSPublicKey: hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&o.PSSPublicKey)),
		}),
	)
	jsonhttptest.Request(t, client, http.MethodGet, "/node", http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(debugapi.NodeResponse{
			HopMode:           hopMode.String(),
			GatewayMode:       gatewayMode,
			ChequebookEnabled: true,
			SwapEnabled:       true,
		}),
	)
}

func testBasicRouter(t *testing.T, client *http.Client) {
	t.Helper()

	jsonhttptest.Request(t, client, http.MethodGet, "/health", http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(debugapi.StatusResponse{
			Status:          "ok",
			Version:         version.Version,
			APIVersion:      api.Version,
			DebugAPIVersion: debugapi.Version,
		}),
	)

	for _, path := range []string{
		"/metrics",
		"/debug/pprof",
		"/debug/pprof/cmdline",
		"/debug/pprof/profile?seconds=1", // profile for only 1 second to check only the status code
		"/debug/pprof/symbol",
		"/debug/pprof/trace",
		"/debug/vars",
	} {
		jsonhttptest.Request(t, client, http.MethodGet, path, http.StatusOK)
	}
}
