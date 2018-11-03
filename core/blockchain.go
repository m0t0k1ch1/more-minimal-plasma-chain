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

type Blockchain struct {
	mu           *sync.RWMutex
	currentBlock *types.Block
	chain        map[*big.Int]string
	lightBlocks  map[string]*types.LightBlock
	blockTxes    map[string]*types.BlockTx
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		mu:           &sync.RWMutex{},
		currentBlock: types.NewBlock(nil, big.NewInt(DefaultBlockNumber)),
		chain:        map[*big.Int]string{},
		lightBlocks:  map[string]*types.LightBlock{},
		blockTxes:    map[string]*types.BlockTx{},
	}
}

func (bc *Blockchain) CurrentBlockNumber() *big.Int {
	return bc.currentBlock.Number
}

func (bc *Blockchain) GetBlockHash(blkNum *big.Int) ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	blkHashStr, ok := bc.chain[blkNum]
	if !ok {
		return nil, ErrBlockNotFound
	}

	return utils.DecodeHex(blkHashStr)
}

func (bc *Blockchain) GetBlock(blkHash common.Hash) (*types.Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	blkHashStr := utils.EncodeToHex(blkHash.Bytes())

	if _, ok := bc.lightBlocks[blkHashStr]; !ok {
		return nil, ErrBlockNotFound
	}

	return bc.getBlock(blkHashStr), nil
}

func (bc *Blockchain) AddBlock(signer *types.Account) (common.Hash, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	blk := bc.currentBlock

	// check block validity
	if len(blk.Txes) == 0 {
		return types.NullHash, ErrEmptyBlock
	}

	// sign block
	if err := blk.Sign(signer); err != nil {
		return types.NullHash, err
	}

	// add block
	blkHash, err := bc.addBlock(blk)
	if err != nil {
		return types.NullHash, err
	}

	// reset current block
	bc.currentBlock = types.NewBlock(nil, blk.Number.Add(blk.Number, big.NewInt(1)))

	return blkHash, nil
}

func (bc *Blockchain) AddDepositBlock(ownerAddr common.Address, amount *big.Int, signer *types.Account) (common.Hash, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	tx := types.NewTx()
	tx.Outputs[0] = types.NewTxOut(ownerAddr, amount)

	blk := types.NewBlock([]*types.Tx{tx}, bc.currentBlock.Number)

	// sign block
	if err := blk.Sign(signer); err != nil {
		return types.NullHash, err
	}

	// add block
	blkHash, err := bc.addBlock(blk)
	if err != nil {
		return types.NullHash, err
	}

	// increment current block number
	bc.currentBlock.Number.Add(bc.currentBlock.Number, big.NewInt(1))

	return blkHash, nil
}

func (bc *Blockchain) GetTx(txHash common.Hash) (*types.Tx, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHash.Bytes())

	btx, ok := bc.blockTxes[txHashStr]
	if !ok {
		return nil, ErrTxNotFound
	}

	return btx.Tx, nil
}

func (bc *Blockchain) GetTxProof(txHash common.Hash) ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	txHashStr := utils.EncodeToHex(txHash.Bytes())

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
	return tree.CreateMembershipProof(btx.TxIndex.Uint64())
}

func (bc *Blockchain) AddTxToMempool(tx *types.Tx) (common.Hash, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := bc.validateTx(tx); err != nil {
		return types.NullHash, err
	}

	return bc.addTxToMempool(tx)
}

func (bc *Blockchain) ConfirmTx(txHash common.Hash, iIndex *big.Int, confSig types.Signature) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	txHashStr := utils.EncodeToHex(txHash.Bytes())

	// check tx existence
	btx, ok := bc.blockTxes[txHashStr]
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

	inTxOut := bc.getTxOut(txIn.BlockNumber, txIn.TxIndex, txIn.OutputIndex)

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
	bc.blockTxes[txHashStr].Inputs[iIndex.Uint64()].ConfirmationSignature = confSig

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

func (bc *Blockchain) addBlock(blk *types.Block) (common.Hash, error) {
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
	bc.chain[blk.Number] = blkHashStr

	// store block
	bc.lightBlocks[blkHashStr] = lblk

	// store txes
	for i, tx := range blk.Txes {
		bc.blockTxes[lblk.TxHashes[i]] = tx.InBlock(blk.Number, big.NewInt(int64(i)))
	}

	return blkHash, nil
}

func (bc *Blockchain) validateTx(tx *types.Tx) error {
	nullTxInNum := 0
	iAmount, oAmount := big.NewInt(0), big.NewInt(0)

	for _, txOut := range tx.Outputs {
		oAmount.Add(oAmount, txOut.Amount)
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

func (bc *Blockchain) addTxToMempool(tx *types.Tx) (common.Hash, error) {
	txHash, err := tx.Hash()
	if err != nil {
		return types.NullHash, err
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

	return txHash, nil
}

func (bc *Blockchain) isExistTxOut(blkNum, txIndex, oIndex *big.Int) bool {
	blkHashStr, ok := bc.chain[blkNum]
	if !ok {
		return false
	}

	lblk, ok := bc.lightBlocks[blkHashStr]
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

func (bc *Blockchain) getTxOut(blkNum, txIndex, oIndex *big.Int) *types.TxOut {
	return bc.blockTxes[bc.lightBlocks[bc.chain[blkNum]].TxHashes[txIndex.Uint64()]].Outputs[oIndex.Uint64()]
}
