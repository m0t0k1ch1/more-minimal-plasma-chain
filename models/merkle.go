package models

import (
	"sync"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	merkle "github.com/m0t0k1ch1/fixed-merkle"
)

const (
	merkleTreeDepth = 10
	merkleLeafSize  = 32
)

var (
	initonce     sync.Once
	merkleConfig *merkle.Config
)

func initMerkleConfig() {
	merkleConfig, _ = merkle.NewConfig(
		sha3.NewKeccak256(),
		merkleTreeDepth,
		merkleLeafSize,
	)
}

func MerkleConfig() *merkle.Config {
	initonce.Do(initMerkleConfig)
	return merkleConfig
}
