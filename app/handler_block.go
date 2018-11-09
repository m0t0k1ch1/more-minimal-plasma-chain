package app

import (
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

	// BEGIN TXN
	txn := p.db.NewTransaction(true)
	defer txn.Discard()

	currentBlkNum, err := p.childChain.GetCurrentBlockNumber(txn)
	if err != nil {
		return c.JSONError(err)
	}
	if rootBlkNum != currentBlkNum {
		return c.JSONError(ErrBlockchainNotSynchronized)
	}

	newBlkNum, err := p.childChain.AddBlock(txn, p.operator)
	if err != nil {
		if err == core.ErrEmptyBlock {
			return c.JSONError(ErrEmptyBlock)
		}
		return c.JSONError(err)
	}

	newBlk, err := p.childChain.GetBlock(txn, newBlkNum)
	if err != nil {
		if err == core.ErrBlockNotFound {
			return c.JSONError(ErrBlockNotFound)
		}
		return c.JSONError(err)
	}

	// COMMIT TXN
	if err := txn.Commit(nil); err != nil {
		return c.JSONError(err)
	}

	newBlkRootHash, err := newBlk.Root()
	if err != nil {
		return c.JSONError(err)
	}

	if _, err := p.rootChain.CommitPlasmaBlockRoot(p.operator, newBlkRootHash); err != nil {
		return c.JSONError(err)
	}
	p.Logger().Infof("[COMMIT] root: %s", utils.HashToHex(newBlkRootHash))

	return c.JSONSuccess(map[string]uint64{
		"blknum": newBlkNum,
	})
}

func (p *Plasma) GetBlockHandler(c *Context) error {
	blkNum, err := c.GetBlockNumberFromPath()
	if err != nil {
		return c.JSONError(err)
	}

	// BEGIN TXN
	txn := p.db.NewTransaction(false)
	defer txn.Discard()

	blk, err := p.childChain.GetBlock(txn, blkNum)
	if err != nil {
		if err == core.ErrBlockNotFound {
			return c.JSONError(ErrBlockNotFound)
		}
		return c.JSONError(err)
	}

	// COMMIT TXN
	if err := txn.Commit(nil); err != nil {
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
