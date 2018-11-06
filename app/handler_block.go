package app

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

func (p *Plasma) PostBlockHandler(c *Context) error {
	c.Request().ParseForm()

	blkNum, err := p.rootChain.CurrentPlasmaBlockNumber()
	if err != nil {
		return c.JSONError(err)
	}
	if blkNum.Cmp(p.childChain.CurrentBlockNumber()) != 0 {
		return c.JSONError(ErrBlockchainNotSynchronized)
	}

	blkHash, err := p.childChain.AddBlock(p.operator)
	if err != nil {
		if err == core.ErrEmptyBlock {
			return c.JSONError(ErrEmptyBlock)
		}
		return c.JSONError(err)
	}

	blk, err := p.childChain.GetBlock(blkHash)
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

	return c.JSONSuccess(map[string]interface{}{
		"blkhash": utils.HashToHex(blkHash),
	})
}

func (p *Plasma) GetBlockHandler(c *Context) error {
	blkHash, err := c.GetBlockHashFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	blk, err := p.childChain.GetBlock(blkHash)
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
