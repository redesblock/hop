package api_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/logging"
	statestore "github.com/redesblock/hop/core/statestore/mock"
	"github.com/redesblock/hop/core/storage/mock"
	testingc "github.com/redesblock/hop/core/storage/testing"
	"github.com/redesblock/hop/core/tags"
)

func TestGatewayMode(t *testing.T) {
	logger := logging.New(ioutil.Discard, 0)
	chunk := testingc.GenerateTestRandomChunk()
	client, _, _ := newTestServer(t, testServerOptions{
		Storer:      mock.NewStorer(),
		Tags:        tags.NewTags(statestore.NewStateStore(), logger),
		Logger:      logger,
		GatewayMode: true,
	})

	forbiddenResponseOption := jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
		Message: http.StatusText(http.StatusForbidden),
		Code:    http.StatusForbidden,
	})

	t.Run("pinning endpoints", func(t *testing.T) {
		path := "/pin/chunks/0773a91efd6547c754fc1d95fb1c62c7d1b47f959c2caa685dfec8736da95c1c"
		jsonhttptest.Request(t, client, http.MethodGet, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodPost, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodDelete, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodGet, "/pin/chunks", http.StatusForbidden, forbiddenResponseOption)
	})

	t.Run("tags endpoints", func(t *testing.T) {
		path := "/tags/42"
		jsonhttptest.Request(t, client, http.MethodGet, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodDelete, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodPatch, path, http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodGet, "/tags", http.StatusForbidden, forbiddenResponseOption)
	})

	t.Run("pss endpoints", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodPost, "/pss/send/test-topic/ab", http.StatusForbidden, forbiddenResponseOption)
		jsonhttptest.Request(t, client, http.MethodGet, "/pss/subscribe/test-topic", http.StatusForbidden, forbiddenResponseOption)
	})

	t.Run("pinning", func(t *testing.T) {
		headerOption := jsonhttptest.WithRequestHeader(api.SwarmPinHeader, "true")

		forbiddenResponseOption := jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
			Message: "pinning is disabled",
			Code:    http.StatusForbidden,
		})

		// should work without pinning
		jsonhttptest.Request(t, client, http.MethodPost, "/chunks/"+chunk.Address().String(), http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader(chunk.Data())),
		)

		jsonhttptest.Request(t, client, http.MethodPost, "/chunks/0773a91efd6547c754fc1d95fb1c62c7d1b47f959c2caa685dfec8736da95c1c", http.StatusForbidden, forbiddenResponseOption, headerOption)

		jsonhttptest.Request(t, client, http.MethodPost, "/bytes", http.StatusOK) // should work without pinning
		jsonhttptest.Request(t, client, http.MethodPost, "/bytes", http.StatusForbidden, forbiddenResponseOption, headerOption)
		jsonhttptest.Request(t, client, http.MethodPost, "/files", http.StatusForbidden, forbiddenResponseOption, headerOption)
		jsonhttptest.Request(t, client, http.MethodPost, "/dirs", http.StatusForbidden, forbiddenResponseOption, headerOption)
	})

	t.Run("encryption", func(t *testing.T) {
		headerOption := jsonhttptest.WithRequestHeader(api.SwarmEncryptHeader, "true")

		forbiddenResponseOption := jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
			Message: "encryption is disabled",
			Code:    http.StatusForbidden,
		})

		jsonhttptest.Request(t, client, http.MethodPost, "/bytes", http.StatusOK) // should work without pinning
		jsonhttptest.Request(t, client, http.MethodPost, "/bytes", http.StatusForbidden, forbiddenResponseOption, headerOption)
		jsonhttptest.Request(t, client, http.MethodPost, "/files", http.StatusForbidden, forbiddenResponseOption, headerOption)
		jsonhttptest.Request(t, client, http.MethodPost, "/dirs", http.StatusForbidden, forbiddenResponseOption, headerOption)
	})
}
