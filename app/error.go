package app

import (
	"fmt"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
)

const (
	PathParamErrorCode = 20001
	FormParamErrorCode = 20002
)

var (
	ErrUnexpected = NewError(10000, "unexpected error")

	ErrBlockNotFound = NewError(11001, "block is not found")
	ErrTxNotFound    = NewError(11002, "tx is not found")

	ErrInvalidTxInput      = NewError(12001, core.ErrInvalidTxInput.Error())
	ErrTxInputAlreadySpent = NewError(12002, core.ErrTxInputAlreadySpent.Error())
	ErrInvalidTxSignature  = NewError(12003, core.ErrInvalidTxSignature.Error())
	ErrInvalidTxBalance    = NewError(12004, core.ErrInvalidTxBalance.Error())
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

func NewInvalidPathParamError(key string) *Error {
	return NewError(
		PathParamErrorCode,
		fmt.Sprintf("'%s' is invalid", key),
	)
}

func NewRequiredFormParamError(key string) *Error {
	return NewError(
		FormParamErrorCode,
		fmt.Sprintf("'%s' is required", key),
	)
}

func NewInvalidFormParamError(key string) *Error {
	return NewError(
		FormParamErrorCode,
		fmt.Sprintf("'%s' is invalid", key),
	)
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s [%d]", err.Message, err.Code)
}
