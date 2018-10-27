package models

import "errors"

var (
	ErrMempoolFull        = errors.New("mempool is full")
	ErrBlockAlreadyExists = errors.New("block already exists")
)
