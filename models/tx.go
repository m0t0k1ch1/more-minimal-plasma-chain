package models

import (
	"bytes"
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type TxIn struct {
	BlockNum    uint64 `json:"blknum"`
	TxIndex     uint64 `json:"txindex"`
	OutputIndex uint64 `json:"oindex"`
}

func NewTxIn(blkNum, txIndex, oIndex uint64) *TxIn {
	return &TxIn{
		BlockNum:    blkNum,
		TxIndex:     txIndex,
		OutputIndex: oIndex,
	}
}

func (txIn *TxIn) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txIn.BlockNum, txIn.TxIndex, txIn.OutputIndex,
	})
}

type TxOut struct {
	Owner  common.Address `json:"owner"`
	Amount uint64         `json:"amount"`
}

func NewTxOut(owner common.Address, amount uint64) *TxOut {
	return &TxOut{
		Owner:  owner,
		Amount: amount,
	}
}

func (txOut *TxOut) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		txOut.Owner.Bytes(), txOut.Amount,
	})
}

type TxSummary struct {
	Inputs                 [txElementsNum]*TxIn  `json:"ins"`
	Outputs                [txElementsNum]*TxOut `json:"outs"`
	Signatures             [txElementsNum]string `json:"sigs"`
	ConfirmationSignatures [txElementsNum]string `json:"confsigs"`
	Spents                 [txElementsNum]bool   `json:"spents"`
}

type Tx struct {
	Inputs                 [txElementsNum]*TxIn
	Outputs                [txElementsNum]*TxOut
	Signatures             [txElementsNum][]byte
	ConfirmationSignatures [txElementsNum][]byte
	Spents                 [txElementsNum]bool
}

func NewTx() *Tx {
	tx := &Tx{
		Inputs:                 [txElementsNum]*TxIn{},
		Outputs:                [txElementsNum]*TxOut{},
		Signatures:             [txElementsNum][]byte{},
		ConfirmationSignatures: [txElementsNum][]byte{},
		Spents:                 [txElementsNum]bool{},
	}

	for i := 0; i < txElementsNum; i++ {
		tx.Inputs[i] = nullTxIn
		tx.Outputs[i] = nullTxOut
		tx.Signatures[i] = nullSignature
		tx.ConfirmationSignatures[i] = nullSignature
		tx.Spents[i] = false
	}

	return tx
}

func (tx *Tx) IsDeposit() bool {
	for _, in := range tx.Inputs {
		if in.BlockNum != 0 {
			return false
		}
	}
	return true
}

// implements RLP Encoder interface
// ref. https://godoc.org/github.com/ethereum/go-ethereum/rlp#Encoder
func (tx *Tx) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		tx.Inputs, tx.Outputs,
	})
}

func (tx *Tx) Hash() ([]byte, error) {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(b), nil
}

func (tx *Tx) MerkleLeaf() ([]byte, error) {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	for _, sig := range tx.Signatures {
		if _, err := buf.Write(sig); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (tx *Tx) Summary() (*TxSummary, error) {
	summary := &TxSummary{
		Inputs:                 tx.Inputs,
		Outputs:                tx.Outputs,
		Signatures:             [txElementsNum]string{},
		ConfirmationSignatures: [txElementsNum]string{},
		Spents:                 tx.Spents,
	}

	for i, sig := range tx.Signatures {
		summary.Signatures[i] = common.Bytes2Hex(sig)
	}
	for i, confsig := range tx.ConfirmationSignatures {
		summary.ConfirmationSignatures[i] = common.Bytes2Hex(confsig)
	}

	return summary, nil
}

func (tx *Tx) Sign(idx int, privKey *ecdsa.PrivateKey) error {
	hashBytes, err := tx.Hash()
	if err != nil {
		return err
	}

	sig, err := crypto.Sign(hashBytes, privKey)
	if err != nil {
		return err
	}
	tx.Signatures[idx] = sig

	return nil
}

func (tx *Tx) Signers() ([]common.Address, error) {
	hashBytes, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	signers := make([]common.Address, txElementsNum)
	for i, sig := range tx.Signatures {
		if bytes.Equal(sig, nullSignature) {
			signers[i] = nullAddress
			continue
		}

		pubKey, err := crypto.SigToPub(hashBytes, sig)
		if err != nil {
			return nil, err
		}

		signers[i] = crypto.PubkeyToAddress(*pubKey)
	}

	return signers, nil
}
