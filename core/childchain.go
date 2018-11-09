package core

import (
	"bytes"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	DefaultBlockNumber = 1
)

type ChildChain struct{}

func NewChildChain(txn *badger.Txn) (*ChildChain, error) {
	cc := &ChildChain{}

	if _, err := cc.getCurrentBlockNumber(txn); err != nil {
		if err == badger.ErrKeyNotFound {
			if err := cc.setCurrentBlockNumber(txn, DefaultBlockNumber); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return cc, nil
}

func (cc *ChildChain) GetCurrentBlockNumber(txn *badger.Txn) (uint64, error) {
	return cc.getCurrentBlockNumber(txn)
}

func (cc *ChildChain) GetBlock(txn *badger.Txn, blkNum uint64) (*types.Block, error) {
	blk, err := cc.getBlock(txn, blkNum)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrBlockNotFound
		} else {
			return nil, err
		}
	}

	return blk, nil
}

func (cc *ChildChain) AddBlock(txn *badger.Txn, signer *types.Account) (uint64, error) {
	// get current block
	blk, err := cc.fixCurrentBlock(txn)
	if err != nil {
		return 0, err
	}

	// check block validity
	if len(blk.Txes) == 0 {
		return 0, ErrEmptyBlock
	}

	// sign block
	if err := blk.Sign(signer); err != nil {
		return 0, err
	}

	// add block
	if err := cc.addBlock(txn, blk); err != nil {
		return 0, err
	}

	// increment current block number
	if _, err := cc.incrementCurrentBlockNumber(txn); err != nil {
		return 0, err
	}

	return blk.Number, nil
}

func (cc *ChildChain) AddDepositBlock(txn *badger.Txn, ownerAddr common.Address, amount uint64, signer *types.Account) (uint64, error) {
	// create deposit tx
	tx := types.NewTx()
	txOut := types.NewTxOut(ownerAddr, amount)
	if err := tx.SetOutput(0, txOut); err != nil {
		return 0, err
	}

	// get current block number
	currentBlkNum, err := cc.getCurrentBlockNumber(txn)
	if err != nil {
		return 0, err
	}

	// create deposit block
	blk, err := types.NewBlock([]*types.Tx{tx}, currentBlkNum)
	if err != nil {
		return 0, err
	}

	// sign deposit block
	if err := blk.Sign(signer); err != nil {
		return 0, err
	}

	// add deposit block
	if err := cc.addBlock(txn, blk); err != nil {
		return 0, err
	}

	// increment current block number
	if _, err := cc.incrementCurrentBlockNumber(txn); err != nil {
		return 0, err
	}

	return blk.Number, nil
}

func (cc *ChildChain) GetTx(txn *badger.Txn, txPos types.Position) (*types.Tx, error) {
	blkNum, txIndex := types.ParseTxPosition(txPos)

	tx, err := cc.getTx(txn, blkNum, txIndex)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrTxNotFound
		} else {
			return nil, err
		}
	}

	return tx, nil
}

func (cc *ChildChain) GetTxProof(txn *badger.Txn, txPos types.Position) ([]byte, error) {
	blkNum, txIndex := types.ParseTxPosition(txPos)

	// get block
	blk, err := cc.getBlock(txn, blkNum)
	if err != nil {
		return nil, err
	}

	// build tx Merkle tree
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	// create tx proof
	return tree.CreateMembershipProof(txIndex)
}

func (cc *ChildChain) AddTxToMempool(txn *badger.Txn, tx *types.Tx) error {
	// check mempool capacity
	if cc.countTxesInMempool(txn) >= types.MaxBlockTxesNum {
		return types.ErrBlockTxesNumExceedsLimit
	}

	// validate tx
	if err := cc.validateTx(txn, tx); err != nil {
		return err
	}

	for _, txIn := range tx.Inputs {
		// skip if txin is null
		if txIn.IsNull() {
			continue
		}

		// get tx
		tx, err := cc.getTx(txn, txIn.BlockNumber, txIn.TxIndex)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrInvalidTxIn
			} else {
				return err
			}
		}

		// spend UTXO
		if err := tx.SpendOutput(txIn.OutputIndex); err != nil {
			if err == types.ErrInvalidTxOutIndex {
				return ErrInvalidTxIn
			} else {
				return err
			}
		}

		// update tx
		if err := cc.setTx(txn, txIn.BlockNumber, txIn.TxIndex, tx); err != nil {
			return err
		}
	}

	// add tx to mempool
	if err := cc.addTxToMempool(txn, tx); err != nil {
		return err
	}

	return nil
}

