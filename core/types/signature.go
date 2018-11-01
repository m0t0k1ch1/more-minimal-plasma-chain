package types

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	SignatureSize = 65 // bytes
)

var (
	NullSignature = Signature{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

var (
	ErrInvalidSignatureSize = fmt.Errorf("signature must be %d bytes", SignatureSize)
)

type Signature [SignatureSize]byte

func NewSignatureFromBytes(b []byte) (Signature, error) {
	if len(b) != SignatureSize {
		return NullSignature, ErrInvalidSignatureSize
	}

	sig := Signature{}
	copy(sig[:], b[:])

	return sig, nil
}

func NewSignatureFromHex(s string) (Signature, error) {
	b, err := utils.DecodeHex(s)
	if err != nil {
		return NullSignature, err
	}

	return NewSignatureFromBytes(b)
}

func (sig Signature) Bytes() []byte {
	return sig[:]
}

func (sig Signature) MarshalText() ([]byte, error) {
	b := make([]byte, len(sig[:])*2+2)
	copy(b, `0x`)
	hex.Encode(b[2:], sig[:])

	return b, nil
}

func (sig Signature) SignerAddress(b []byte) (common.Address, error) {
	pubKey, err := crypto.SigToPub(b, sig.Bytes())
	if err != nil {
		return NullAddress, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
