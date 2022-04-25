package libp2p_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/redesblock/hop/core/addressbook"
	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/logging"
	"github.com/redesblock/hop/core/p2p"
	"github.com/redesblock/hop/core/p2p/libp2p"
	"github.com/redesblock/hop/core/statestore/mock"
	"github.com/redesblock/hop/core/swarm"
)

// newService constructs a new libp2p service.
func newService(t *testing.T, networkID uint64, o libp2p.Options) (s *libp2p.Service, overlay swarm.Address, cleanup func()) {
	t.Helper()

	privateKey, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}

	overlay = crypto.NewOverlayAddress(privateKey.PublicKey, networkID)

	addr := ":0"

	if o.Logger == nil {
		o.Logger = logging.New(ioutil.Discard, 0)
	}

	if o.Addressbook == nil {
		statestore := mock.NewStateStore()
		o.Addressbook = addressbook.New(statestore)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s, err = libp2p.New(ctx, crypto.NewDefaultSigner(privateKey), networkID, overlay, addr, o)
	if err != nil {
		t.Fatal(err)
	}
	return s, overlay, func() {
		cancel()
		s.Close()
	}
}

// expectPeers validates that peers with addresses are connected.
func expectPeers(t *testing.T, s *libp2p.Service, addrs ...swarm.Address) {
	t.Helper()

	peers := s.Peers()

	if len(peers) != len(addrs) {
		t.Fatalf("got peers %v, want %v", len(peers), len(addrs))
	}

	sort.Slice(addrs, func(i, j int) bool {
		return bytes.Compare(addrs[i].Bytes(), addrs[j].Bytes()) == -1
	})
	sort.Slice(peers, func(i, j int) bool {
		return bytes.Compare(peers[i].Address.Bytes(), peers[j].Address.Bytes()) == -1
	})

	for i, got := range peers {
		want := addrs[i]
		if !got.Address.Equal(want) {
			t.Errorf("got %v peer %s, want %s", i, got.Address, want)
		}
	}
}

// expectPeersEventually validates that peers with addresses are connected with
// retires. It is supposed to be used to validate asynchronous connecting on the
// peer that is connected to.
func expectPeersEventually(t *testing.T, s *libp2p.Service, addrs ...swarm.Address) {
	t.Helper()

	var peers []p2p.Peer
	for i := 0; i < 100; i++ {
		peers = s.Peers()
		if len(peers) == len(addrs) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if len(peers) != len(addrs) {
		t.Fatalf("got peers %v, want %v", len(peers), len(addrs))
	}

	sort.Slice(addrs, func(i, j int) bool {
		return bytes.Compare(addrs[i].Bytes(), addrs[j].Bytes()) == -1
	})
	sort.Slice(peers, func(i, j int) bool {
		return bytes.Compare(peers[i].Address.Bytes(), peers[j].Address.Bytes()) == -1
	})

	for i, got := range peers {
		want := addrs[i]
		if !got.Address.Equal(want) {
			t.Errorf("got %v peer %s, want %s", i, got.Address, want)
		}
	}
}

func serviceUnderlayAddress(t *testing.T, s *libp2p.Service) multiaddr.Multiaddr {
	t.Helper()

	addrs, err := s.Addresses()
	if err != nil {
		t.Fatal(err)
	}
	return addrs[0]
}