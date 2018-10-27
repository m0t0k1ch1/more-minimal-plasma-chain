package main

import "fmt"

var (
	ErrUnexpected = NewError(
		1000,
		"unexpected error",
	)

	ErrInvalidAddressHex = NewError(
		2001,
		"address is invalid hex",
	)

	ErrBlockNotFound = NewError(
		3001,
		"block is not found",
	)
	ErrTxNotFound = NewError(
		3002,
		"tx is not found",
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
