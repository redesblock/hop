package debugapi_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/redesblock/hop/core/debugapi"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/p2p/mock"
)

func TestGetWelcomeMessage(t *testing.T) {
	const DefaultTestWelcomeMessage = "Hello World!"

	srv := newTestServer(t, testServerOptions{
		P2P: mock.New(mock.WithGetWelcomeMessageFunc(func() string {
			return DefaultTestWelcomeMessage
		}))})

	jsonhttptest.ResponseDirect(t, srv.Client, http.MethodGet, "/welcome-message", nil, http.StatusOK, debugapi.WelcomeMessageResponse{
		WelcomeMesssage: DefaultTestWelcomeMessage,
	})
}

func TestSetWelcomeMessage(t *testing.T) {
	testCases := []struct {
		desc        string
		message     string
		wantFail    bool
		wantStatus  int
		wantMessage string
	}{
		{
			desc:       "OK",
			message:    "Changed value",
			wantStatus: http.StatusOK,
		},
		{
			desc:       "OK - welcome message empty",
			message:    "",
			wantStatus: http.StatusOK,
		},
		{
			desc:     "fails - request entity too large",
			wantFail: true,
			message: `pPPopopopoHpHHPoHoPoooHpppPHPHoppHHHoHpHppPooHpHHpHHoPp
oPPPHHooPooPpHopHopoPHPpHPPoPpPpPpooPPHPppoHPHpPppHHpPPppPoPPPpooopp
oHpPPHoHPHpPpHPHpopHHopHHpopppHoHoPpPHPHPpHPPooPPHPPHpPpHPopPHpPoHpP
oooHooPpPopoPpPpopppoppopPPpoopoHPPoHoHPHHPpPppoHHHHHPoPpHppHopHooop
HHpooPHopopHPoppHpoPHppopoooHPHpHPpHPopHpPpHPPPHpPPHpHPPpopPoppPpHHp
PPpoPppPPPHPHHoPPoPpHHHpopPPooPPHPPHHHoHPpPoPHPHHHppPHoooopHpoopHHHp
oPHHoPpHoPPHpoHoPPHpHpHpHHopppPHopoPHopHoPpooHPHHooPoHHooHPopoPpoHpH
oHooPPopPpooHopPoPPPPoppPPoHpPPoPpPppHpoPP`, // 513 characters
			wantStatus: http.StatusRequestEntityTooLarge,
		},
	}
	testURL := "/welcome-message"

	srv := newTestServer(t, testServerOptions{
		P2P: mock.New(),
	})

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.wantMessage == "" {
				tC.wantMessage = http.StatusText(tC.wantStatus)
			}
			data, _ := json.Marshal(debugapi.WelcomeMessageRequest{
				WelcomeMesssage: tC.message,
			})
			body := bytes.NewReader(data)
			wantResponse := jsonhttp.StatusResponse{
				Message: tC.wantMessage,
				Code:    tC.wantStatus,
			}
			jsonhttptest.ResponseDirect(t, srv.Client, http.MethodPost, testURL, body, tC.wantStatus, wantResponse)
			if !tC.wantFail {
				got := srv.P2PMock.GetWelcomeMessage()
				if got != tC.message {
					t.Fatalf("could not set dynamic welcome message: want %s, got %s", tC.message, got)
				}
			}
		})
	}
}

func TestSetWelcomeMessageInternalServerError(t *testing.T) {
	testMessage := "NO CHANCE BYE"
	testError := errors.New("Could not set value")
	testURL := "/welcome-message"

	srv := newTestServer(t, testServerOptions{
		P2P: mock.New(mock.WithSetWelcomeMessageFunc(func(string) error {
			return testError
		})),
	})

	data, _ := json.Marshal(debugapi.WelcomeMessageRequest{
		WelcomeMesssage: testMessage,
	})
	body := bytes.NewReader(data)
	t.Run("internal server error - failed to store", func(t *testing.T) {
		wantCode := http.StatusInternalServerError
		wantResp := jsonhttp.StatusResponse{
			Message: testError.Error(),
			Code:    wantCode,
		}
		jsonhttptest.ResponseDirect(t, srv.Client, http.MethodPost, testURL, body, wantCode, wantResp)
	})

}