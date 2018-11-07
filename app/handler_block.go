package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (p *Plasma) PostBlockHandler(c *Context) error {
	c.Request().ParseForm()

	rootBlkNum, err := p.rootChain.CurrentPlasmaBlockNumber()
	if err != nil {
		return c.JSONError(err)
	}
	if rootBlkNum.Cmp(p.childChain.CurrentBlockNumber()) != 0 {
		return c.JSONError(ErrBlockchainNotSynchronized)
	}

	blkNum, err := p.childChain.AddBlock(p.operator)
	if err != nil {
		if err == core.ErrEmptyBlock {
			return c.JSONError(ErrEmptyBlock)
		}
		return c.JSONError(err)
	}

	blk, err := p.childChain.GetBlock(blkNum)
	if err != nil {
		if err == core.ErrBlockNotFound {
			return c.JSONError(ErrBlockNotFound)
		}
		return c.JSONError(err)
	}

	blkRootHash, err := blk.Root()
	if err != nil {
		return c.JSONError(err)
	}

	if _, err := p.rootChain.CommitPlasmaBlockRoot(p.operator, blkRootHash); err != nil {
		return c.JSONError(err)
	}
	p.Logger().Infof("[COMMIT] root: %s", utils.HashToHex(blkRootHash))

	return c.JSONSuccess(map[string]*big.Int{
		"blknum": blkNum,
	})
}

func (p *Plasma) GetBlockHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	blk, err := p.childChain.GetBlock(blkNum)
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

	return c.JSONSuccess(map[string]string{
		"blk": utils.EncodeToHex(blkBytes),
	})
}
