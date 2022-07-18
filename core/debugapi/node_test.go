package debugapi_test

import (
	"testing"

	"github.com/redesblock/hop/core/debugapi"
)

func TestHopNodeMode_String(t *testing.T) {
	const nonExistingMode debugapi.HopNodeMode = 4

	mapping := map[string]string{
		debugapi.LightMode.String(): "light",
		debugapi.FullMode.String():  "full",
		debugapi.DevMode.String():   "dev",
		nonExistingMode.String():    "unknown",
	}

	for have, want := range mapping {
		if have != want {
			t.Fatalf("unexpected node mode: have %q; want %q", have, want)
		}
	}
}
