package app

import "github.com/dgraph-io/badger"

type DB struct {
	*badger.DB
}

func NewDB(conf DBConfig) (*DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = conf.Dir
	opts.ValueDir = conf.Dir

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
