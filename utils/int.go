package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"strconv"
)

var (
	ErrInvalidDecimalString = errors.New("invalid decimal string")
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
	return strconv.ParseUint(s, 10, 64)
}

func StringToBigInt(s string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, ErrInvalidDecimalString
	}
	return i, nil
}
