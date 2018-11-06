package core

import (
	"bytes"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	DefaultBlockNumber = 1
)

var (
	ErrBlockNotFound = errors.New("block is not found")
	ErrEmptyBlock    = errors.New("block is empty")

	ErrTxNotFound                     = errors.New("tx is not found")
	ErrInvalidTxSignature             = errors.New("tx signature is invalid")
	ErrInvalidTxConfirmationSignature = errors.New("tx confirmation signature is invalid")
	ErrInvalidTxBalance               = errors.New("tx balance is invalid")

	ErrTxInNotFound         = errors.New("txin is not found")
	ErrInvalidTxIn          = errors.New("txin is invalid")
	ErrNullTxInConfirmation = errors.New("null txin cannot be confirmed")

	ErrTxOutAlreadySpent = errors.New("txout is already spent")
)

type ChildChain struct {
	mu           *sync.RWMutex
	currentBlock *types.Block
	chain        map[string]common.Hash
	lightBlocks  map[string]*types.LightBlock
	blockTxes    map[string]*types.BlockTx
}

func NewChildChain() (*ChildChain, error) {
	blk, err := types.NewBlock(nil, big.NewInt(DefaultBlockNumber))
	if err != nil {
		return nil, err
	}

	return &ChildChain{
		mu:           &sync.RWMutex{},
		currentBlock: blk,
		chain:        map[string]common.Hash{},
		lightBlocks:  map[string]*types.LightBlock{},
		blockTxes:    map[string]*types.BlockTx{},
	}, nil
}

func (cc *ChildChain) CurrentBlockNumber() *big.Int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.currentBlockNumber()
}

func (cc *ChildChain) GetBlockHash(blkNum *big.Int) (common.Hash, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistBlockHash(blkNum) {
		return types.NullHash, ErrBlockNotFound
	}

	return cc.getBlockHash(blkNum), nil
}

func (cc *ChildChain) GetBlock(blkHash common.Hash) (*types.Block, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistLightBlock(blkHash) {
		return nil, ErrBlockNotFound
	}

	return cc.getBlock(blkHash)
}

func (cc *ChildChain) AddBlock(signer *types.Account) (common.Hash, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	blk := cc.currentBlock

	// check block validity
	if len(blk.Txes) == 0 {
		return types.NullHash, ErrEmptyBlock
	}

	// sign block
	if err := blk.Sign(signer); err != nil {
		return types.NullHash, err
	}

	// add block
	blkHash, err := cc.addBlock(blk)
	if err != nil {
		return types.NullHash, err
	}

	// reset current block
	blkNext, err := types.NewBlock(nil, cc.newNextBlockNumber())
	if err != nil {
		return types.NullHash, err
	}
	cc.currentBlock = blkNext

	return blkHash, nil
}

func (cc *ChildChain) AddDepositBlock(ownerAddr common.Address, amount *big.Int, signer *types.Account) (common.Hash, common.Hash, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// create deposit tx
	tx := types.NewTx()
	txOut := types.NewTxOut(ownerAddr, amount)
	if err := tx.SetOutput(big.NewInt(0), txOut); err != nil {
		return types.NullHash, types.NullHash, err
	}

	// create deposit block
	blk, err := types.NewBlock([]*types.Tx{tx}, cc.newCurrentBlockNumber())
	if err != nil {
		return types.NullHash, types.NullHash, err
	}

	// sign deposit block
	if err := blk.Sign(signer); err != nil {
		return types.NullHash, types.NullHash, err
	}

	// add deposit block
	blkHash, err := cc.addBlock(blk)
	if err != nil {
		return types.NullHash, types.NullHash, err
	}

	txHash := cc.getTxHash(cc.currentBlockNumber(), big.NewInt(0))

	// increment current block number
	cc.incrementBlockNumber()

	return blkHash, txHash, nil
}

func (cc *ChildChain) GetTx(txHash common.Hash) (*types.Tx, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistBlockTx(txHash) {
		return nil, ErrTxNotFound
	}

	return cc.getTx(txHash), nil
}

func (cc *ChildChain) GetTxIndex(txHash common.Hash) (*big.Int, *big.Int, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if !cc.isExistBlockTx(txHash) {
		return nil, nil, ErrTxNotFound
	}

	btx := cc.getBlockTx(txHash)

	return btx.BlockNumber, btx.TxIndex, nil
}

