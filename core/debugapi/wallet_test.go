package debugapi_test

import (
	"context"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	erc20mock "github.com/redesblock/hop/core/settlement/swap/erc20/mock"
	"github.com/redesblock/hop/core/transaction/backendmock"
	"github.com/redesblock/hop/core/util/bigint"
)

func TestWallet(t *testing.T) {

	t.Run("Okay", func(t *testing.T) {

		srv := newTestServer(t, testServerOptions{
			Erc20Opts: []erc20mock.Option{
				erc20mock.WithBalanceOfFunc(func(ctx context.Context, address common.Address) (*big.Int, error) {
					return big.NewInt(10000000000000000), nil
				}),
			},
			BackendOpts: []backendmock.Option{
				backendmock.WithBalanceAt(func(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {
					return big.NewInt(2000000000000000000), nil
				}),
			},
			ChainID: 1,
		})

		jsonhttptest.Request(t, srv.Client, http.MethodGet, "/wallet", http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(debugapi.WalletResponse{
				HOP:     bigint.Wrap(big.NewInt(10000000000000000)),
				BSC:     bigint.Wrap(big.NewInt(2000000000000000000)),
				ChainID: 1,
			}),
		)
	})

	t.Run("500 - erc20 error", func(t *testing.T) {
		srv := newTestServer(t, testServerOptions{
			BackendOpts: []backendmock.Option{
				backendmock.WithBalanceAt(func(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {
					return new(big.Int), nil
				}),
			},
		})

		jsonhttptest.Request(t, srv.Client, http.MethodGet, "/wallet", http.StatusInternalServerError,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "unable to acquire erc20 balance",
				Code:    500,
			}))
	})

	t.Run("500 - chain backend error", func(t *testing.T) {
		srv := newTestServer(t, testServerOptions{
			Erc20Opts: []erc20mock.Option{
				erc20mock.WithBalanceOfFunc(func(ctx context.Context, address common.Address) (*big.Int, error) {
					return new(big.Int), nil
				}),
			},
		})

		jsonhttptest.Request(t, srv.Client, http.MethodGet, "/wallet", http.StatusInternalServerError,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "unable to acquire balance from the chain backend",
				Code:    500,
			}))
	})
}
