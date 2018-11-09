package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

/*
blknum_current                   => uint64
blk_header_<block number>        => *types.BlockHeader
tx_<block_number>_<tx index>     => *types.Tx
tx_mempool_<tx hash>             => *types.Tx
token_<address>_<txout position> => types.Position
*/

const (
	FirstBlockNumber = 1
	MempoolSize      = 99999 // must be less than or equal to types.MaxBlockTxesNum

	currentBlockNumberKey = "current_blknum"
	blockHeaderKeyPrefix  = "blk_header"
	txKeyPrefix           = "tx"
	mempoolTxKeyPrefix    = "mempool_tx"
	tokenKeyPrefix        = "token"
)

type ChildChain struct{}

func NewChildChain(txn *badger.Txn) (*ChildChain, error) {
	cc := &ChildChain{}

	if _, err := cc.getCurrentBlockNumber(txn); err != nil {
		if err == badger.ErrKeyNotFound {
			if err := cc.setCurrentBlockNumber(txn, FirstBlockNumber); err != nil {
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

	// check tx existence
	if _, err := cc.getTx(txn, blkNum, txIndex); err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrTxNotFound
		} else {
			return nil, err
		}
	}

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
	if cc.countTxesInMempool(txn) >= MempoolSize {
		return ErrMempoolFull
	}

	// validate tx
	if err := cc.ValidateTx(txn, tx); err != nil {
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

		// spend txout
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

func (cc *ChildChain) ValidateTx(txn *badger.Txn, tx *types.Tx) error {
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

func (cc *ChildChain) GetUTXOPositions(txn *badger.Txn, addr common.Address) ([]types.Position, error) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // key-only iteration

	it := txn.NewIterator(opts)
	defer it.Close()

	prefix, poses := cc.tokenKeyPrefix(addr), []types.Position{}
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		// get position
		pos, err := types.StrToPosition(strings.TrimPrefix(string(it.Item().Key()), string(prefix)))
		if err != nil {
			return nil, err
		}

		blkNum, txIndex, outIndex := types.ParseTxOutPosition(pos)

		// get txout
		txOut, err := cc.getTxOut(txn, blkNum, txIndex, outIndex)
		if err != nil {
			return nil, err
		}

		// skip if txout was spent
		if txOut.IsSpent {
			continue
		}

		poses = append(poses, pos)
	}

	return poses, nil
}

func (cc *ChildChain) currentBlockNumberKey() []byte {
	return []byte(currentBlockNumberKey)
}

func (cc *ChildChain) getCurrentBlockNumber(txn *badger.Txn) (uint64, error) {
	item, err := txn.Get(cc.currentBlockNumberKey())
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
	return txn.Set(cc.currentBlockNumberKey(), utils.Uint64ToBytes(blkNum))
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

func (cc *ChildChain) blockHeaderKey(blkNum uint64) []byte {
	return []byte(fmt.Sprintf("%s_%d", blockHeaderKeyPrefix, blkNum))
}

func (cc *ChildChain) setBlockHeader(txn *badger.Txn, blkNum uint64, blkHeader *types.BlockHeader) error {
	blkHeaderBytes, err := rlp.EncodeToBytes(blkHeader)
	if err != nil {
		return err
	}

	return txn.Set(cc.blockHeaderKey(blkNum), blkHeaderBytes)
}

func (cc *ChildChain) getBlockHeader(txn *badger.Txn, blkNum uint64) (*types.BlockHeader, error) {
	item, err := txn.Get(cc.blockHeaderKey(blkNum))
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

	// convert block header to block
	blk := &types.Block{
		BlockHeader: blkHeader,
		Txes:        nil,
	}

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := cc.txKeyPrefix(blkNum)
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

	prefix := cc.mempoolTxKeyPrefix()
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
		for j, txIn := range tx.Inputs {
			if txIn.IsNull() {
				continue
			}

			inTxOut, err := cc.getTxOut(txn, txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)
			if err != nil {
				return err
			}

			// update position of txin by which token was spent
			if err := cc.setToken(txn,
				inTxOut.OwnerAddress,
				types.NewTxOutPosition(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex),
				types.NewTxInPosition(blk.Number, uint64(i), uint64(j)),
			); err != nil {
				return err
			}
		}

		for j, txOut := range tx.Outputs {
			// store unspent token
			if err := cc.setToken(
				txn,
				txOut.OwnerAddress,
				types.NewTxOutPosition(blk.Number, uint64(i), uint64(j)),
				0,
			); err != nil {
				return err
			}
		}

		// store tx
		if err := cc.setTx(txn, blk.Number, uint64(i), tx); err != nil {
			return err
		}
	}

	// store block header
	return cc.setBlockHeader(txn, blk.Number, blk.BlockHeader)
}

func (cc *ChildChain) txKeyPrefix(blkNum uint64) []byte {
	return []byte(fmt.Sprintf("%s_%d_", txKeyPrefix, blkNum))
}

func (cc *ChildChain) txKey(blkNum, txIndex uint64) []byte {
	return []byte(fmt.Sprintf("%s_%d_%d", txKeyPrefix, blkNum, txIndex))
}

func (cc *ChildChain) setTx(txn *badger.Txn, blkNum, txIndex uint64, tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	return txn.Set(cc.txKey(blkNum, txIndex), txBytes)
}

func (cc *ChildChain) getTx(txn *badger.Txn, blkNum, txIndex uint64) (*types.Tx, error) {
	item, err := txn.Get(cc.txKey(blkNum, txIndex))
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

func (cc *ChildChain) mempoolTxKeyPrefix() []byte {
	return []byte(fmt.Sprintf("%s_", mempoolTxKeyPrefix))
}

func (cc *ChildChain) mempoolTxKey(tx *types.Tx) ([]byte, error) {
	txHash, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%s_%s", mempoolTxKeyPrefix, utils.HashToHex(txHash))), nil
}

func (cc *ChildChain) countTxesInMempool(txn *badger.Txn) uint64 {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // key-only iteration

	it := txn.NewIterator(opts)
	defer it.Close()

	prefix, cnt := cc.mempoolTxKeyPrefix(), uint64(0)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		cnt++
	}

	return cnt
}

func (cc *ChildChain) addTxToMempool(txn *badger.Txn, tx *types.Tx) error {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	key, err := cc.mempoolTxKey(tx)
	if err != nil {
		return err
	}

	return txn.Set(key, txBytes)
}

func (cc *ChildChain) getTxOut(txn *badger.Txn, blkNum, txIndex, outIndex uint64) (*types.TxOut, error) {
	tx, err := cc.getTx(txn, blkNum, txIndex)
	if err != nil {
		return nil, err
	}

	return tx.GetOutput(outIndex), nil
}

func (cc *ChildChain) tokenKeyPrefix(addr common.Address) []byte {
	return []byte(fmt.Sprintf("%s_%s_", tokenKeyPrefix, utils.AddressToHex(addr)))
}

func (cc *ChildChain) tokenKey(addr common.Address, txOutPos types.Position) []byte {
	return []byte(fmt.Sprintf("%s_%s_%d", tokenKeyPrefix, utils.AddressToHex(addr), txOutPos))
}

func (cc *ChildChain) setToken(txn *badger.Txn, addr common.Address, txOutPos types.Position, spendingTxInPos types.Position) error {
	return txn.Set(cc.tokenKey(addr, txOutPos), spendingTxInPos.Bytes())
}
