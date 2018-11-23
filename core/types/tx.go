package types

import (
	"bytes"
	"fmt"
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
	ErrInvalidTxInIndex  = fmt.Errorf("txin index must be less than %d", TxElementsNum)
	ErrInvalidTxOutIndex = fmt.Errorf("txout index must be less than %d", TxElementsNum)
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
		tx.Inputs[i] = NewTxIn(0, 0, 0)
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

func (tx *Tx) GetInput(inIndex uint64) *TxIn {
	if !tx.IsExistInput(inIndex) {
		return nil
	}

	return tx.Inputs[inIndex]
}

func (tx *Tx) SetInput(inIndex uint64, txIn *TxIn) error {
	if !tx.IsExistInput(inIndex) {
		return ErrInvalidTxInIndex
	}

	tx.Inputs[inIndex] = txIn

	return nil
}

func (tx *Tx) IsExistInput(inIndex uint64) bool {
	return inIndex < TxElementsNum
}

func (tx *Tx) GetOutput(outIndex uint64) *TxOut {
	if !tx.IsExistOutput(outIndex) {
		return nil
	}

	return tx.Outputs[outIndex]
}

func (tx *Tx) SetOutput(outIndex uint64, txOut *TxOut) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.Outputs[outIndex] = txOut

	return nil
}

func (tx *Tx) SpendOutput(outIndex uint64) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.GetOutput(outIndex).IsSpent = true

	return nil
}

func (tx *Tx) ExitOutput(outIndex uint64) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.GetOutput(outIndex).IsExited = true

	return nil
}

func (tx *Tx) RestoreOutput(outIndex uint64) error {
	if !tx.IsExistOutput(outIndex) {
		return ErrInvalidTxOutIndex
	}

	tx.GetOutput(outIndex).IsExited = false

	return nil
}

func (tx *Tx) IsExistOutput(outIndex uint64) bool {
	return outIndex < TxElementsNum
}

func (tx *Tx) Sign(inIndex uint64, signer *Account) error {
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

	tx.GetInput(inIndex).Signature = sig

	return nil
}

func (tx *Tx) Confirm(inIndex uint64, signer *Account) error {
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

	tx.GetInput(inIndex).ConfirmationSignature = sig

	return nil
}

func (tx *Tx) SignerAddress(inIndex uint64) (common.Address, error) {
	if !tx.IsExistInput(inIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.Hash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.GetInput(inIndex).Signature)
}

func (tx *Tx) ConfirmationSignerAddress(inIndex uint64) (common.Address, error) {
	if !tx.IsExistInput(inIndex) {
		return NullAddress, ErrInvalidTxInIndex
	}

	h, err := tx.ConfirmationHash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(h, tx.GetInput(inIndex).ConfirmationSignature)
}

func (tx *Tx) signerAddress(h common.Hash, sig Signature) (common.Address, error) {
	if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return sig.SignerAddress(h)
}

func (tx *Tx) SetConfirmationSignature(inIndex uint64, confSig Signature) error {
	if !tx.IsExistInput(inIndex) {
		return ErrInvalidTxInIndex
	}

	tx.GetInput(inIndex).ConfirmationSignature = confSig

	return nil
}
