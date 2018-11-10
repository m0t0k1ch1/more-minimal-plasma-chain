package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func IsHexAddress(s string) bool {
	return common.IsHexAddress(s)
}

func IsHexHash(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*common.HashLength && isHex(s)
}

func BytesToAddress(b []byte) common.Address {
	return common.BytesToAddress(b)
}

func BytesToHash(b []byte) common.Hash {
	return common.BytesToHash(b)
}

func HexToAddress(s string) common.Address {
	return common.HexToAddress(s)
}

func HexToHash(s string) common.Hash {
	return common.HexToHash(s)
}

func HexToPrivateKey(s string) (*ecdsa.PrivateKey, error) {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return crypto.HexToECDSA(s)
}

func AddressToHex(addr common.Address) string {
	return EncodeToHex(addr.Bytes())
}

func HashToHex(h common.Hash) string {
	return EncodeToHex(h.Bytes())
}

func DecodeHex(s string) ([]byte, error) {
	return hexutil.Decode(s)
}

func EncodeToHex(b []byte) string {
	return hexutil.Encode(b)
}

func hasHexPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isHex(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for _, c := range []byte(s) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
