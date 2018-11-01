package app

import (
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (cc *ChildChain) GetChainHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	blkHashBytes, err := cc.blockchain.GetBlockHash(blkNum)
	if err != nil {
		if err == core.ErrBlockNotFound {
			return c.JSONError(ErrBlockNotFound)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blkhash": utils.EncodeToHex(blkHashBytes),
	})
}
