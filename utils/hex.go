package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func HashToHex(h common.Hash) string {
	return EncodeToHex(h.Bytes())
}

func DecodeHex(s string) ([]byte, error) {
	return hexutil.Decode(s)
}

func EncodeToHex(b []byte) string {
	return hexutil.Encode(b)
}