func (cc *ChildChain) ConfirmTx(txn *badger.Txn, txInPos types.Position, confSig types.Signature) error {
	blkNum, txIndex, inIndex := types.ParseTxInPosition(txInPos)

	// check tx existence
	tx, err := cc.getTx(txn, blkNum, txIndex)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrTxNotFound
		} else {
			return err
		}
	}

	// check txin existence
	if !tx.IsExistInput(inIndex) {
		return ErrTxInNotFound
	}

	// get txin
	txIn := tx.GetInput(inIndex)

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	// get input txout
	inTxOut, err := cc.getTxOut(txn, txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)
	if err != nil || inTxOut == nil {
		return ErrInvalidTxIn
	}

	// verify confirmation signature
	h, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}
	signerAddr, err := confSig.SignerAddress(h)
	if err != nil {
		return ErrInvalidTxConfirmationSignature
	}
	if !bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
		return ErrInvalidTxConfirmationSignature
	}

	// update confirmation signature
	if err := tx.SetConfirmationSignature(inIndex, confSig); err != nil {
		if err == types.ErrInvalidTxInIndex {
			return ErrInvalidTxIn
		} else {
			return err
		}
	}

	// update tx
	return cc.setTx(txn, blkNum, txIndex, tx)
}

func (cc *ChildChain) getCurrentBlockNumber(txn *badger.Txn) (uint64, error) {
	item, err := txn.Get([]byte("blknum_current"))
	if err != nil {
		return 0, err
	}

	blkNumBytes, err := item.Value()
	if err != nil {
		return 0, err
	}

	return utils.BytesToUint64(blkNumBytes)
}

func (cc *ChildChain) setCurrentBlockNumber(txn *badger.Txn, blkNum uint64) error {
	return txn.Set([]byte("blknum_current"), utils.Uint64ToBytes(blkNum))
}

func (cc *ChildChain) getNextBlockNumber(txn *badger.Txn) (uint64, error) {
	currentBlkNum, err := cc.getCurrentBlockNumber(txn)
	if err != nil {
		return 0, err
	}

	return currentBlkNum + 1, nil
}

func (cc *ChildChain) incrementCurrentBlockNumber(txn *badger.Txn) (uint64, error) {
	nextBlkNum, err := cc.getNextBlockNumber(txn)
	if err != nil {
		return 0, err
	}

	if err := cc.setCurrentBlockNumber(txn, nextBlkNum); err != nil {
		return 0, err
	}

	return nextBlkNum, nil
}

func (cc *ChildChain) setBlockHeader(txn *badger.Txn, blkNum uint64, blkHeader *types.BlockHeader) error {
	blkHeaderBytes, err := rlp.EncodeToBytes(blkHeader)
	if err != nil {
		return err
	}

	return txn.Set([]byte(fmt.Sprintf("blk_header_%d", blkNum)), blkHeaderBytes)
}

func (cc *ChildChain) getBlockHeader(txn *badger.Txn, blkNum uint64) (*types.BlockHeader, error) {
	item, err := txn.Get([]byte(fmt.Sprintf("blk_header_%d", blkNum)))
	if err != nil {
		return nil, err
	}

	blkHeaderBytes, err := item.Value()
	if err != nil {
		return nil, err
	}

	var blkHeader types.BlockHeader
	if err := rlp.DecodeBytes(blkHeaderBytes, &blkHeader); err != nil {
		return nil, err
	}

	return &blkHeader, nil
}

func (cc *ChildChain) getBlock(txn *badger.Txn, blkNum uint64) (*types.Block, error) {
	// get block header
	blkHeader, err := cc.getBlockHeader(txn, blkNum)
	if err != nil {
		return nil, err
	}

	blk := &types.Block{
		BlockHeader: blkHeader,
		Txes:        nil,
	}

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(fmt.Sprintf("tx_%d_", blk.Number))
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		// get tx
		var tx types.Tx
		txBytes, err := it.Item().Value()
		if err != nil {
			return nil, err
		}
		if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
			return nil, err
		}

		// add tx to block
		if err := blk.AddTx(&tx); err != nil {
			return nil, err
		}
	}

	return blk, nil
}

