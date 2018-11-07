package types

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	TxElementsNum = 2
)

var (
	TxElementsNumBig = big.NewInt(TxElementsNum)
)

var (
	ErrInvalidTxInIndex  = errors.New("tx input index is invalid")
	ErrInvalidTxOutIndex = errors.New("tx ouptut index is invalid")
)

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
		tx.Inputs[i] = NewTxIn(big.NewInt(0), big.NewInt(0), big.NewInt(0))
		tx.Outputs[i] = NewTxOut(NullAddress, big.NewInt(0))
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

	return utils.BytesToHash(crypto.Keccak256(b)), nil
}

func (tx *Tx) ConfirmationHash() (common.Hash, error) {
	h, err := tx.Hash()
	if err != nil {
		return NullHash, err
	}

	return utils.BytesToHash(crypto.Keccak256(h.Bytes())), nil
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

func (tx *Tx) SignaturesBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, txIn := range tx.Inputs {
		if _, err := buf.Write(txIn.Signature.Bytes()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (tx *Tx) ConfirmationSignaturesBytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, txIn := range tx.Inputs {
		if _, err := buf.Write(txIn.ConfirmationSignature.Bytes()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (tx *Tx) GetInput(inIndex *big.Int) *TxIn {
	if !tx.IsExistInput(inIndex) {
		return nil
	}

	return tx.Inputs[inIndex.Uint64()]
}

func (tx *Tx) SetInput(inIndex *big.Int, txIn *TxIn) error {
	if !tx.IsExistInput(inIndex) {
		return ErrInvalidTxInIndex
	}

	tx.Inputs[inIndex.Uint64()] = txIn

	return nil
}

func (tx *Tx) IsExistInput(inIndex *big.Int) bool {
	return inIndex.Cmp(TxElementsNumBig) < 0
}

func (tx *Tx) GetOutput(outIndex *big.Int) *TxOut {
	if !tx.IsExistOutput(outIndex) {
		return nil
	}

	return tx.Outputs[outIndex.Uint64()]
}

func (tx *Tx) SetOutput(outIndex *big.Int, txOut *TxOut) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.Outputs[outIndex.Uint64()] = txOut

	return nil
}

func (tx *Tx) SpendOutput(outIndex *big.Int) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.Outputs[outIndex.Uint64()].IsSpent = true

	return nil
}

func (tx *Tx) IsExistOutput(outIndex *big.Int) bool {
	return outIndex.Cmp(TxElementsNumBig) < 0
}

func (tx *Tx) Sign(inIndex *big.Int, signer *Account) error {
	if !tx.IsExistInput(inIndex) {
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
	sig, err := BytesToSignature(sigBytes)
	if err != nil {
		return err
	}

	tx.Inputs[inIndex.Uint64()].Signature = sig

	return nil
}

func (tx *Tx) Confirm(inIndex *big.Int, signer *Account) error {
	if !tx.IsExistInput(inIndex) {
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
	sig, err := BytesToSignature(sigBytes)
	if err != nil {
		return err
	}

	tx.Inputs[inIndex.Uint64()].ConfirmationSignature = sig

	return nil
}

func (tx *Tx) SignerAddress(inIndex *big.Int) (common.Address, error) {
	if !tx.IsExistInput(inIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.Hash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[inIndex.Uint64()].Signature)
}

func (tx *Tx) ConfirmationSignerAddress(inIndex *big.Int) (common.Address, error) {
	if !tx.IsExistInput(inIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.ConfirmationHash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[inIndex.Uint64()].ConfirmationSignature)
}

func (tx *Tx) signerAddress(h common.Hash, sig Signature) (common.Address, error) {
	if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return sig.SignerAddress(h)
}

func (tx *Tx) SetConfirmationSignature(inIndex *big.Int, confSig Signature) error {
	if !tx.IsExistInput(inIndex) {
		return ErrInvalidTxInIndex
	}

	tx.Inputs[inIndex.Uint64()].ConfirmationSignature = confSig

	return nil
}