func (cc *ChildChain) GetTxProof(txHash common.Hash) ([]byte, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	// check tx existence
	if !cc.isExistBlockTx(txHash) {
		return nil, ErrTxNotFound
	}

	btx := cc.getBlockTx(txHash)

	blk, err := cc.getBlockByIndex(btx.BlockNumber)
	if err != nil {
		return nil, err
	}

	// build tx Merkle tree
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	// create tx proof
	return tree.CreateMembershipProof(btx.TxIndex.Uint64())
}

func (cc *ChildChain) AddTxToMempool(tx *types.Tx) (common.Hash, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if err := cc.validateTx(tx); err != nil {
		return types.NullHash, err
	}

	return cc.addTxToMempool(tx)
}

func (cc *ChildChain) ConfirmTx(txHash common.Hash, iIndex *big.Int, confSig types.Signature) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// check tx existence
	if !cc.isExistBlockTx(txHash) {
		return ErrTxNotFound
	}

	btx := cc.getBlockTx(txHash)

	// check txin existence
	if !btx.IsExistInput(iIndex) {
		return ErrTxInNotFound
	}

	txIn := btx.GetInput(iIndex)

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	inTxOut := cc.getTxOutByIndex(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

	// verify confirmation signature
	h, err := btx.ConfirmationHash()
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
	if err := cc.setConfirmationSignature(btx.BlockNumber, btx.TxIndex, iIndex, confSig); err != nil {
		return err
	}

	return nil
}

func (cc *ChildChain) currentBlockNumber() *big.Int {
	return cc.currentBlock.Number
}

func (cc *ChildChain) newCurrentBlockNumber() *big.Int {
	return new(big.Int).Set(cc.currentBlockNumber())
}

func (cc *ChildChain) newNextBlockNumber() *big.Int {
	return new(big.Int).Add(cc.currentBlockNumber(), big.NewInt(1))
}

func (cc *ChildChain) incrementBlockNumber() {
	cc.currentBlockNumber().Add(cc.currentBlockNumber(), big.NewInt(1))
}

func (cc *ChildChain) getBlockHash(blkNum *big.Int) common.Hash {
	return cc.chain[blkNum.String()]
}

func (cc *ChildChain) isExistBlockHash(blkNum *big.Int) bool {
	_, ok := cc.chain[blkNum.String()]
	return ok
}

func (cc *ChildChain) getBlock(blkHash common.Hash) (*types.Block, error) {
	lblk := cc.getLightBlock(blkHash)

	// get txes in block
	txes := make([]*types.Tx, len(lblk.TxHashes))
	for i, txHash := range lblk.TxHashes {
		txes[i] = cc.getTx(txHash)
	}

	// build block
	blk, err := types.NewBlock(txes, lblk.Number)
	if err != nil {
		return nil, err
	}
	blk.Signature = lblk.Signature

	return blk, nil
}

func (cc *ChildChain) getBlockByIndex(blkNum *big.Int) (*types.Block, error) {
	return cc.getBlock(cc.getBlockHash(blkNum))
}

func (cc *ChildChain) addBlock(blk *types.Block) (common.Hash, error) {
	lblk, err := blk.Lighten()
	if err != nil {
		return types.NullHash, err
	}

	blkHash, err := blk.Hash()
	if err != nil {
		return types.NullHash, err
	}
	blkHashStr := utils.HashToHex(blkHash)

	// update chain
	cc.chain[blk.Number.String()] = blkHash

	// store block
	cc.lightBlocks[blkHashStr] = lblk

	// store txes
	for i, tx := range blk.Txes {
		iBig := big.NewInt(int64(i))
		txHash := lblk.GetTxHash(iBig)
		cc.blockTxes[utils.HashToHex(txHash)] = tx.InBlock(blk.Number, iBig)
	}

	return blkHash, nil
}

func (cc *ChildChain) getLightBlock(blkHash common.Hash) *types.LightBlock {
	return cc.lightBlocks[utils.HashToHex(blkHash)]
}

func (cc *ChildChain) getLightBlockByIndex(blkNum *big.Int) *types.LightBlock {
	return cc.getLightBlock(cc.getBlockHash(blkNum))
}

func (cc *ChildChain) isExistLightBlock(blkHash common.Hash) bool {
	_, ok := cc.lightBlocks[utils.HashToHex(blkHash)]
	return ok
}

func (cc *ChildChain) isExistLightBlockByIndex(blkNum *big.Int) bool {
	if !cc.isExistBlockHash(blkNum) {
		return false
	}

	return cc.isExistLightBlock(cc.getBlockHash(blkNum))
}

