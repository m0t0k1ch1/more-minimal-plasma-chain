package app

import "fmt"

var (
	ErrUnexpected = NewError(
		1000,
		"unexpected error",
	)

	ErrInvalidAddressHex = NewError(
		2001,
		"invalid address hex",
	)
	ErrInvalidBlockNumber = NewError(
		2002,
		"invalid block number",
	)
	ErrInvalidTxIndex = NewError(
		2003,
		"invalid tx index",
	)

	ErrBlockNotFound = NewError(
		3001,
		"block not found",
	)
	ErrTxNotFound = NewError(
		3002,
		"tx not found",
	)
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s [%d]", err.Message, err.Code)
}
