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

func (tx *Tx) GetInput(iIndex *big.Int) *TxIn {
	if !tx.IsExistInput(iIndex) {
		return nil
	}

	return tx.Inputs[iIndex.Uint64()]
}

func (tx *Tx) SetInput(iIndex *big.Int, txIn *TxIn) error {
	if !tx.IsExistInput(iIndex) {
		return ErrInvalidTxInIndex
	}

	tx.Inputs[iIndex.Uint64()] = txIn

	return nil
}

func (tx *Tx) IsExistInput(iIndex *big.Int) bool {
	return iIndex.Cmp(TxElementsNumBig) < 0
}

func (tx *Tx) GetOutput(oIndex *big.Int) *TxOut {
	if !tx.IsExistOutput(oIndex) {
		return nil
	}

	return tx.Outputs[oIndex.Uint64()]
}

func (tx *Tx) SetOutput(oIndex *big.Int, txOut *TxOut) error {
	if !tx.IsExistOutput(oIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.Outputs[oIndex.Uint64()] = txOut

	return nil
}

func (tx *Tx) SpendOutput(oIndex *big.Int) error {
	if !tx.IsExistOutput(oIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.Outputs[oIndex.Uint64()].IsSpent = true

	return nil
}

func (tx *Tx) IsExistOutput(oIndex *big.Int) bool {
	return oIndex.Cmp(TxElementsNumBig) < 0
}

func (tx *Tx) Sign(iIndex *big.Int, signer *Account) error {
	if !tx.IsExistInput(iIndex) {
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

	tx.Inputs[iIndex.Uint64()].Signature = sig

	return nil
}

func (tx *Tx) Confirm(iIndex *big.Int, signer *Account) error {
	if !tx.IsExistInput(iIndex) {
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

	tx.Inputs[iIndex.Uint64()].ConfirmationSignature = sig

	return nil
}

func (tx *Tx) SignerAddress(iIndex *big.Int) (common.Address, error) {
	if !tx.IsExistInput(iIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.Hash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[iIndex.Uint64()].Signature)
}

func (tx *Tx) ConfirmationSignerAddress(iIndex *big.Int) (common.Address, error) {
	if !tx.IsExistInput(iIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.ConfirmationHash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.Inputs[iIndex.Uint64()].ConfirmationSignature)
}

func (tx *Tx) signerAddress(h common.Hash, sig Signature) (common.Address, error) {
	if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return sig.SignerAddress(h)
}

func (tx *Tx) SetConfirmationSignature(iIndex *big.Int, confSig Signature) error {
	if !tx.IsExistInput(iIndex) {
		return ErrInvalidTxInIndex
	}

	tx.Inputs[iIndex.Uint64()].ConfirmationSignature = confSig

	return nil
}
