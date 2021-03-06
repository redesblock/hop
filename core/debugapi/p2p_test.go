package debugapi_test

import (
	"encoding/hex"
	"errors"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/multiformats/go-multiaddr"
	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/p2p/mock"
	"github.com/redesblock/hop/core/swarm"
)

func TestAddresses(t *testing.T) {
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

	testServer := newTestServer(t, testServerOptions{
		PublicKey:       privateKey.PublicKey,
		PSSPublicKey:    pssPrivateKey.PublicKey,
		Overlay:         overlay,
		EthereumAddress: ethereumAddress,
		P2P: mock.New(mock.WithAddressesFunc(func() ([]multiaddr.Multiaddr, error) {
			return addresses, nil
		})),
	})

	t.Run("ok", func(t *testing.T) {
		jsonhttptest.Request(t, testServer.Client, http.MethodGet, "/addresses", http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(debugapi.AddressesResponse{
				Overlay:      &overlay,
				Underlay:     addresses,
				Ethereum:     ethereumAddress,
				PublicKey:    hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&privateKey.PublicKey)),
				PSSPublicKey: hex.EncodeToString(crypto.EncodeSecp256k1PublicKey(&pssPrivateKey.PublicKey)),
			}),
		)
	})

	t.Run("post method not allowed", func(t *testing.T) {
		jsonhttptest.Request(t, testServer.Client, http.MethodPost, "/addresses", http.StatusMethodNotAllowed,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: http.StatusText(http.StatusMethodNotAllowed),
			}),
		)
	})
}

func TestAddresses_error(t *testing.T) {
	testErr := errors.New("test error")

	testServer := newTestServer(t, testServerOptions{
		P2P: mock.New(mock.WithAddressesFunc(func() ([]multiaddr.Multiaddr, error) {
			return nil, testErr
		})),
	})

	jsonhttptest.Request(t, testServer.Client, http.MethodGet, "/addresses", http.StatusInternalServerError,
		jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
			Code:    http.StatusInternalServerError,
			Message: testErr.Error(),
		}),
	)
}
