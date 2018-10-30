package types

import "crypto/ecdsa"

type Account struct {
	privateKey *ecdsa.PrivateKey
}

func NewAccount(privKey *ecdsa.PrivateKey) *Account {
	return &Account{
		privateKey: privKey,
	}
}
