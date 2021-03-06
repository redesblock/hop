package debugapi

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/postage/postagecontract"
	"github.com/redesblock/hop/core/settlement"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/util/bigint"
)

var (
	errCantSettlements     = "can not get settlements"
	errCantSettlementsPeer = "can not get settlements for peer"
)

type settlementResponse struct {
	Peer               string         `json:"peer"`
	SettlementReceived *bigint.BigInt `json:"received"`
	SettlementSent     *bigint.BigInt `json:"sent"`
}

type settlementsResponse struct {
	TotalSettlementReceived *bigint.BigInt       `json:"totalReceived"`
	TotalSettlementSent     *bigint.BigInt       `json:"totalSent"`
	Settlements             []settlementResponse `json:"settlements"`
}

func (s *Service) settlementsHandler(w http.ResponseWriter, r *http.Request) {

	settlementsSent, err := s.swap.SettlementsSent()
	if errors.Is(err, postagecontract.ErrChainDisabled) {
		s.logger.Debugf("debug api: sent settlements: %v", err)
		s.logger.Error("debug api: can not get sent settlements")
		jsonhttp.MethodNotAllowed(w, err)
		return
	}
	if err != nil {
		s.logger.Debugf("debug api: sent settlements: %v", err)
		s.logger.Error("debug api: can not get sent settlements")
		jsonhttp.InternalServerError(w, errCantSettlements)
		return
	}
	settlementsReceived, err := s.swap.SettlementsReceived()
	if err != nil {
		s.logger.Debugf("debug api: received settlements: %v", err)
		s.logger.Error("debug api: can not get received settlements")
		jsonhttp.InternalServerError(w, errCantSettlements)
		return
	}

	totalReceived := big.NewInt(0)
	totalSent := big.NewInt(0)

	settlementResponses := make(map[string]settlementResponse)

	for a, b := range settlementsSent {
		settlementResponses[a] = settlementResponse{
			Peer:               a,
			SettlementSent:     bigint.Wrap(b),
			SettlementReceived: bigint.Wrap(big.NewInt(0)),
		}
		totalSent.Add(b, totalSent)
	}

	for a, b := range settlementsReceived {
		if _, ok := settlementResponses[a]; ok {
			t := settlementResponses[a]
			t.SettlementReceived = bigint.Wrap(b)
			settlementResponses[a] = t
		} else {
			settlementResponses[a] = settlementResponse{
				Peer:               a,
				SettlementSent:     bigint.Wrap(big.NewInt(0)),
				SettlementReceived: bigint.Wrap(b),
			}
		}
		totalReceived.Add(b, totalReceived)
	}

	settlementResponsesArray := make([]settlementResponse, len(settlementResponses))
	i := 0
	for k := range settlementResponses {
		settlementResponsesArray[i] = settlementResponses[k]
		i++
	}

	jsonhttp.OK(w, settlementsResponse{TotalSettlementReceived: bigint.Wrap(totalReceived), TotalSettlementSent: bigint.Wrap(totalSent), Settlements: settlementResponsesArray})
}

func (s *Service) peerSettlementsHandler(w http.ResponseWriter, r *http.Request) {
	addr := mux.Vars(r)["peer"]
	peer, err := swarm.ParseHexAddress(addr)
	if err != nil {
		s.logger.Debugf("debug api: settlements peer: invalid peer address %s: %v", addr, err)
		s.logger.Errorf("debug api: settlements peer: invalid peer address %s", addr)
		jsonhttp.NotFound(w, errInvalidAddress)
		return
	}

	peerexists := false

	received, err := s.swap.TotalReceived(peer)
	if errors.Is(err, postagecontract.ErrChainDisabled) {
		s.logger.Debugf("debug api: settlements peer: %v", err)
		s.logger.Errorf("debug api: settlements peer: can't get peer %s received settlement", peer.String())
		jsonhttp.MethodNotAllowed(w, err)
		return
	}
	if err != nil {
		if !errors.Is(err, settlement.ErrPeerNoSettlements) {
			s.logger.Debugf("debug api: settlements peer: get peer %s received settlement: %v", peer.String(), err)
			s.logger.Errorf("debug api: settlements peer: can't get peer %s received settlement", peer.String())
			jsonhttp.InternalServerError(w, errCantSettlementsPeer)
			return
		} else {
			received = big.NewInt(0)
		}
	}

	if err == nil {
		peerexists = true
	}

	sent, err := s.swap.TotalSent(peer)
	if err != nil {
		if !errors.Is(err, settlement.ErrPeerNoSettlements) {
			s.logger.Debugf("debug api: settlements peer: get peer %s sent settlement: %v", peer.String(), err)
			s.logger.Errorf("debug api: settlements peer: can't get peer %s sent settlement", peer.String())
			jsonhttp.InternalServerError(w, errCantSettlementsPeer)
			return
		} else {
			sent = big.NewInt(0)
		}
	}

	if err == nil {
		peerexists = true
	}

	if !peerexists {
		jsonhttp.NotFound(w, settlement.ErrPeerNoSettlements)
		return
	}

	jsonhttp.OK(w, settlementResponse{
		Peer:               peer.String(),
		SettlementReceived: bigint.Wrap(received),
		SettlementSent:     bigint.Wrap(sent),
	})
}

func (s *Service) settlementsHandlerPseudosettle(w http.ResponseWriter, r *http.Request) {

	settlementsSent, err := s.pseudosettle.SettlementsSent()
	if err != nil {
		jsonhttp.InternalServerError(w, errCantSettlements)
		s.logger.Debugf("debug api: sent settlements: %v", err)
		s.logger.Error("debug api: can not get sent settlements")
		return
	}
	settlementsReceived, err := s.pseudosettle.SettlementsReceived()
	if err != nil {
		jsonhttp.InternalServerError(w, errCantSettlements)
		s.logger.Debugf("debug api: received settlements: %v", err)
		s.logger.Error("debug api: can not get received settlements")
		return
	}

	totalReceived := big.NewInt(0)
	totalSent := big.NewInt(0)

	settlementResponses := make(map[string]settlementResponse)

	for a, b := range settlementsSent {
		settlementResponses[a] = settlementResponse{
			Peer:               a,
			SettlementSent:     bigint.Wrap(b),
			SettlementReceived: bigint.Wrap(big.NewInt(0)),
		}
		totalSent.Add(b, totalSent)
	}

	for a, b := range settlementsReceived {
		if _, ok := settlementResponses[a]; ok {
			t := settlementResponses[a]
			t.SettlementReceived = bigint.Wrap(b)
			settlementResponses[a] = t
		} else {
			settlementResponses[a] = settlementResponse{
				Peer:               a,
				SettlementSent:     bigint.Wrap(big.NewInt(0)),
				SettlementReceived: bigint.Wrap(b),
			}
		}
		totalReceived.Add(b, totalReceived)
	}

	settlementResponsesArray := make([]settlementResponse, len(settlementResponses))
	i := 0
	for k := range settlementResponses {
		settlementResponsesArray[i] = settlementResponses[k]
		i++
	}

	jsonhttp.OK(w, settlementsResponse{TotalSettlementReceived: bigint.Wrap(totalReceived), TotalSettlementSent: bigint.Wrap(totalSent), Settlements: settlementResponsesArray})
}
