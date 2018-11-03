package types

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	privateKey *ecdsa.PrivateKey
}

func NewAccount(privKey *ecdsa.PrivateKey) *Account {
	return &Account{
		privateKey: privKey,
	}
}

func (a *Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.privateKey.PublicKey)
}

func (a *Account) TransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(a.privateKey)
}

func (a *Account) Sign(b []byte) ([]byte, error) {
	return crypto.Sign(b, a.privateKey)
}
