package app

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	c.Request().ParseForm()

	bt, err := c.GetBlockTypeFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	if bt.IsDeposit() {
		return cc.postDepositBlockHandler(c)
	}

	return cc.postBlockHandler(c)
}

func (cc *ChildChain) postBlockHandler(c *Context) error {
	blkHashBytes, err := cc.blockchain.AddBlock(cc.operator)
	if err != nil {
		if err == core.ErrEmptyBlock {
			return c.JSONError(ErrEmptyBlock)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blkhash": utils.EncodeToHex(blkHashBytes),
	})
}

func (cc *ChildChain) postDepositBlockHandler(c *Context) error {
	ownerAddr, err := c.GetOwnerFromForm()
	if err != nil {
		return c.JSONError(err)
	}
	amount, err := c.GetAmountFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	blkHashBytes, err := cc.blockchain.AddDepositBlock(ownerAddr, amount, cc.operator)
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blkhash": utils.EncodeToHex(blkHashBytes),
	})
}

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	blkHashBytes, err := c.GetBlockHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	blk, err := cc.blockchain.GetBlock(blkHashBytes)
	if err != nil {
		if err == core.ErrBlockNotFound {
			return c.JSONError(ErrBlockNotFound)
		}
		return c.JSONError(err)
	}

	blkBytes, err := rlp.EncodeToBytes(blk)
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blk": utils.EncodeToHex(blkBytes),
	})
}
