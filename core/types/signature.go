package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	b, err := hexutil.Decode(s)
	if err != nil {
		return NullSignature, err
	}

	return NewSignatureFromBytes(b)
}

func (sig Signature) Bytes() []byte {
	return sig[:]
}

func (sig Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(sig[:]).MarshalText()
}

func (sig Signature) SignerAddress(b []byte) (common.Address, error) {
	pubKey, err := crypto.SigToPub(b, sig.Bytes())
	if err != nil {
		return NullAddress, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
