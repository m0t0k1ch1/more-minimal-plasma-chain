package app

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
)

func (cc *ChildChain) GetBlockHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	return cc.getBlockHandler(c, blkNum)
}

func (cc *ChildChain) PostBlockHandler(c *Context) error {
	c.Request().ParseForm()

	bt, err := c.GetBlockTypeFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	// TODO: delete
	if bt.IsDeposit() {
		return cc.postDepositBlockHandler(c)
	}

	return cc.postBlockHandler(c)
}

func (cc *ChildChain) postBlockHandler(c *Context) error {
	blkNum, err := cc.blockchain.AddBlock(cc.operator)
	if err != nil {
		if err == core.ErrEmptyBlock {
			return c.JSONError(ErrEmptyBlock)
		}
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blknum": blkNum,
	})
}

// TODO: delete
func (cc *ChildChain) postDepositBlockHandler(c *Context) error {
	ownerAddr, err := c.GetOwnerFromForm()
	if err != nil {
		return c.JSONError(err)
	}
	amount, err := c.GetAmountFromForm()
	if err != nil {
		return c.JSONError(err)
	}

	blkNum, err := cc.blockchain.AddDepositBlock(ownerAddr, amount, cc.operator)
	if err != nil {
		return c.JSONError(err)
	}

	return c.JSONSuccess(map[string]interface{}{
		"blknum": blkNum,
	})
}

func (cc *ChildChain) getBlockHandler(c *Context, blkNum uint64) error {
	blk, err := cc.blockchain.GetBlock(blkNum)
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
		"hex": hexutil.Encode(blkBytes),
	})
}
