package file_test

import (
	"testing"

	"github.com/redesblock/hop/core/keystore/file"
	"github.com/redesblock/hop/core/keystore/test"
)

func TestService(t *testing.T) {
	dir := t.TempDir()

	test.Service(t, file.New(dir))
}
