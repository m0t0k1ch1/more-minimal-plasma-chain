package utils

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, i)
	return b[:n]
}

func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func BytesToUint64(b []byte) (uint64, error) {
	return binary.ReadUvarint(bytes.NewReader(b))
}

func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 64, 10)
}
