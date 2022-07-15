package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redesblock/hop/core/bigint"
	"github.com/redesblock/hop/core/jsonhttp"
)

type walletResponse struct {
	HOP             *bigint.BigInt `json:"hop"`             // the HOP balance of the wallet associated with the eth address of the node
	XDai            *bigint.BigInt `json:"xDai"`            // the xDai balance of the wallet associated with the eth address of the node
	ChainID         int64          `json:"chainID"`         // the id of the block chain
	ContractAddress common.Address `json:"contractAddress"` // the address of the chequebook contract
}

func (s *Service) walletHandler(w http.ResponseWriter, r *http.Request) {

	xdai, err := s.chainBackend.BalanceAt(r.Context(), s.ethereumAddress, nil)
	if err != nil {
		s.logger.Debugf("wallet: unable to acquire balance from the chain backend: %v", err)
		s.logger.Error("wallet: unable to acquire balance from the chain backend")
		jsonhttp.InternalServerError(w, "unable to acquire balance from the chain backend")
		return
	}

	hop, err := s.erc20Service.BalanceOf(r.Context(), s.ethereumAddress)
	if err != nil {
		s.logger.Debugf("wallet: unable to acquire erc20 balance: %v", err)
		s.logger.Error("wallet: unable to acquire erc20 balance")
		jsonhttp.InternalServerError(w, "unable to acquire erc20 balance")
		return
	}

	jsonhttp.OK(w, walletResponse{
		HOP:             bigint.Wrap(hop),
		XDai:            bigint.Wrap(xdai),
		ChainID:         s.chainID,
		ContractAddress: s.chequebook.Address(),
	})
}
