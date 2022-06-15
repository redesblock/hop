package clef

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/redesblock/hop/core/crypto"
	"github.com/redesblock/hop/core/crypto/eip712"
)

var (
	ErrNoAccounts       = errors.New("no accounts found in clef")
	clefRecoveryMessage = []byte("public key recovery message")
)

// ExternalSignerInterface is the interface for the clef client from go-ethereum.
type ExternalSignerInterface interface {
	SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error)
	SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	Accounts() []accounts.Account
}

// Client is the interface for rpc.RpcClient.
type Client interface {
	Call(result interface{}, method string, args ...interface{}) error
}

type clefSigner struct {
	client  Client // low-level rpc client to clef as ExternalSigner does not implement account_signTypedData
	clef    ExternalSignerInterface
	account accounts.Account // the account this signer will use
	pubKey  *ecdsa.PublicKey // the public key for the account
}

// DefaultIpcPath returns the os-dependent default ipc path for clef.
func DefaultIpcPath() (string, error) {
	socket := "clef.ipc"
	// on windows clef uses top level pipes
	if runtime.GOOS == "windows" {
		return `\\.\pipe\` + socket, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// on mac os clef defaults to ~/Library/Signer/clef.ipc
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Signer", socket), nil
	}

	// on unix clef defaults to ~/.clef/clef.ipc
	return filepath.Join(home, ".clef", socket), nil
}

// NewSigner creates a new connection to the signer at endpoint.
// As clef does not expose public keys it signs a test message to recover the public key.
func NewSigner(clef ExternalSignerInterface, client Client, recoverFunc crypto.RecoverFunc) (signer crypto.Signer, err error) {
	// get the list of available ethereum accounts
	clefAccounts := clef.Accounts()
	if len(clefAccounts) == 0 {
		return nil, ErrNoAccounts
	}

	// pick the first account as the one we use
	account := clefAccounts[0]

	// clef currently does not expose the public key
	// sign some data so we can recover it
	sig, err := clef.SignData(account, accounts.MimetypeTextPlain, clefRecoveryMessage)
	if err != nil {
		return nil, err
	}

	pubKey, err := recoverFunc(sig, clefRecoveryMessage)
	if err != nil {
		return nil, err
	}

	return &clefSigner{
		client:  client,
		clef:    clef,
		account: account,
		pubKey:  pubKey,
	}, nil
}

// PublicKey returns the public key recovered during creation.
func (c *clefSigner) PublicKey() (*ecdsa.PublicKey, error) {
	return c.pubKey, nil
}

// SignData signs with the text/plain type which is the standard Ethereum prefix method.
func (c *clefSigner) Sign(data []byte) ([]byte, error) {
	return c.clef.SignData(c.account, accounts.MimetypeTextPlain, data)
}

// SignTx signs an ethereum transaction.
func (c *clefSigner) SignTx(transaction *types.Transaction) (*types.Transaction, error) {
	// chainId is nil here because it is set on the clef side
	return c.clef.SignTx(c.account, transaction, nil)
}

// EthereumAddress returns the ethereum address this signer uses.
func (c *clefSigner) EthereumAddress() (common.Address, error) {
	return c.account.Address, nil
}

// SignTypedData signs data according to eip712.
func (c *clefSigner) SignTypedData(typedData *eip712.TypedData) ([]byte, error) {
	var sig hexutil.Bytes
	err := c.client.Call(&sig, "account_signTypedData", c.account.Address, typedData)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
