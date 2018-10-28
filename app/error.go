package app

import "fmt"

var (
	ErrUnexpected = NewError(
		1000,
		"unexpected error",
	)

	ErrInvalidBlockNumber = NewError(
		2001,
		"block number is invalid",
	)
	ErrInvalidTxIndex = NewError(
		2002,
		"tx index is invalid",
	)
	ErrInvalidBlockType = NewError(
		2101,
		"block type is invalid",
	)
	ErrOwnerRequired = NewError(
		2102,
		"'owner' param is required",
	)
	ErrInvalidOwnerHex = NewError(
		2103,
		"owner hex is invalid",
	)
	ErrAmountRequired = NewError(
		2104,
		"'amount' param is required",
	)
	ErrInvalidAmount = NewError(
		2105,
		"amount is invalid",
	)
	ErrTxRequired = NewError(
		2106,
		"'tx' param is required",
	)
	ErrInvalidTxHex = NewError(
		2107,
		"tx hex is invalid",
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
