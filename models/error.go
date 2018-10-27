package models

import "errors"

var (
	ErrMempoolFull = errors.New("mempool is full")
)
