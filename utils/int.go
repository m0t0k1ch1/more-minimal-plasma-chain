package utils

import (
	"bytes"
	"encoding/binary"
)

func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, i)
	return b[:n]
}

func BytesToUint64(b []byte) (uint64, error) {
	return binary.ReadUvarint(bytes.NewReader(b))
}
