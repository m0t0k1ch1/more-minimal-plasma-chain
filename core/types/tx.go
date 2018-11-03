package types

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	TxElementsNum = 2
)

var (
	ErrInvalidTxInIndex = errors.New("invalid txin index")
)

type BlockTx struct {
	*Tx
	BlockNumber uint64
	TxIndex     uint64
}

type Tx struct {
	Inputs  [TxElementsNum]*TxIn  `json:"ins"`
	Outputs [TxElementsNum]*TxOut `json:"outs"`
}

func NewTx() *Tx {
	tx := &Tx{
		Inputs:  [TxElementsNum]*TxIn{},
		Outputs: [TxElementsNum]*TxOut{},
	}

	for i := 0; i < TxElementsNum; i++ {
		tx.Inputs[i] = NewTxIn(0, 0, 0)
		tx.Outputs[i] = NewTxOut(NullAddress, 0)
	}

	return tx
}

func (tx *Tx) inputCores() [TxElementsNum]*TxInCore {
	txInCores := [TxElementsNum]*TxInCore{}
	for i := 0; i < TxElementsNum; i++ {
		txInCores[i] = tx.Inputs[i].TxInCore
	}
	return txInCores
}

func (tx *Tx) outputCores() [TxElementsNum]*TxOutCore {
	txOutCores := [TxElementsNum]*TxOutCore{}
	for i := 0; i < TxElementsNum; i++ {
		txOutCores[i] = tx.Outputs[i].TxOutCore
	}
	return txOutCores
}

func (tx *Tx) signatures() [TxElementsNum]Signature {
	sigs := [TxElementsNum]Signature{}
	for i := 0; i < TxElementsNum; i++ {
		sigs[i] = tx.Inputs[i].Signature
	}
	return sigs
}

func (tx *Tx) Encode() ([]byte, error) {
	return rlp.EncodeToBytes([]interface{}{
		tx.inputCores(), tx.outputCores(),
	})
}

func (tx *Tx) Hash() (common.Hash, error) {
	b, err := tx.Encode()
	if err != nil {
		return NullHash, err
	}

	return common.BytesToHash(crypto.Keccak256(b)), nil
}

func (tx *Tx) ConfirmationHash() (common.Hash, error) {
	h, err := tx.Hash()
	if err != nil {
		return NullHash, err
	}

	return common.BytesToHash(crypto.Keccak256(h.Bytes())), nil
}

func (tx *Tx) MerkleLeaf() ([]byte, error) {
	b, err := tx.Encode()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	for _, txIn := range tx.Inputs {
		if _, err := buf.Write(txIn.Signature.Bytes()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (tx *Tx) Sign(iIndex uint64, signer *Account) error {
	if iIndex >= uint64(len(tx.Inputs)) {
		return ErrInvalidTxInIndex
	}

	h, err := tx.Hash()
	if err != nil {
		return err
	}

	sigBytes, err := signer.Sign(h)
	if err != nil {
		return err
	}
	sig, err := NewSignatureFromBytes(sigBytes)
	if err != nil {
		return err
	}

	tx.Inputs[iIndex].Signature = sig

	return nil
}

func (tx *Tx) Confirm(iIndex uint64, signer *Account) error {
	if iIndex >= uint64(len(tx.Inputs)) {
		return ErrInvalidTxInIndex
	}

	h, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}

	sigBytes, err := signer.Sign(h)
	if err != nil {
		return err
	}
	sig, err := NewSignatureFromBytes(sigBytes)
	if err != nil {
		return err
	}

	tx.Inputs[iIndex].ConfirmationSignature = sig

	return nil
}

func (tx *Tx) SignerAddress(iIndex uint64) (common.Address, error) {
	if iIndex >= uint64(len(tx.Inputs)) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.Hash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[iIndex].Signature)
}

func (tx *Tx) ConfirmationSignerAddress(iIndex uint64) (common.Address, error) {
	if iIndex >= uint64(len(tx.Inputs)) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.ConfirmationHash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[iIndex].ConfirmationSignature)
}

func (tx *Tx) signerAddress(h common.Hash, sig Signature) (common.Address, error) {
	if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return sig.SignerAddress(h)
}

func (tx *Tx) InBlock(blkNum, txIndex uint64) *BlockTx {
	return &BlockTx{
		Tx:          tx,
		BlockNumber: blkNum,
		TxIndex:     txIndex,
	}
}