func (cc *ChildChain) getTxHash(blkNum, txIndex *big.Int) common.Hash {
	return cc.getLightBlockByIndex(blkNum).GetTxHash(txIndex)
}

func (cc *ChildChain) getTx(txHash common.Hash) *types.Tx {
	return cc.getBlockTx(txHash).Tx
}

func (cc *ChildChain) getTxByIndex(blkNum, txIndex *big.Int) *types.Tx {
	return cc.getBlockTxByIndex(blkNum, txIndex).Tx
}

func (cc *ChildChain) addTxToMempool(tx *types.Tx) (common.Hash, error) {
	txHash, err := tx.Hash()
	if err != nil {
		return types.NullHash, err
	}

	for _, txIn := range tx.Inputs {
		if txIn.IsNull() {
			continue
		}

		// spend utxo
		if err := cc.spendUTXO(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex); err != nil {
			return types.NullHash, err
		}
	}

	// add tx to current block
	if err := cc.currentBlock.AddTx(tx); err != nil {
		return types.NullHash, err
	}

	return txHash, nil
}

func (cc *ChildChain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := big.NewInt(0), big.NewInt(0)

	for _, txOut := range tx.Outputs {
		oAmount.Add(oAmount, txOut.Amount)
	}

	for i, txIn := range tx.Inputs {
		// check spending txout existence
		if !cc.isExistTxOutByIndex(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex) {
			if txIn.IsNull() {
				nullTxInNum++
				continue
			}
			return ErrInvalidTxIn
		}

		inTxOut := cc.getTxOutByIndex(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

		// check double spent
		if inTxOut.IsSpent {
			return ErrTxOutAlreadySpent
		}

		// verify signature
		signerAddr, err := tx.SignerAddress(big.NewInt(int64(i)))
		if err != nil {
			return ErrInvalidTxSignature
		}
		if txIn.Signature == types.NullSignature ||
			!bytes.Equal(signerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
			return ErrInvalidTxSignature
		}

		iAmount.Add(iAmount, inTxOut.Amount)
	}

	// check txins validity
	if nullTxInNum == len(tx.Inputs) {
		return ErrInvalidTxIn
	}

	// check in/out balance
	if iAmount.Cmp(oAmount) < 0 {
		return ErrInvalidTxBalance
	}

	return nil
}

func (cc *ChildChain) getBlockTx(txHash common.Hash) *types.BlockTx {
	return cc.blockTxes[utils.HashToHex(txHash)]
}

func (cc *ChildChain) getBlockTxByIndex(blkNum, txIndex *big.Int) *types.BlockTx {
	return cc.getBlockTx(cc.getTxHash(blkNum, txIndex))
}

func (cc *ChildChain) isExistBlockTx(txHash common.Hash) bool {
	_, ok := cc.blockTxes[utils.HashToHex(txHash)]
	return ok
}

func (cc *ChildChain) isExistBlockTxByIndex(blkNum, txIndex *big.Int) bool {
	if !cc.isExistLightBlockByIndex(blkNum) {
		return false
	}

	lblk := cc.getLightBlockByIndex(blkNum)

	if !lblk.IsExistTxHash(txIndex) {
		return false
	}

	return cc.isExistBlockTx(lblk.GetTxHash(txIndex))
}

func (cc *ChildChain) getTxOutByIndex(blkNum, txIndex, oIndex *big.Int) *types.TxOut {
	return cc.getBlockTxByIndex(blkNum, txIndex).GetOutput(oIndex)
}

func (cc *ChildChain) isExistTxOutByIndex(blkNum, txIndex, oIndex *big.Int) bool {
	if !cc.isExistBlockTxByIndex(blkNum, txIndex) {
		return false
	}

	return cc.getBlockTxByIndex(blkNum, txIndex).IsExistOutput(oIndex)
}

func (cc *ChildChain) spendUTXO(blkNum, txIndex, oIndex *big.Int) error {
	return cc.blockTxes[utils.HashToHex(cc.getTxHash(blkNum, txIndex))].SpendOutput(oIndex)
}

func (cc *ChildChain) setConfirmationSignature(blkNum, txIndex, iIndex *big.Int, confSig types.Signature) error {
	return cc.blockTxes[utils.HashToHex(cc.getTxHash(blkNum, txIndex))].SetConfirmationSignature(iIndex, confSig)
}
