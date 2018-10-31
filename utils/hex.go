package utils

import "github.com/ethereum/go-ethereum/common/hexutil"

func DecodeHex(s string) ([]byte, error) {
	return hexutil.Decode(s)
}

func EncodeToHex(b []byte) string {
	return hexutil.Encode(b)
}
