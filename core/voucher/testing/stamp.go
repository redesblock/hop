package testing

import (
	crand "crypto/rand"
	"io"

	"github.com/redesblock/hop/core/voucher"
)

const signatureSize = 65

// MustNewSignature will create a new random signature (65 byte slice). Panics
// on errors.
func MustNewSignature() []byte {
	sig := make([]byte, signatureSize)
	_, err := io.ReadFull(crand.Reader, sig)
	if err != nil {
		panic(err)
	}
	return sig
}

// MustNewStamp will generate a voucher stamp with random data. Panics on
// errors.
func MustNewStamp() *voucher.Stamp {
	return voucher.NewStamp(MustNewID(), MustNewID()[:8], MustNewID()[:8], MustNewSignature())
}
