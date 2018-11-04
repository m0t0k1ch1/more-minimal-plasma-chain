package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	NullHash = common.BytesToHash([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	})
)

var (
	ErrInvalidHashSize = fmt.Errorf("hash size must be %d bytes", common.HashLength)
)

func HexToHash(b []byte) (common.Hash, error) {
	if len(b) != common.HashLength {
		return NullHash, ErrInvalidHashSize
	}

	h := common.Hash{}
	copy(h[:], b[:])

	return h, nil
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
