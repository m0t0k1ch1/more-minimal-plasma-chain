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

func BytesToSignature(b []byte) (Signature, error) {
	if len(b) != SignatureSize {
		return NullSignature, ErrInvalidSignatureSize
	}

	sig := Signature{}
	copy(sig[:], b)

	return sig, nil
}

func HexToSignature(s string) (Signature, error) {
	b, err := utils.DecodeHex(s)
	if err != nil {
		return NullSignature, err
	}

	return BytesToSignature(b)
}

func (sig Signature) Bytes() []byte {
	return sig[:]
}

func (sig Signature) Hex() string {
	return utils.EncodeToHex(sig.Bytes())
}

func (sig Signature) MarshalText() ([]byte, error) {
	b := make([]byte, len(sig[:])*2+2)
	copy(b, `0x`)
	hex.Encode(b[2:], sig[:])

	return b, nil
}

func (sig Signature) SignerAddress(h common.Hash) (common.Address, error) {
	pubKey, err := crypto.SigToPub(h.Bytes(), sig.Bytes())
	if err != nil {
		return NullAddress, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
