package api_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/redesblock/hop/core/logging"
	statestore "github.com/redesblock/hop/core/statestore/mock"

	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/storage/mock"
	testingc "github.com/redesblock/hop/core/storage/testing"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/swarm/test"
	"github.com/redesblock/hop/core/tags"
	"gitlab.com/nolash/go-mockbytes"
)

type fileUploadResponse struct {
	Reference swarm.Address `json:"reference"`
}

func tagsWithIdResource(id uint32) string { return fmt.Sprintf("/tags/%d", id) }

func TestTags(t *testing.T) {
	var (
		filesResource  = "/files"
		dirResource    = "/dirs"
		bytesResource  = "/bytes"
		chunksResource = "/chunks"
		tagsResource   = "/tags"
		chunk          = testingc.GenerateTestRandomChunk()
		someTagName    = "file.jpg"
		mockStatestore = statestore.NewStateStore()
		logger         = logging.New(ioutil.Discard, 0)
		tag            = tags.NewTags(mockStatestore, logger)
		client, _, _   = newTestServer(t, testServerOptions{
			Storer: mock.NewStorer(),
			Tags:   tag,
		})
	)

	t.Run("create unnamed tag", func(t *testing.T) {
		tr := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagRequest{}),
			jsonhttptest.WithUnmarshalJSONResponse(&tr),
		)

		if !strings.Contains(tr.Name, "unnamed_tag_") {
			t.Fatalf("expected tag name to contain %s but is %s instead", "unnamed_tag_", tr.Name)
		}
	})

	t.Run("create tag with name", func(t *testing.T) {
		tr := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagRequest{
				Name: someTagName,
			}),
			jsonhttptest.WithUnmarshalJSONResponse(&tr),
		)

		if tr.Name != someTagName {
			t.Fatalf("expected tag name to be %s but is %s instead", someTagName, tr.Name)
		}
	})

	t.Run("create tag with invalid id", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodPost, chunksResource, http.StatusBadRequest,
			jsonhttptest.WithRequestBody(bytes.NewReader(chunk.Data())),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "cannot get tag",
				Code:    http.StatusBadRequest,
			}),
			jsonhttptest.WithRequestHeader(api.SwarmTagUidHeader, "invalid_id.jpg"), // the value should be uint32
		)
	})

	t.Run("get invalid tags", func(t *testing.T) {
		// invalid tag
		jsonhttptest.Request(t, client, http.MethodGet, tagsResource+"/foobar", http.StatusBadRequest,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "invalid id",
				Code:    http.StatusBadRequest,
			}),
		)

		// non-existent tag
		jsonhttptest.Request(t, client, http.MethodDelete, tagsWithIdResource(uint32(333)), http.StatusNotFound,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "tag not present",
				Code:    http.StatusNotFound,
			}),
		)
	})

	t.Run("create tag upload chunk", func(t *testing.T) {
		// create a tag using the API
		tr := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{
				Name: someTagName,
			}),
			jsonhttptest.WithUnmarshalJSONResponse(&tr),
		)

		if tr.Name != someTagName {
			t.Fatalf("sent tag name %s does not match received tag name %s", someTagName, tr.Name)
		}

		_ = jsonhttptest.Request(t, client, http.MethodPost, chunksResource, http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader(chunk.Data())),
			jsonhttptest.WithExpectedJSONResponse(api.ChunkAddressResponse{Reference: chunk.Address()}),
		)

		rcvdHeaders := jsonhttptest.Request(t, client, http.MethodPost, chunksResource, http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader(chunk.Data())),
			jsonhttptest.WithExpectedJSONResponse(api.ChunkAddressResponse{Reference: chunk.Address()}),
			jsonhttptest.WithRequestHeader(api.SwarmTagUidHeader, strconv.FormatUint(uint64(tr.Uid), 10)),
		)

		isTagFoundInResponse(t, rcvdHeaders, &tr)
		tagValueTest(t, tr.Uid, 1, 1, 1, 0, 0, 0, swarm.ZeroAddress, client)
	})

	t.Run("delete tag error", func(t *testing.T) {
		// try to delete invalid tag
		jsonhttptest.Request(t, client, http.MethodDelete, tagsResource+"/foobar", http.StatusBadRequest,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "invalid id",
				Code:    http.StatusBadRequest,
			}),
		)

		// try to delete non-existent tag
		jsonhttptest.Request(t, client, http.MethodDelete, tagsWithIdResource(uint32(333)), http.StatusNotFound,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "tag not present",
				Code:    http.StatusNotFound,
			}),
		)
	})

	t.Run("delete tag", func(t *testing.T) {
		// create a tag through API
		tRes := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{
				Name: someTagName,
			}),
			jsonhttptest.WithUnmarshalJSONResponse(&tRes),
		)

		// delete tag through API
		jsonhttptest.Request(t, client, http.MethodDelete, tagsWithIdResource(tRes.Uid), http.StatusNoContent,
			jsonhttptest.WithNoResponseBody(),
		)

		// try to get tag
		jsonhttptest.Request(t, client, http.MethodGet, tagsWithIdResource(tRes.Uid), http.StatusNotFound,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "tag not present",
				Code:    http.StatusNotFound,
			}),
		)
	})

	t.Run("done split error", func(t *testing.T) {
		// invalid tag
		jsonhttptest.Request(t, client, http.MethodPatch, tagsResource+"/foobar", http.StatusBadRequest,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{}),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "invalid id",
				Code:    http.StatusBadRequest,
			}),
		)

		// non-existent tag
		jsonhttptest.Request(t, client, http.MethodPatch, tagsWithIdResource(uint32(333)), http.StatusNotFound,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{}),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "tag not present",
				Code:    http.StatusNotFound,
			}),
		)
	})

	t.Run("done split", func(t *testing.T) {
		// create a tag through API
		tRes := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{
				Name: someTagName,
			}),
			jsonhttptest.WithUnmarshalJSONResponse(&tRes),
		)
		tagId := tRes.Uid

		// generate address to be supplied to the done split
		addr := test.RandomAddress()

		// upload content with tag
		jsonhttptest.Request(t, client, http.MethodPost, chunksResource, http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader(chunk.Data())),
			jsonhttptest.WithRequestHeader(api.SwarmTagUidHeader, fmt.Sprint(tagId)),
		)

		// call done split
		jsonhttptest.Request(t, client, http.MethodPatch, tagsWithIdResource(tagId), http.StatusOK,
			jsonhttptest.WithJSONRequestBody(api.TagRequest{
				Address: addr,
			}),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "ok",
				Code:    http.StatusOK,
			}),
		)
		tagValueTest(t, tagId, 1, 1, 1, 0, 0, 1, addr, client)

		// try different address value
		addr = test.RandomAddress()

		// call done split
		jsonhttptest.Request(t, client, http.MethodPatch, tagsWithIdResource(tagId), http.StatusOK,
			jsonhttptest.WithJSONRequestBody(api.TagRequest{
				Address: addr,
			}),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "ok",
				Code:    http.StatusOK,
			}),
		)
		tagValueTest(t, tagId, 1, 1, 1, 0, 0, 1, addr, client)
	})

	t.Run("file tags", func(t *testing.T) {
		// upload a file without supplying tag
		expectedHash := swarm.MustParseHexAddress("8e27bb803ff049e8c2f4650357026723220170c15ebf9b635a7026539879a1a8")
		expectedResponse := api.FileUploadResponse{Reference: expectedHash}

		respHeaders := jsonhttptest.Request(t, client, http.MethodPost, filesResource, http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader([]byte("some data"))),
			jsonhttptest.WithExpectedJSONResponse(expectedResponse),
			jsonhttptest.WithRequestHeader("Content-Type", "application/octet-stream"),
		)

		tagId, err := strconv.Atoi(respHeaders.Get(api.SwarmTagUidHeader))
		if err != nil {
			t.Fatal(err)
		}
		tagValueTest(t, uint32(tagId), 3, 3, 0, 0, 0, 3, expectedHash, client)
	})

	t.Run("dir tags", func(t *testing.T) {
		// upload a dir without supplying tag
		tarReader := tarFiles(t, []f{{
			data: []byte("some dir data"),
			name: "binary-file",
		}})
		expectedHash := swarm.MustParseHexAddress("3dc643abeb3db60a4dfb72008b577dd9a573abaa74c6afe37a75c63ceea829f6")
		expectedResponse := api.FileUploadResponse{Reference: expectedHash}

		respHeaders := jsonhttptest.Request(t, client, http.MethodPost, dirResource, http.StatusOK,
			jsonhttptest.WithRequestBody(tarReader),
			jsonhttptest.WithExpectedJSONResponse(expectedResponse),
			jsonhttptest.WithRequestHeader("Content-Type", api.ContentTypeTar),
		)

		tagId, err := strconv.Atoi(respHeaders.Get(api.SwarmTagUidHeader))
		if err != nil {
			t.Fatal(err)
		}
		tagValueTest(t, uint32(tagId), 7, 7, 0, 0, 0, 7, expectedHash, client)
	})

	t.Run("bytes tags", func(t *testing.T) {
		// create a tag using the API
		tr := api.TagResponse{}
		jsonhttptest.Request(t, client, http.MethodPost, tagsResource, http.StatusCreated,
			jsonhttptest.WithJSONRequestBody(api.TagResponse{
				Name: someTagName,
			}),
			jsonhttptest.WithUnmarshalJSONResponse(&tr),
		)
		if tr.Name != someTagName {
			t.Fatalf("sent tag name %s does not match received tag name %s", someTagName, tr.Name)
		}

		sentHeaders := make(http.Header)
		sentHeaders.Set(api.SwarmTagUidHeader, strconv.FormatUint(uint64(tr.Uid), 10))

		g := mockbytes.New(0, mockbytes.MockTypeStandard).WithModulus(255)
		dataChunk, err := g.SequentialBytes(swarm.ChunkSize)
		if err != nil {
			t.Fatal(err)
		}

		rootAddress := swarm.MustParseHexAddress("5e2a21902f51438be1adbd0e29e1bd34c53a21d3120aefa3c7275129f2f88de9")

		content := make([]byte, swarm.ChunkSize*2)
		copy(content[swarm.ChunkSize:], dataChunk)
		copy(content[:swarm.ChunkSize], dataChunk)

		rcvdHeaders := jsonhttptest.Request(t, client, http.MethodPost, bytesResource, http.StatusOK,
			jsonhttptest.WithRequestBody(bytes.NewReader(content)),
			jsonhttptest.WithExpectedJSONResponse(fileUploadResponse{
				Reference: rootAddress,
			}),
			jsonhttptest.WithRequestHeader(api.SwarmTagUidHeader, strconv.FormatUint(uint64(tr.Uid), 10)),
		)
		id := isTagFoundInResponse(t, rcvdHeaders, nil)

		tagToVerify, err := tag.Get(id)
		if err != nil {
			t.Fatal(err)
		}

		if tagToVerify.Uid != tr.Uid {
			t.Fatalf("expected tag id to be %d but is %d", tagToVerify.Uid, tr.Uid)
		}
		tagValueTest(t, id, 3, 3, 1, 0, 0, 0, swarm.ZeroAddress, client)
	})
}

