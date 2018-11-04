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
	chain        map[string]string
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
		chain:        map[string]string{},
		lightBlocks:  map[string]*types.LightBlock{},
		blockTxes:    map[string]*types.BlockTx{},
	}, nil
}

func (cc *ChildChain) CurrentBlockNumber() *big.Int {
	return cc.currentBlock.Number
}

func (cc *ChildChain) GetBlockHash(blkNum *big.Int) ([]byte, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkHashStr, ok := cc.chain[blkNum.String()]
	if !ok {
		return nil, ErrBlockNotFound
	}

	return utils.DecodeHex(blkHashStr)
}

func (cc *ChildChain) GetBlock(blkHash common.Hash) (*types.Block, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	blkHashStr := utils.EncodeToHex(blkHash.Bytes())

	if _, ok := cc.lightBlocks[blkHashStr]; !ok {
		return nil, ErrBlockNotFound
	}

	return cc.getBlock(blkHashStr)
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
	blkNext, err := types.NewBlock(nil, blk.Number.Add(blk.Number, big.NewInt(1)))
	if err != nil {
		return types.NullHash, err
	}
	cc.currentBlock = blkNext

	return blkHash, nil
}

func (cc *ChildChain) AddDepositBlock(ownerAddr common.Address, amount *big.Int, signer *types.Account) (common.Hash, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	tx := types.NewTx()
	tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

	blk, err := types.NewBlock([]*types.Tx{tx}, cc.currentBlock.Number)
	if err != nil {
		return types.NullHash, err
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

	// increment current block number
	cc.currentBlock.Number.Add(cc.currentBlock.Number, big.NewInt(1))

	return blkHash, nil
}

func (cc *ChildChain) GetTx(txHash common.Hash) (*types.Tx, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHash.Bytes())

	btx, ok := cc.blockTxes[txHashStr]
	if !ok {
		return nil, ErrTxNotFound
	}

	return btx.Tx, nil
}

func (cc *ChildChain) GetTxProof(txHash common.Hash) ([]byte, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHash.Bytes())

	// check tx existence
	btx, ok := cc.blockTxes[txHashStr]
	if !ok {
		return nil, ErrTxNotFound
	}

	blk, err := cc.getBlock(cc.chain[btx.BlockNumber.String()])
	if err != nil {
		return nil, err
	}

	// build tx Merkle tree
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	// create proof
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

	txHashStr := utils.EncodeToHex(txHash.Bytes())

	// check tx existence
	btx, ok := cc.blockTxes[txHashStr]
	if !ok {
		return ErrTxNotFound
	}

	// check txin existence
	if iIndex.Cmp(types.TxElementsNumBig) >= 0 {
		return ErrTxInNotFound
	}

	txIn := btx.Inputs[iIndex.Uint64()]

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	inTxOut := cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

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
	cc.blockTxes[txHashStr].Inputs[iIndex.Uint64()].ConfirmationSignature = confSig

	return nil
}

func (cc *ChildChain) getBlock(blkHashStr string) (*types.Block, error) {
	lblk := cc.lightBlocks[blkHashStr]

	// get txes in block
	txes := make([]*types.Tx, len(lblk.TxHashes))
	for i, txHashStr := range lblk.TxHashes {
		txes[i] = cc.blockTxes[txHashStr].Tx
	}

	// build block
	blk, err := types.NewBlock(txes, lblk.Number)
	if err != nil {
		return nil, err
	}
	blk.Signature = lblk.Signature

	return blk, nil
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
	blkHashStr := utils.EncodeToHex(blkHash.Bytes())

	// update chain
	cc.chain[blk.Number.String()] = blkHashStr

	// store block
	cc.lightBlocks[blkHashStr] = lblk

	// store txes
	for i, tx := range blk.Txes {
		cc.blockTxes[lblk.TxHashes[i]] = tx.InBlock(blk.Number, big.NewInt(int64(i)))
	}

	return blkHash, nil
}

func (cc *ChildChain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := big.NewInt(0), big.NewInt(0)

	for _, txOut := range tx.Outputs {
		oAmount.Add(oAmount, txOut.Amount)
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
		cc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex).IsSpent = true
	}

	// add tx to current block
	cc.currentBlock.Txes = append(cc.currentBlock.Txes, tx)

	return txHash, nil
}

func (cc *ChildChain) isExistTxOut(blkNum, txIndex, oIndex *big.Int) bool {
	blkHashStr, ok := cc.chain[blkNum.String()]
	if !ok {
		return false
	}

	lblk, ok := cc.lightBlocks[blkHashStr]
	if !ok {
		return false
	}

	if txIndex.Cmp(big.NewInt(int64(len(lblk.TxHashes)))) >= 0 {
		return false
	}

	if oIndex.Cmp(types.TxElementsNumBig) >= 0 {
		return false
	}

	return true
}

func (cc *ChildChain) getTxOut(blkNum, txIndex, oIndex *big.Int) *types.TxOut {
	return cc.blockTxes[cc.lightBlocks[cc.chain[blkNum.String()]].TxHashes[txIndex.Uint64()]].Outputs[oIndex.Uint64()]
}
