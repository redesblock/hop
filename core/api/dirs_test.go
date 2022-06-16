package api_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/redesblock/hop/core/api"
	"github.com/redesblock/hop/core/collection/entry"
	"github.com/redesblock/hop/core/file"
	"github.com/redesblock/hop/core/file/joiner"
	"github.com/redesblock/hop/core/jsonhttp"
	"github.com/redesblock/hop/core/jsonhttp/jsonhttptest"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/manifest"
	statestore "github.com/redesblock/hop/core/statestore/mock"
	"github.com/redesblock/hop/core/storage/mock"
	"github.com/redesblock/hop/core/swarm"
	"github.com/redesblock/hop/core/tags"
)

func TestDirs(t *testing.T) {
	var (
		dirUploadResource    = "/dirs"
		fileDownloadResource = func(addr string) string { return "/files/" + addr }
		hopDownloadResource  = func(addr, path string) string { return "/hop/" + addr + "/" + path }
		storer               = mock.NewStorer()
		mockStatestore       = statestore.NewStateStore()
		logger               = logging.New(ioutil.Discard, 0)
		client, _, _         = newTestServer(t, testServerOptions{
			Storer:          storer,
			Tags:            tags.NewTags(mockStatestore, logger),
			Logger:          logging.New(ioutil.Discard, 5),
			PreventRedirect: true,
		})
	)

	t.Run("empty request body", func(t *testing.T) {
		jsonhttptest.Request(t, client, http.MethodPost, dirUploadResource, http.StatusBadRequest,
			jsonhttptest.WithRequestBody(bytes.NewReader(nil)),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "could not validate request",
				Code:    http.StatusBadRequest,
			}),
			jsonhttptest.WithRequestHeader("Content-Type", api.ContentTypeTar),
		)
	})

	t.Run("non tar file", func(t *testing.T) {
		file := bytes.NewReader([]byte("some data"))

		jsonhttptest.Request(t, client, http.MethodPost, dirUploadResource, http.StatusInternalServerError,
			jsonhttptest.WithRequestBody(file),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "could not store dir",
				Code:    http.StatusInternalServerError,
			}),
			jsonhttptest.WithRequestHeader("Content-Type", api.ContentTypeTar),
		)
	})

	t.Run("wrong content type", func(t *testing.T) {
		tarReader := tarFiles(t, []f{{
			data: []byte("some data"),
			name: "binary-file",
		}})

		// submit valid tar, but with wrong content-type
		jsonhttptest.Request(t, client, http.MethodPost, dirUploadResource, http.StatusBadRequest,
			jsonhttptest.WithRequestBody(tarReader),
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: "could not validate request",
				Code:    http.StatusBadRequest,
			}),
			jsonhttptest.WithRequestHeader("Content-Type", "other"),
		)
	})

	// valid tars
	for _, tc := range []struct {
		name                string
		wantIndexFilename   string
		wantErrorFilename   string
		indexFilenameOption jsonhttptest.Option
		errorFilenameOption jsonhttptest.Option
		files               []f // files in dir for test case
	}{
		{
			name: "non-nested files without extension",
			files: []f{
				{
					data:      []byte("first file data"),
					name:      "file1",
					dir:       "",
					reference: swarm.MustParseHexAddress("3c07cd2cf5c46208d69d554b038f4dce203f53ac02cb8a313a0fe1e3fe6cc3cf"),
					header: http.Header{
						"Content-Type": {""},
					},
				},
				{
					data:      []byte("second file data"),
					name:      "file2",
					dir:       "",
					reference: swarm.MustParseHexAddress("47e1a2a8f16e02da187fac791d57e6794f3e9b5d2400edd00235da749ad36683"),
					header: http.Header{
						"Content-Type": {""},
					},
				},
			},
		},
		{
			name: "nested files with extension",
			files: []f{
				{
					data:      []byte("robots text"),
					name:      "robots.txt",
					dir:       "",
					reference: swarm.MustParseHexAddress("17b96d0a800edca59aaf7e40c6053f7c4c0fb80dd2eb3f8663d51876bf350b12"),
					header: http.Header{
						"Content-Type": {"text/plain; charset=utf-8"},
					},
				},
				{
					data:      []byte("image 1"),
					name:      "1.png",
					dir:       "img",
					reference: swarm.MustParseHexAddress("3c1b3fc640e67f0595d9c1db23f10c7a2b0bdc9843b0e27c53e2ac2a2d6c4674"),
					header: http.Header{
						"Content-Type": {"image/png"},
					},
				},
				{
					data:      []byte("image 2"),
					name:      "2.png",
					dir:       "img",
					reference: swarm.MustParseHexAddress("b234ea7954cab7b2ccc5e07fe8487e932df11b2275db6b55afcbb7bad0be73fb"),
					header: http.Header{
						"Content-Type": {"image/png"},
					},
				},
			},
		},
		{
			name: "no index filename",
			files: []f{
				{
					data:      []byte("<h1>Swarm"),
					name:      "index.html",
					dir:       "",
					reference: swarm.MustParseHexAddress("bcb1bfe15c36f1a529a241f4d0c593e5648aa6d40859790894c6facb41a6ef28"),
					header: http.Header{
						"Content-Type": {"text/html; charset=utf-8"},
					},
				},
			},
		},
		{
			name:                "explicit index filename",
			wantIndexFilename:   "index.html",
			indexFilenameOption: jsonhttptest.WithRequestHeader(api.SwarmIndexDocumentHeader, "index.html"),
			files: []f{
				{
					data:      []byte("<h1>Swarm"),
					name:      "index.html",
					dir:       "",
					reference: swarm.MustParseHexAddress("bcb1bfe15c36f1a529a241f4d0c593e5648aa6d40859790894c6facb41a6ef28"),
					header: http.Header{
						"Content-Type": {"text/html; charset=utf-8"},
					},
				},
			},
		},
		{
			name:                "nested index filename",
			wantIndexFilename:   "index.html",
			indexFilenameOption: jsonhttptest.WithRequestHeader(api.SwarmIndexDocumentHeader, "index.html"),
			files: []f{
				{
					data:      []byte("<h1>Swarm"),
					name:      "index.html",
					dir:       "dir",
					reference: swarm.MustParseHexAddress("bcb1bfe15c36f1a529a241f4d0c593e5648aa6d40859790894c6facb41a6ef28"),
					header: http.Header{
						"Content-Type": {"text/html; charset=utf-8"},
					},
				},
			},
		},
		{
			name:                "explicit index and error filename",
			wantIndexFilename:   "index.html",
			wantErrorFilename:   "error.html",
			indexFilenameOption: jsonhttptest.WithRequestHeader(api.SwarmIndexDocumentHeader, "index.html"),
			errorFilenameOption: jsonhttptest.WithRequestHeader(api.SwarmErrorDocumentHeader, "error.html"),
			files: []f{
				{
					data:      []byte("<h1>Swarm"),
					name:      "index.html",
					dir:       "",
					reference: swarm.MustParseHexAddress("bcb1bfe15c36f1a529a241f4d0c593e5648aa6d40859790894c6facb41a6ef28"),
					header: http.Header{
						"Content-Type": {"text/html; charset=utf-8"},
					},
				},
				{
					data:      []byte("<h2>404"),
					name:      "error.html",
					dir:       "",
					reference: swarm.MustParseHexAddress("b1f309c095d650521b75760b23122a9c59c2b581af28fc6daaf9c58da86a204d"),
					header: http.Header{
						"Content-Type": {"text/html; charset=utf-8"},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// tar all the test case files
			tarReader := tarFiles(t, tc.files)

			var respBytes []byte

			options := []jsonhttptest.Option{
				jsonhttptest.WithRequestBody(tarReader),
				jsonhttptest.WithRequestHeader("Content-Type", api.ContentTypeTar),
				jsonhttptest.WithPutResponseBody(&respBytes),
			}
			if tc.indexFilenameOption != nil {
				options = append(options, tc.indexFilenameOption)
			}
			if tc.errorFilenameOption != nil {
				options = append(options, tc.errorFilenameOption)
			}

			// verify directory tar upload response
			jsonhttptest.Request(t, client, http.MethodPost, dirUploadResource, http.StatusOK, options...)

			read := bytes.NewReader(respBytes)

			// get the reference as everytime it will change because of random encryption key
			var resp api.FileUploadResponse
			err := json.NewDecoder(read).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			// NOTE: reference will be different each time, due to manifest randomness

			if resp.Reference.String() == "" {
				t.Fatalf("expected file reference, did not got any")
			}

			// read manifest metadata
			j, _, err := joiner.New(context.Background(), storer, resp.Reference)
			if err != nil {
				t.Fatal(err)
			}

			buf := bytes.NewBuffer(nil)
			_, err = file.JoinReadAll(context.Background(), j, buf)
			if err != nil {
				t.Fatal(err)
			}
			e := &entry.Entry{}
			err = e.UnmarshalBinary(buf.Bytes())
			if err != nil {
				t.Fatal(err)
			}

			// verify manifest content
			verifyManifest, err := manifest.NewManifestReference(
				context.Background(),
				manifest.DefaultManifestType,
				e.Reference(),
				false,
				storer,
			)
			if err != nil {
				t.Fatal(err)
			}

			validateFile := func(t *testing.T, file f, filePath string) {
				t.Helper()

				entry, err := verifyManifest.Lookup(filePath)
				if err != nil {
					t.Fatal(err)
				}

				fileReference := entry.Reference()

				if !bytes.Equal(file.reference.Bytes(), fileReference.Bytes()) {
					t.Fatalf("expected file reference to match %s, got %s", file.reference, fileReference)
				}

				jsonhttptest.Request(t, client, http.MethodGet, fileDownloadResource(fileReference.String()), http.StatusOK,
					jsonhttptest.WithExpectedResponse(file.data),
					jsonhttptest.WithRequestHeader("Content-Type", file.header.Get("Content-Type")),
				)
			}

			validateIsPermanentRedirect := func(t *testing.T, fromPath, toPath string) {
				t.Helper()

				expectedResponse := fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.\n\n", hopDownloadResource(resp.Reference.String(), toPath))

				jsonhttptest.Request(t, client, http.MethodGet, hopDownloadResource(resp.Reference.String(), fromPath), http.StatusPermanentRedirect,
					jsonhttptest.WithExpectedResponse([]byte(expectedResponse)),
				)
			}

			validateHopPath := func(t *testing.T, fromPath, toPath string) {
				t.Helper()

				toEntry, err := verifyManifest.Lookup(toPath)
				if err != nil {
					t.Fatal(err)
				}

				var respBytes []byte

				jsonhttptest.Request(t, client, http.MethodGet, fileDownloadResource(toEntry.Reference().String()), http.StatusOK,
					jsonhttptest.WithPutResponseBody(&respBytes),
				)

				jsonhttptest.Request(t, client, http.MethodGet, hopDownloadResource(resp.Reference.String(), fromPath), http.StatusOK,
					jsonhttptest.WithExpectedResponse(respBytes),
				)
			}

			// check if each file can be located and read
			for _, file := range tc.files {
				validateFile(t, file, path.Join(file.dir, file.name))
			}

			// check index filename
			if tc.wantIndexFilename != "" {
				entry, err := verifyManifest.Lookup(api.ManifestRootPath)
				if err != nil {
					t.Fatal(err)
				}

				manifestRootMetadata := entry.Metadata()
				indexDocumentSuffixPath, ok := manifestRootMetadata[api.ManifestWebsiteIndexDocumentSuffixKey]
				if !ok {
					t.Fatalf("expected index filename '%s', did not find any", tc.wantIndexFilename)
				}

				// check index suffix for each dir
				for _, file := range tc.files {
					if file.dir != "" {
						validateIsPermanentRedirect(t, file.dir, file.dir+"/")
						validateHopPath(t, file.dir+"/", path.Join(file.dir, indexDocumentSuffixPath))
					}
				}
			}

			// check error filename
			if tc.wantErrorFilename != "" {
				entry, err := verifyManifest.Lookup(api.ManifestRootPath)
				if err != nil {
					t.Fatal(err)
				}

				manifestRootMetadata := entry.Metadata()
				errorDocumentPath, ok := manifestRootMetadata[api.ManifestWebsiteErrorDocumentPathKey]
				if !ok {
					t.Fatalf("expected error filename '%s', did not find any", tc.wantErrorFilename)
				}

				// check error document
				validateHopPath(t, "_non_existent_file_path_", errorDocumentPath)
			}

		})
	}
}

// tarFiles receives an array of test case files and creates a new tar with those files as a collection
// it returns a bytes.Buffer which can be used to read the created tar
func tarFiles(t *testing.T, files []f) *bytes.Buffer {
	t.Helper()

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
	header    http.Header
}
