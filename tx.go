package main

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	TxElementsNum = 2
)

var (
	NullAddress = common.BytesToAddress([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
	NullSignature = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	NullTxIn  = NewTxIn(0, 0, 0)
	NullTxOut = NewTxOut(NullAddress, 0)
	NullTx    = NewTx()
)

type TxIn struct {
	BlockNum    uint
	TxIndex     uint
	OutputIndex uint
}

func NewTxIn(blkNum, txIndex, oIndex uint) *TxIn {
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
	Owner  common.Address
	Amount uint
}

func NewTxOut(owner common.Address, amount uint) *TxOut {
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

type Tx struct {
	Inputs        [TxElementsNum]*TxIn
	Outputs       [TxElementsNum]*TxOut
	Signatures    [TxElementsNum][]byte
	Confirmations [TxElementsNum][]byte
}

func NewTx() *Tx {
	tx := &Tx{
		Inputs:        [TxElementsNum]*TxIn{},
		Outputs:       [TxElementsNum]*TxOut{},
		Signatures:    [TxElementsNum][]byte{},
		Confirmations: [TxElementsNum][]byte{},
	}

	for i := 0; i < TxElementsNum; i++ {
		tx.Inputs[i] = NullTxIn
		tx.Outputs[i] = NullTxOut
		tx.Signatures[i] = NullSignature
		tx.Confirmations[i] = NullSignature
	}

	return tx
}

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
