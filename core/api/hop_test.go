package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"testing"

	statestore "github.com/redesblock/hop/core/statestore/mock"

	"github.com/redesblock/hop/core/collection/entry"
	"github.com/redesblock/hop/core/file/pipeline"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/manifest"
	"github.com/redesblock/hop/core/storage"
	smock "github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
)

func TestHop(t *testing.T) {
	var (
		hopDownloadResource = func(addr, path string) string { return "/hop/" + addr + "/" + path }
		storer              = smock.NewStorer()
		ctx                 = context.Background()
		mockStatestore      = statestore.NewStateStore()
		logger              = logging.New(ioutil.Discard, 0)
		client, _, _        = newTestServer(t, testServerOptions{
			Storer: storer,
			Tags:   tags.NewTags(mockStatestore, logger),
			Logger: logging.New(ioutil.Discard, 5),
		})
		pipeWriteAll = func(r io.Reader, l int64) (swarm.Address, error) {
			pipe := pipeline.NewPipelineBuilder(ctx, storer, storage.ModePutUpload, false)
			return pipeline.FeedPipeline(ctx, pipe, r, l)
		}
	)
	t.Run("download-file-by-path", func(t *testing.T) {
		fileName := "sample.html"
		filePath := "test/" + fileName
		missingFilePath := "test/missing"
		sampleHtml := `<!DOCTYPE html>
		<html>
		<body>
	
		<h1>My First Heading</h1>
	
		<p>My first paragraph.</p>
	
		</body>
		</html>`

		var err error
		var fileContentReference swarm.Address
		var fileReference swarm.Address
		var manifestFileReference swarm.Address

		// save file
		fileContentReference, err = pipeWriteAll(strings.NewReader(sampleHtml), int64(len(sampleHtml)))

		if err != nil {
			t.Fatal(err)
		}

		fileMetadata := entry.NewMetadata(fileName)
		fileMetadata.MimeType = "text/html; charset=utf-8"
		fileMetadataBytes, err := json.Marshal(fileMetadata)
		if err != nil {
			t.Fatal(err)
		}

		fileMetadataReference, err := pipeWriteAll(bytes.NewReader(fileMetadataBytes), int64(len(fileMetadataBytes)))

		if err != nil {
			t.Fatal(err)
		}

		fe := entry.New(fileContentReference, fileMetadataReference)
		fileEntryBytes, err := fe.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		fileReference, err = pipeWriteAll(bytes.NewReader(fileEntryBytes), int64(len(fileEntryBytes)))

		if err != nil {
			t.Fatal(err)
		}

		// save manifest
		m, err := manifest.NewDefaultManifest(false, storer)
		if err != nil {
			t.Fatal(err)
		}

		e := manifest.NewEntry(fileReference)

		err = m.Add(filePath, e)
		if err != nil {
			t.Fatal(err)
		}

		manifestBytesReference, err := m.Store(context.Background(), storage.ModePutUpload)
		if err != nil {
			t.Fatal(err)
		}

		metadata := entry.NewMetadata(manifestBytesReference.String())
		metadata.MimeType = m.Type()
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			t.Fatal(err)
		}

		mr, err := pipeWriteAll(bytes.NewReader(metadataBytes), int64(len(metadataBytes)))
		if err != nil {
			t.Fatal(err)
		}

		// now join both references (fr,mr) to create an entry and store it.
		newEntry := entry.New(manifestBytesReference, mr)
		manifestFileEntryBytes, err := newEntry.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		manifestFileReference, err = pipeWriteAll(bytes.NewReader(manifestFileEntryBytes), int64(len(manifestFileEntryBytes)))
		if err != nil {
			t.Fatal(err)
		}

		// read file from manifest path

		rcvdHeader := jsonhttptest.Request(t, client, http.MethodGet, hopDownloadResource(manifestFileReference.String(), filePath), http.StatusOK,
			jsonhttptest.WithExpectedResponse([]byte(sampleHtml)),
		)
		cd := rcvdHeader.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(cd)
		if err != nil {
			t.Fatal(err)
		}
		if params["filename"] != fileName {
			t.Fatal("Invalid file name detected")
		}
		if rcvdHeader.Get("ETag") != fmt.Sprintf("%q", fileContentReference) {
			t.Fatal("Invalid ETags header received")
		}
		if rcvdHeader.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Fatal("Invalid content type detected")
		}

		// check on invalid path

		jsonhttptest.Request(t, client, http.MethodGet, hopDownloadResource(manifestFileReference.String(), missingFilePath), http.StatusNotFound,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "path address not found",
				Code:    http.StatusNotFound,
			}),
		)
	})

}
