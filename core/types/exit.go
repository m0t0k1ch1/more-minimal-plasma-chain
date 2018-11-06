package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Exit struct {
	Owner     common.Address `json:"owner"`
	Amount    *big.Int       `json:"amount"`
	IsStarted bool           `json:"started"`
	IsValid   bool           `json:"valid"`
}
