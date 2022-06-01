package api_test

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/manifest/jsonmanifest"
	"github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
)

func TestDirs(t *testing.T) {
	var (
		dirUploadResource    = "/dirs"
		fileDownloadResource = func(addr string) string { return "/files/" + addr }
		client               = newTestServer(t, testServerOptions{
			Storer: mock.NewStorer(),
			Tags:   tags.NewTags(),
			Logger: logging.New(ioutil.Discard, 5),
		})
	)

	t.Run("empty request body", func(t *testing.T) {
		jsonhttptest.ResponseDirectSendHeadersAndReceiveHeaders(t, client, http.MethodPost, dirUploadResource, bytes.NewReader(nil), http.StatusBadRequest, jsonhttp.StatusResponse{
			Message: "could not validate request",
			Code:    http.StatusBadRequest,
		}, http.Header{
			"Content-Type": {api.ContentTypeTar},
		})
	})

	t.Run("non tar file", func(t *testing.T) {
		file := bytes.NewReader([]byte("some data"))

		jsonhttptest.ResponseDirectSendHeadersAndReceiveHeaders(t, client, http.MethodPost, dirUploadResource, file, http.StatusInternalServerError, jsonhttp.StatusResponse{
			Message: "could not store dir",
			Code:    http.StatusInternalServerError,
		}, http.Header{
			"Content-Type": {api.ContentTypeTar},
		})
	})

	t.Run("wrong content type", func(t *testing.T) {
		tarReader := tarFiles(t, []f{{
			data: []byte("some data"),
			name: "binary-file",
		}})

		// submit valid tar, but with wrong content-type
		jsonhttptest.ResponseDirectSendHeadersAndReceiveHeaders(t, client, http.MethodPost, dirUploadResource, tarReader, http.StatusBadRequest, jsonhttp.StatusResponse{
			Message: "could not validate request",
			Code:    http.StatusBadRequest,
		}, http.Header{
			"Content-Type": {"other"},
		})
	})

	// valid tars
	for _, tc := range []struct {
		name         string
		expectedHash string
		files        []f // files in dir for test case
	}{
		{
			name:         "non-nested files without extension",
			expectedHash: "2fa041bd35ebff676727eb3023272f43b1e0fa71c8735cc1a7487e9131f963c4",
			files: []f{
				{
					data:      []byte("first file data"),
					name:      "file1",
					dir:       "",
					reference: swarm.MustParseHexAddress("3c07cd2cf5c46208d69d554b038f4dce203f53ac02cb8a313a0fe1e3fe6cc3cf"),
					headers: http.Header{
						"Content-Type": {""},
					},
				},
				{
					data:      []byte("second file data"),
					name:      "file2",
					dir:       "",
					reference: swarm.MustParseHexAddress("47e1a2a8f16e02da187fac791d57e6794f3e9b5d2400edd00235da749ad36683"),
					headers: http.Header{
						"Content-Type": {""},
					},
				},
			},
		},
		{
			name:         "nested files with extension",
			expectedHash: "c3cb9fbe2efa7bbc979245d9bac1400bd4894371776b7560309d49e687514dd6",
			files: []f{
				{
					data:      []byte("robots text"),
					name:      "robots.txt",
					dir:       "",
					reference: swarm.MustParseHexAddress("17b96d0a800edca59aaf7e40c6053f7c4c0fb80dd2eb3f8663d51876bf350b12"),
					headers: http.Header{
						"Content-Type": {"text/plain; charset=utf-8"},
					},
				},
				{
					data:      []byte("image 1"),
					name:      "1.png",
					dir:       "img",
					reference: swarm.MustParseHexAddress("3c1b3fc640e67f0595d9c1db23f10c7a2b0bdc9843b0e27c53e2ac2a2d6c4674"),
					headers: http.Header{
						"Content-Type": {"image/png"},
					},
				},
				{
					data:      []byte("image 2"),
					name:      "2.png",
					dir:       "img",
					reference: swarm.MustParseHexAddress("b234ea7954cab7b2ccc5e07fe8487e932df11b2275db6b55afcbb7bad0be73fb"),
					headers: http.Header{
						"Content-Type": {"image/png"},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// tar all the test case files
			tarReader := tarFiles(t, tc.files)

			// verify directory tar upload response
			jsonhttptest.ResponseDirectSendHeadersAndReceiveHeaders(t, client, http.MethodPost, dirUploadResource, tarReader, http.StatusOK, api.FileUploadResponse{
				Reference: swarm.MustParseHexAddress(tc.expectedHash),
			}, http.Header{
				"Content-Type": {api.ContentTypeTar},
			})

			// create expected manifest
			expectedManifest := jsonmanifest.NewManifest()
			for _, file := range tc.files {
				e := jsonmanifest.NewEntry(file.reference, file.name, file.headers)
				expectedManifest.Add(path.Join(file.dir, file.name), e)
			}

			b, err := expectedManifest.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			// verify directory upload manifest through files api
			jsonhttptest.ResponseDirectCheckBinaryResponse(t, client, http.MethodGet, fileDownloadResource(tc.expectedHash), nil, http.StatusOK, b, nil)
		})
	}
}

// tarFiles receives an array of test case files and creates a new tar with those files as a collection
// it returns a bytes.Buffer which can be used to read the created tar
func tarFiles(t *testing.T, files []f) *bytes.Buffer {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, file := range files {
		// create tar header and write it
		hdr := &tar.Header{
			Name: path.Join(file.dir, file.name),
			Mode: 0600,
			Size: int64(len(file.data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		// write the file data to the tar
		if _, err := tw.Write(file.data); err != nil {
			t.Fatal(err)
		}
	}

	// finally close the tar writer
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	return &buf
}

// struct for dir files for test cases
type f struct {
	data      []byte
	name      string
	dir       string
	reference swarm.Address
	headers   http.Header
}
