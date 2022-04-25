package debugapi_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/p2p/mock"
	"github.com/redesblock/hop/core/swarm"
)

func TestAddresses(t *testing.T) {
	overlay := swarm.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c")
	addresses := []multiaddr.Multiaddr{
		mustMultiaddr(t, "/ip4/127.0.0.1/tcp/7071/p2p/16Uiu2HAmTBuJT9LvNmBiQiNoTsxE5mtNy6YG3paw79m94CRa9sRb"),
		mustMultiaddr(t, "/ip4/192.168.0.101/tcp/7071/p2p/16Uiu2HAmTBuJT9LvNmBiQiNoTsxE5mtNy6YG3paw79m94CRa9sRb"),
		mustMultiaddr(t, "/ip4/127.0.0.1/udp/7071/quic/p2p/16Uiu2HAmTBuJT9LvNmBiQiNoTsxE5mtNy6YG3paw79m94CRa9sRb"),
	}

	testServer := newTestServer(t, testServerOptions{
		Overlay: overlay,
		P2P: mock.New(mock.WithAddressesFunc(func() ([]multiaddr.Multiaddr, error) {
			return addresses, nil
		})),
	})
	defer testServer.Cleanup()

	t.Run("ok", func(t *testing.T) {
		jsonhttptest.ResponseDirect(t, testServer.Client, http.MethodGet, "/addresses", nil, http.StatusOK, debugapi.AddressesResponse{
			Overlay:  overlay,
			Underlay: addresses,
		})
	})

	t.Run("post method not allowed", func(t *testing.T) {
		jsonhttptest.ResponseDirect(t, testServer.Client, http.MethodPost, "/addresses", nil, http.StatusMethodNotAllowed, jsonhttp.StatusResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
	})
}

func TestAddresses_error(t *testing.T) {
	testErr := errors.New("test error")

	testServer := newTestServer(t, testServerOptions{
		P2P: mock.New(mock.WithAddressesFunc(func() ([]multiaddr.Multiaddr, error) {
			return nil, testErr
		})),
	})
	defer testServer.Cleanup()

	jsonhttptest.ResponseDirect(t, testServer.Client, http.MethodGet, "/addresses", nil, http.StatusInternalServerError, jsonhttp.StatusResponse{
		Code:    http.StatusInternalServerError,
		Message: testErr.Error(),
	})
}