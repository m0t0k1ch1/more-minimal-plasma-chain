package app

import (
	"fmt"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
)

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
		3001,
		"block type is invalid",
	)
	ErrOwnerRequired = NewError(
		3002,
		"'owner' param is required",
	)
	ErrInvalidOwnerHex = NewError(
		3003,
		"owner hex is invalid",
	)
	ErrAmountRequired = NewError(
		3004,
		"'amount' param is required",
	)
	ErrInvalidAmount = NewError(
		3005,
		"amount is invalid",
	)
	ErrTxRequired = NewError(
		3006,
		"'tx' param is required",
	)
	ErrInvalidTxHex = NewError(
		3007,
		"tx hex is invalid",
	)

	ErrBlockNotFound = NewError(
		4001,
		"block is not found",
	)
	ErrTxNotFound = NewError(
		4002,
		"tx is not found",
	)

	ErrInvalidTxInput = NewError(
		5001,
		core.ErrInvalidTxInput.Error(),
	)
	ErrTxInputAlreadySpent = NewError(
		5002,
		core.ErrTxInputAlreadySpent.Error(),
	)
	ErrInvalidTxSignature = NewError(
		5003,
		core.ErrInvalidTxSignature.Error(),
	)
	ErrInvalidTxBalance = NewError(
		5004,
		core.ErrInvalidTxBalance.Error(),
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
