package skippeers_test

import (
	"testing"

	"github.com/redesblock/hop/core/skippeers"
	"github.com/redesblock/hop/core/swarm"
)

func TestAddOverdraft(t *testing.T) {
	var (
		p1 = swarm.NewAddress([]byte("0xab"))
		p2 = swarm.NewAddress([]byte("0xbc"))
	)

	sp := new(skippeers.List)
	sp.Add(p1)

	// duplicate entries are ignored
	sp.Add(p1)
	if len(sp.All()) != 1 {
		t.Errorf("expected len: %d, got %d", 1, len(sp.All()))
	}

	// add peer
	sp.Add(p2)
	if len(sp.All()) != 2 {
		t.Errorf("expected len: %d, got %d", 2, len(sp.All()))
	}

	// add overdraft removes from addresses
	sp.AddOverdraft(p2)
	if len(sp.All()) != 2 {
		t.Errorf("expected len: %d, got %d", 2, len(sp.All()))
	}

	sp.ResetOverdraft()

	if !sp.OverdraftListEmpty() {
		t.Errorf("expected empty list, got %s", sp.All())
	}
}
