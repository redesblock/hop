package crypto_test

import (
	"bytes"
	"testing"

	"github.com/redesblock/hop/core/crypto"
)

func TestGenerateSecp256k1Key(t *testing.T) {
	k1, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	if k1 == nil {
		t.Fatal("nil key")
	}
	k2, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	if k2 == nil {
		t.Fatal("nil key")
	}

	if bytes.Equal(k1.D.Bytes(), k2.D.Bytes()) {
		t.Fatal("two generated keys are equal")
	}
}

func TestNewAddress(t *testing.T) {
	k, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	a := crypto.NewAddress(k.PublicKey)
	if l := len(a.Bytes()); l != 32 {
		t.Errorf("got address length %v, want %v", l, 32)
	}
}

func TestEncodeSecp256k1PrivateKey(t *testing.T) {
	k1, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		t.Fatal(err)
	}
	d := crypto.EncodeSecp256k1PrivateKey(k1)
	k2, err := crypto.DecodeSecp256k1PrivateKey(d)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(k1.D.Bytes(), k2.D.Bytes()) {
		t.Fatal("encoded and decoded keys are not equal")
	}
}
