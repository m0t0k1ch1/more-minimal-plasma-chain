package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	DefaultBlockNumber = 1
)

type ChildChain struct {
	mu    *sync.RWMutex
	chain map[uint64]*types.Block // key: blkNum
}

func NewChildChain(txn *badger.Txn) (*ChildChain, error) {
	cc := &ChildChain{
		mu:    &sync.RWMutex{},
		chain: map[uint64]*types.Block{},
	}

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

func (cc *ChildChain) GetBlock(blkNum uint64) (*types.Block, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistBlock(blkNum) {
		return nil, ErrBlockNotFound
	}

	return cc.getBlock(blkNum), nil
}

func (cc *ChildChain) AddBlock(txn *badger.Txn, signer *types.Account) (uint64, error) {
	// get current block
	blk, err := cc.createNewBlock(txn)
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

func (cc *ChildChain) GetTx(txPos types.Position) (*types.Tx, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkNum, txIndex := types.ParseTxPosition(txPos)

	if !cc.isExistTx(blkNum, txIndex) {
		return nil, ErrTxNotFound
	}

	return cc.getTx(blkNum, txIndex), nil
}

func (cc *ChildChain) GetTxProof(txPos types.Position) ([]byte, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkNum, txIndex := types.ParseTxPosition(txPos)

	blk := cc.getBlock(blkNum)

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
	if err := cc.validateTx(tx); err != nil {
		return err
	}

	// spend input utxos
	for _, txIn := range tx.Inputs {
		if txIn.IsNull() {
			continue
		}
		if err := cc.spendUTXO(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex); err != nil {
			return err
		}
	}

	// add tx to mempool
	if err := cc.addTxToMempool(txn, tx); err != nil {
		return err
	}

	return nil
}

func (cc *ChildChain) ConfirmTx(txInPos types.Position, confSig types.Signature) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	blkNum, txIndex, inIndex := types.ParseTxInPosition(txInPos)

	// check tx existence
	if !cc.isExistTx(blkNum, txIndex) {
		return ErrTxNotFound
	}

	tx := cc.getTx(blkNum, txIndex)

	// check txin existence
	if !tx.IsExistInput(inIndex) {
		return ErrTxInNotFound
	}

	txIn := tx.GetInput(inIndex)

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	inTxOut := cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

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
	if err := cc.setConfirmationSignature(blkNum, txIndex, inIndex, confSig); err != nil {
		return err
	}

	return nil
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

func (cc *ChildChain) createNewBlock(txn *badger.Txn) (*types.Block, error) {
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

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

	prefix := []byte("tx_mempool_")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		val, err := item.Value()
		if err != nil {
			return nil, err
		}

		var tx types.Tx
		if err := rlp.DecodeBytes(val, &tx); err != nil {
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
	// TODO: delete
	cc.chain[blk.Number] = blk

	// TODO: save txes
	// TODO: save utxos

	// TODO: remove txes from block
	blkBytes, err := rlp.EncodeToBytes(blk)
	if err != nil {
		return err
	}

	return txn.Set([]byte(fmt.Sprintf("blk_%d", blk.Number)), blkBytes)
}

func (cc *ChildChain) countTxesInMempool(txn *badger.Txn) uint64 {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // key-only iteration

	it := txn.NewIterator(opts)
	defer it.Close()

	cnt := uint64(0)

	prefix := []byte("tx_mempool_")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		cnt++
	}

	return cnt
}

func (cc *ChildChain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := uint64(0), uint64(0)

	for _, txOut := range tx.Outputs {
		oAmount += txOut.Amount
	}

	for i, txIn := range tx.Inputs {
		// check spending txout existence
		if !cc.isExistTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex) {
			if txIn.IsNull() {
				nullTxInNum++
				continue
			}
			return ErrInvalidTxIn
		}

		inTxOut := cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

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

func (cc *ChildChain) getBlock(blkNum uint64) *types.Block {
	return cc.chain[blkNum]
}

func (cc *ChildChain) isExistBlock(blkNum uint64) bool {
	_, ok := cc.chain[blkNum]
	return ok
}

func (cc *ChildChain) getTx(blkNum, txIndex uint64) *types.Tx {
	return cc.getBlock(blkNum).GetTx(txIndex)
}

func (cc *ChildChain) isExistTx(blkNum, txIndex uint64) bool {
	if !cc.isExistBlock(blkNum) {
		return false
	}

	return cc.getBlock(blkNum).IsExistTx(txIndex)
}

func (cc *ChildChain) getTxOut(blkNum, txIndex, outIndex uint64) *types.TxOut {
	return cc.getTx(blkNum, txIndex).GetOutput(outIndex)
}

func (cc *ChildChain) isExistTxOut(blkNum, txIndex, outIndex uint64) bool {
	if !cc.isExistTx(blkNum, txIndex) {
		return false
	}

	return cc.getTx(blkNum, txIndex).IsExistOutput(outIndex)
}

func (cc *ChildChain) spendUTXO(blkNum, txIndex, outIndex uint64) error {
	return cc.getTx(blkNum, txIndex).SpendOutput(outIndex)
}

func (cc *ChildChain) setConfirmationSignature(blkNum, txIndex, inIndex uint64, confSig types.Signature) error {
	return cc.getTx(blkNum, txIndex).SetConfirmationSignature(inIndex, confSig)
}
