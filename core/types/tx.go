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

func (tx *Tx) IsDeposit() bool {
	for _, txIn := range tx.Inputs {
		if txIn.BlockNumber != 0 {
			return false
		}
	}
	return true
}

func (tx *Tx) Hash() ([]byte, error) {
	b, err := rlp.EncodeToBytes([]interface{}{
		tx.inputCores(), tx.outputCores(),
	})
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(b), nil
}

func (tx *Tx) ConfirmationHash() ([]byte, error) {
	hashBytes, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(hashBytes), nil
}

func (tx *Tx) MerkleLeaf() ([]byte, error) {
	b, err := rlp.EncodeToBytes([]interface{}{
		tx.inputCores(), tx.outputCores(),
	})
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

	hashBytes, err := tx.Hash()
	if err != nil {
		return err
	}

	sigBytes, err := signer.Sign(hashBytes)
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

	confHashBytes, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}

	confSigBytes, err := signer.Sign(confHashBytes)
	if err != nil {
		return err
	}
	confSig, err := NewSignatureFromBytes(confSigBytes)
	if err != nil {
		return err
	}

	tx.Inputs[iIndex].ConfirmationSignature = confSig

	return nil
}

func (tx *Tx) SignerAddress(iIndex uint64) (common.Address, error) {
	if iIndex >= uint64(len(tx.Inputs)) {
		return NullAddress, ErrInvalidTxInIndex
	}

	hashBytes, err := tx.Hash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(hashBytes, tx.Inputs[iIndex].Signature)
}

func (tx *Tx) ConfirmationSignerAddress(iIndex uint64) (common.Address, error) {
	if iIndex >= uint64(len(tx.Inputs)) {
		return NullAddress, ErrInvalidTxInIndex
	}

	confHashBytes, err := tx.ConfirmationHash()
	if err != nil {
		return NullAddress, err
	}

	return tx.signerAddress(confHashBytes, tx.Inputs[iIndex].ConfirmationSignature)
}

func (tx *Tx) signerAddress(b []byte, sig Signature) (common.Address, error) {
	if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
		return NullAddress, nil
	}

	return sig.SignerAddress(b)
}
