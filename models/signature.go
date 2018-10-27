package models

import "github.com/ethereum/go-ethereum/common/hexutil"

type Signature [SignatureSize]byte

func newSignatureFromBytes(b []byte) Signature {
	sig := Signature{}
	copy(sig[:], b[:])

	return sig
}

func (sig Signature) Bytes() []byte {
	return sig[:]
}

func (sig Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(sig[:]).MarshalText()
}