// isTagFoundInResponse verifies that the tag id is found in the supplied HTTP headers
// if an API tag response is supplied, it also verifies that it contains an id which matches the headers
func isTagFoundInResponse(t *testing.T, headers http.Header, tr *api.TagResponse) uint32 {
	t.Helper()

	idStr := headers.Get(api.SwarmTagUidHeader)
	if idStr == "" {
		t.Fatalf("could not find tag id header in chunk upload response")
	}
	nId, err := strconv.Atoi(idStr)
	id := uint32(nId)
	if err != nil {
		t.Fatal(err)
	}
	if tr != nil {
		if id != tr.Uid {
			t.Fatalf("expected created tag id to be %d, but got %d when uploading chunk", tr.Uid, id)
		}
	}
	return id
}

func tagValueTest(t *testing.T, id uint32, split, stored, seen, sent, synced, total int64, address swarm.Address, client *http.Client) {
	t.Helper()
	tag := api.TagResponse{}
	jsonhttptest.Request(t, client, http.MethodGet, tagsWithIdResource(id), http.StatusOK,
		jsonhttptest.WithUnmarshalJSONResponse(&tag),
	)

	if tag.Split != split {
		t.Errorf("tag split count mismatch. got %d want %d", tag.Split, split)
	}
	if tag.Stored != stored {
		t.Errorf("tag stored count mismatch. got %d want %d", tag.Stored, stored)
	}
	if tag.Seen != seen {
		t.Errorf("tag seen count mismatch. got %d want %d", tag.Seen, seen)
	}
	if tag.Sent != sent {
		t.Errorf("tag sent count mismatch. got %d want %d", tag.Sent, sent)
	}
	if tag.Synced != synced {
		t.Errorf("tag synced count mismatch. got %d want %d", tag.Synced, synced)
	}
	if tag.Total != total {
		t.Errorf("tag total count mismatch. got %d want %d", tag.Total, total)
	}

	if !tag.Address.Equal(address) {
		t.Errorf("address mismatch: expected %s got %s", address.String(), tag.Address.String())
	}
}
