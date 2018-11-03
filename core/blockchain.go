package core

import (
	"bytes"
	"errors"
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

type Blockchain struct {
	mu           *sync.RWMutex
	currentBlock *types.Block
	chain        map[uint64]string
	lightBlocks  map[string]*types.LightBlock
	blockTxes    map[string]*types.BlockTx
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		mu:           &sync.RWMutex{},
		currentBlock: types.NewBlock(nil, DefaultBlockNumber),
		chain:        map[uint64]string{},
		lightBlocks:  map[string]*types.LightBlock{},
		blockTxes:    map[string]*types.BlockTx{},
	}
}

func (bc *Blockchain) CurrentBlockNumber() uint64 {
	return bc.currentBlock.Number
}

func (bc *Blockchain) GetBlockHash(blkNum uint64) ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	blkHashStr, ok := bc.chain[blkNum]
	if !ok {
		return nil, ErrBlockNotFound
	}

	return utils.DecodeHex(blkHashStr)
}

func (bc *Blockchain) GetBlock(blkHashBytes []byte) (*types.Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	blkHashStr := utils.EncodeToHex(blkHashBytes)

	if _, ok := bc.lightBlocks[blkHashStr]; !ok {
		return nil, ErrBlockNotFound
	}

	return bc.getBlock(blkHashStr), nil
}

func (bc *Blockchain) AddBlock(signer *types.Account) ([]byte, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	blk := bc.currentBlock

	// check block validity
	if len(blk.Txes) == 0 {
		return nil, ErrEmptyBlock
	}

	// sign block
	if err := blk.Sign(signer); err != nil {
		return nil, err
	}

	// add block
	blkHashBytes, err := bc.addBlock(blk)
	if err != nil {
		return nil, err
	}

	// reset current block
	bc.currentBlock = types.NewBlock(nil, blk.Number+1)

	return blkHashBytes, nil
}

func (bc *Blockchain) AddDepositBlock(ownerAddr common.Address, amount uint64, signer *types.Account) ([]byte, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	tx := types.NewTx()
	tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

	blk := types.NewBlock([]*types.Tx{tx}, bc.currentBlock.Number)

	// sign block
	if err := blk.Sign(signer); err != nil {
		return nil, err
	}

	// add block
	blkHashBytes, err := bc.addBlock(blk)
	if err != nil {
		return nil, err
	}

	// increment current block number
	bc.currentBlock.Number++

	return blkHashBytes, nil
}

func (bc *Blockchain) GetTx(txHashBytes []byte) (*types.Tx, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHashBytes)

	btx, ok := bc.blockTxes[txHashStr]
	if !ok {
		return nil, ErrTxNotFound
	}

	return btx.Tx, nil
}

func (bc *Blockchain) GetTxProof(txHashBytes []byte) ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHashBytes)

	// check tx existence
	btx, ok := bc.blockTxes[txHashStr]
	if !ok {
		return nil, ErrTxNotFound
	}

	blk := bc.getBlock(bc.chain[btx.BlockNumber])

	// build tx Merkle tree
	tree, err := blk.MerkleTree()
	if err != nil {
		return nil, err
	}

	// create proof
	return tree.CreateMembershipProof(btx.TxIndex)
}

func (bc *Blockchain) AddTxToMempool(tx *types.Tx) ([]byte, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := bc.validateTx(tx); err != nil {
		return nil, err
	}

	return bc.addTxToMempool(tx)
}

func (bc *Blockchain) ConfirmTx(txHashBytes []byte, iIndex uint64, confSig types.Signature) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	txHashStr := utils.EncodeToHex(txHashBytes)

	// check tx existence
	btx, ok := bc.blockTxes[txHashStr]
	if !ok {
		return ErrTxNotFound
	}

	// check txin existence
	if iIndex >= uint64(len(btx.Inputs)) {
		return ErrTxInNotFound
	}

	txIn := btx.Inputs[iIndex]

	// check txin validity
	if txIn.IsNull() {
		return ErrNullTxInConfirmation
	}

	inTxOut := bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

	// verify confirmation signature
	confHashBytes, err := btx.ConfirmationHash()
	if err != nil {
		return err
	}
	confSignerAddr, err := confSig.SignerAddress(confHashBytes)
	if err != nil {
		return ErrInvalidTxConfirmationSignature
	}
	if !bytes.Equal(confSignerAddr.Bytes(), inTxOut.OwnerAddress.Bytes()) {
		return ErrInvalidTxConfirmationSignature
	}

	// update confirmation signature
	bc.blockTxes[txHashStr].Inputs[iIndex].ConfirmationSignature = confSig

	return nil
}

func (bc *Blockchain) getBlock(blkHashStr string) *types.Block {
	lblk := bc.lightBlocks[blkHashStr]

	// get txes in block
	txes := make([]*types.Tx, len(lblk.TxHashes))
	for i, txHashStr := range lblk.TxHashes {
		txes[i] = bc.blockTxes[txHashStr].Tx
	}

	// build block
	blk := types.NewBlock(txes, lblk.Number)
	blk.Signature = lblk.Signature

	return blk
}

func (bc *Blockchain) addBlock(blk *types.Block) ([]byte, error) {
	lblk, err := blk.Lighten()
	if err != nil {
		return nil, err
	}

	blkHashBytes, err := blk.Hash()
	if err != nil {
		return nil, err
	}
	blkHashStr := utils.EncodeToHex(blkHashBytes)

	// update chain
	bc.chain[blk.Number] = blkHashStr

	// store block
	bc.lightBlocks[blkHashStr] = lblk

	// store txes
	for i, tx := range blk.Txes {
		bc.blockTxes[lblk.TxHashes[i]] = tx.InBlock(blk.Number, uint64(i))
	}

	return blkHashBytes, nil
}

func (bc *Blockchain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := uint64(0), uint64(0)

	for _, txOut := range tx.Outputs {
		oAmount += txOut.Amount
	}

	for i, txIn := range tx.Inputs {
		// check spending txout existence
		if !bc.isExistTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex) {
			if txIn.IsNull() {
				nullTxInNum++
				continue
			}
			return ErrInvalidTxIn
		}

		inTxOut := bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

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

func (bc *Blockchain) addTxToMempool(tx *types.Tx) ([]byte, error) {
	txHashBytes, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	for _, txIn := range tx.Inputs {
		if txIn.IsNull() {
			continue
		}

		// spend utxo
		bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex).IsSpent = true
	}

	// add tx to current block
	bc.currentBlock.Txes = append(bc.currentBlock.Txes, tx)

	return txHashBytes, nil
}

func (bc *Blockchain) isExistTxOut(blkNum, txIndex, oIndex uint64) bool {
	blkHashStr, ok := bc.chain[blkNum]
	if !ok {
		return false
	}

	lblk, ok := bc.lightBlocks[blkHashStr]
	if !ok {
		return false
	}

	if txIndex >= uint64(len(lblk.TxHashes)) {
		return false
	}
	txHashStr := lblk.TxHashes[txIndex]

	btx, ok := bc.blockTxes[txHashStr]
	if !ok {
		return false
	}

	if oIndex >= uint64(len(btx.Outputs)) {
		return false
	}

	return true
}

func (bc *Blockchain) getTxOut(blkNum, txIndex, oIndex uint64) *types.TxOut {
	return bc.blockTxes[bc.lightBlocks[bc.chain[blkNum]].TxHashes[txIndex]].Outputs[oIndex]
}
