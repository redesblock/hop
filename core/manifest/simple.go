package manifest

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ethersphere/manifest/simple"
	"github.com/redesblock/hop/core/file"
	"github.com/redesblock/hop/core/file/pipeline/builder"
	"github.com/redesblock/hop/core/file/seekjoiner"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/swarm"
)

const (
	// ManifestSimpleContentType represents content type used for noting that
	// specific file should be processed as 'simple' manifest
	ManifestSimpleContentType = "application/hop-manifest-simple+json"
)

type simpleManifest struct {
	manifest simple.Manifest

	encrypted bool
	storer    storage.Storer
}

// NewSimpleManifest creates a new simple manifest.
func NewSimpleManifest(
	encrypted bool,
	storer storage.Storer,
) (Interface, error) {
	return &simpleManifest{
		manifest:  simple.NewManifest(),
		encrypted: encrypted,
		storer:    storer,
	}, nil
}

// NewSimpleManifestReference loads existing simple manifest.
func NewSimpleManifestReference(
	ctx context.Context,
	reference swarm.Address,
	encrypted bool,
	storer storage.Storer,
) (Interface, error) {
	m := &simpleManifest{
		manifest:  simple.NewManifest(),
		encrypted: encrypted,
		storer:    storer,
	}
	err := m.load(ctx, reference)
	return m, err
}

func (m *simpleManifest) Type() string {
	return ManifestSimpleContentType
}

func (m *simpleManifest) Add(path string, entry Entry) error {
	e := entry.Reference().String()

	return m.manifest.Add(path, e)
}

func (m *simpleManifest) Remove(path string) error {

	err := m.manifest.Remove(path)
	if err != nil {
		if errors.Is(err, simple.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (m *simpleManifest) Lookup(path string) (Entry, error) {

	n, err := m.manifest.Lookup(path)
	if err != nil {
		return nil, ErrNotFound
	}

	address, err := swarm.ParseHexAddress(n.Reference())
	if err != nil {
		return nil, fmt.Errorf("parse swarm address: %w", err)
	}

	entry := NewEntry(address)

	return entry, nil
}

func (m *simpleManifest) Store(ctx context.Context, mode storage.ModePut) (swarm.Address, error) {

	data, err := m.manifest.MarshalBinary()
	if err != nil {
		return swarm.ZeroAddress, fmt.Errorf("manifest marshal error: %w", err)
	}

	pipe := builder.NewPipelineBuilder(ctx, m.storer, mode, m.encrypted)
	address, err := builder.FeedPipeline(ctx, pipe, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return swarm.ZeroAddress, fmt.Errorf("manifest save error: %w", err)
	}

	return address, nil
}

func (m *simpleManifest) load(ctx context.Context, reference swarm.Address) error {
	j := seekjoiner.NewSimpleJoiner(m.storer)

	buf := bytes.NewBuffer(nil)
	_, err := file.JoinReadAll(ctx, j, reference, buf)
	if err != nil {
		return fmt.Errorf("manifest load error: %w", err)
	}

	err = m.manifest.UnmarshalBinary(buf.Bytes())
	if err != nil {
		return fmt.Errorf("manifest unmarshal error: %w", err)
	}

	return nil
}
