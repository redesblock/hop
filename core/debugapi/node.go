package debugapi

import (
	"net/http"

	"github.com/redesblock/hop/core/jsonhttp"
)

type HopNodeMode uint

const (
	LightMode HopNodeMode = iota
	FullMode
	DevMode
	UltraLightMode
)

type nodeResponse struct {
	HopMode           string `json:"hopMode"`
	GatewayMode       bool   `json:"gatewayMode"`
	ChequebookEnabled bool   `json:"chequebookEnabled"`
	SwapEnabled       bool   `json:"swapEnabled"`
}

func (b HopNodeMode) String() string {
	switch b {
	case LightMode:
		return "light"
	case FullMode:
		return "full"
	case DevMode:
		return "dev"
	case UltraLightMode:
		return "ultra-light"
	}
	return "unknown"
}

// nodeGetHandler gives back information about the node configuration.
func (s *Service) nodeGetHandler(w http.ResponseWriter, r *http.Request) {
	jsonhttp.OK(w, nodeResponse{
		HopMode:           s.hopMode.String(),
		GatewayMode:       s.gatewayMode,
		ChequebookEnabled: s.chequebookEnabled,
		SwapEnabled:       s.swapEnabled,
	})
}
