package core

import "errors"

var (
	ErrMempoolFull = errors.New("mempool is full")

	ErrBlockNotFound = errors.New("block is not found")
	ErrEmptyBlock    = errors.New("block is empty")

	ErrTxNotFound                     = errors.New("tx is not found")
	ErrInvalidTxSignature             = errors.New("tx signature is invalid")
	ErrInvalidTxConfirmationSignature = errors.New("tx confirmation signature is invalid")
	ErrInvalidTxBalance               = errors.New("tx balance is invalid")

	ErrTxInNotFound         = errors.New("txin is not found")
	ErrInvalidTxIn          = errors.New("txin is invalid")
	ErrNullTxInConfirmation = errors.New("null txin cannot be confirmed")

	ErrTxOutNotFound      = errors.New("txout is not found")
	ErrTxOutAlreadySpent  = errors.New("txout was already spent")
	ErrTxOutAlreadyExited = errors.New("txout was already exited")

	ErrNullConfirmationSignature = errors.New("confirmation signature is null")
)
