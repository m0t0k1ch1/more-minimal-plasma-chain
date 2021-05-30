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

	ErrBlockchainNotSynchronized = NewError(10001, "blockchain is not synchronized")

	ErrMempoolFull                    = NewError(11001, core.ErrMempoolFull.Error())
	ErrBlockNotFound                  = NewError(11002, core.ErrBlockNotFound.Error())
	ErrEmptyBlock                     = NewError(11003, core.ErrEmptyBlock.Error())
	ErrTxNotFound                     = NewError(11004, core.ErrTxNotFound.Error())
	ErrInvalidTxSignature             = NewError(11005, core.ErrInvalidTxSignature.Error())
	ErrInvalidTxConfirmationSignature = NewError(11006, core.ErrInvalidTxConfirmationSignature.Error())
	ErrInvalidTxBalance               = NewError(11007, core.ErrInvalidTxBalance.Error())
	ErrTxInNotFound                   = NewError(11008, core.ErrTxInNotFound.Error())
	ErrInvalidTxIn                    = NewError(11009, core.ErrInvalidTxIn.Error())
	ErrNullTxInConfirmation           = NewError(11010, core.ErrNullTxInConfirmation.Error())
	ErrTxOutNotFound                  = NewError(11011, core.ErrTxOutNotFound.Error())
	ErrTxOutAlreadySpent              = NewError(11012, core.ErrTxOutAlreadySpent.Error())
	ErrTxOutAlreadyExited             = NewError(11013, core.ErrTxOutAlreadyExited.Error())
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