func (cc *ChildChain) fixCurrentBlock(txn *badger.Txn) (*types.Block, error) {
	// get current block number
	currentBlkNum, err := cc.getCurrentBlockNumber(txn)
	if err != nil {
		return nil, err
	}

	// create new block
	blk, err := types.NewBlock(nil, currentBlkNum)
	if err != nil {
		return nil, err
	}

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte("tx_mempool_")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()

		// get tx
		var tx types.Tx
		txBytes, err := item.Value()
		if err != nil {
			return nil, err
		}
		if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
			return nil, err
		}

		// add tx to block
		if err := blk.AddTx(&tx); err != nil {
			return nil, err
		}

		// remove tx from mempool
		if err := txn.Delete(item.Key()); err != nil {
			return nil, err
		}
	}

	return blk, nil
}

func (cc *ChildChain) addBlock(txn *badger.Txn, blk *types.Block) error {
	for i, tx := range blk.Txes {
		// store tx
		if err := cc.setTx(txn, blk.Number, uint64(i), tx); err != nil {
			return err
		}

		// TODO: store UTXO
	}

	// store block header
	return cc.setBlockHeader(txn, blk.Number, blk.BlockHeader)
}

func (cc *ChildChain) setTx(txn *badger.Txn, blkNum, txIndex uint64, tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	return txn.Set([]byte(fmt.Sprintf("tx_%d_%d", blkNum, txIndex)), txBytes)
}

func (cc *ChildChain) getTx(txn *badger.Txn, blkNum, txIndex uint64) (*types.Tx, error) {
	item, err := txn.Get([]byte(fmt.Sprintf("tx_%d_%d", blkNum, txIndex)))
	if err != nil {
		return nil, err
	}

	txBytes, err := item.Value()
	if err != nil {
		return nil, err
	}

	var tx types.Tx
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

func (cc *ChildChain) validateTx(txn *badger.Txn, tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := uint64(0), uint64(0)

	for _, txOut := range tx.Outputs {
		oAmount += txOut.Amount
	}

	for i, txIn := range tx.Inputs {
		// skip validation if txin is null (deposit)
		if txIn.IsNull() {
			nullTxInNum++
			continue
		}

		// get input txout
		inTxOut, err := cc.getTxOut(txn, txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)
		if err != nil {
			if err == badger.ErrKeyNotFound { // tx is not found
				return ErrInvalidTxIn
			} else {
				return err
			}
		} else if inTxOut == nil { // tx does not have the output
			return ErrInvalidTxIn
		}

		// check double spent
		if inTxOut.IsSpent {
			return ErrTxOutAlreadySpent
		}

		// verify signature
		signerAddr, err := tx.SignerAddress(uint64(i))
		if err != nil {
			return ErrInvalidTxSignature
		}
		if txIn.Signature == types.NullSignature ||
			!bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
			return ErrInvalidTxSignature
		}

		iAmount += inTxOut.Amount
	}

	// check txins validity
	if nullTxInNum == len(tx.Inputs) {
		return ErrInvalidTxIn
	}

	// check in/out balance
	if iAmount < oAmount {
		return ErrInvalidTxBalance
	}

	return nil
}

func (cc *ChildChain) countTxesInMempool(txn *badger.Txn) uint64 {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // key-only iteration

	it := txn.NewIterator(opts)
	defer it.Close()

	prefix, cnt := []byte("tx_mempool_"), uint64(0)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		cnt++
	}

	return cnt
}

func (cc *ChildChain) addTxToMempool(txn *badger.Txn, tx *types.Tx) error {
	txHash, err := tx.Hash()
	if err != nil {
		return err
	}

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	return txn.Set([]byte(fmt.Sprintf("tx_mempool_%s", utils.HashToHex(txHash))), txBytes)
}

func (cc *ChildChain) getTxOut(txn *badger.Txn, blkNum, txIndex, outIndex uint64) (*types.TxOut, error) {
	tx, err := cc.getTx(txn, blkNum, txIndex)
	if err != nil {
		return nil, err
	}

	return tx.GetOutput(outIndex), nil
}
