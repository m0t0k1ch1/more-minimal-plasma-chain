package types

import (
	"bytes"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	TxElementsNum = 2
)

type TxCore struct {
	Inputs     [TxElementsNum]*TxIn     `json:"ins"`
	Outputs    [TxElementsNum]*TxOut    `json:"outs"`
	Signatures [TxElementsNum]Signature `json:"sigs"`
}

type Tx struct {
	*TxCore
	ConfirmationSignatures [TxElementsNum]Signature `json:"confsigs"`
	Spents                 [TxElementsNum]bool      `json:"spents"`
}

func NewTx() *Tx {
	tx := &Tx{
		TxCore: &TxCore{
			Inputs:     [TxElementsNum]*TxIn{},
			Outputs:    [TxElementsNum]*TxOut{},
			Signatures: [TxElementsNum]Signature{},
		},
		ConfirmationSignatures: [TxElementsNum]Signature{},
		Spents:                 [TxElementsNum]bool{},
	}

	for i := 0; i < TxElementsNum; i++ {
		tx.Inputs[i] = nullTxIn
		tx.Outputs[i] = nullTxOut
		tx.Signatures[i] = NullSignature
		tx.ConfirmationSignatures[i] = NullSignature
		tx.Spents[i] = false
	}

	return tx
}

func (tx *Tx) IsDeposit() bool {
	for _, txIn := range tx.Inputs {
		if txIn.BlockNumber != 0 {
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

func (tx *Tx) ConfirmationHash() ([]byte, error) {
	hashBytes, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(hashBytes), nil
}

func (tx *Tx) MerkleLeaf() ([]byte, error) {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	for _, sig := range tx.Signatures {
		if _, err := buf.Write(sig.Bytes()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (tx *Tx) Sign(iIndex uint64, signer *Account) error {
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

	tx.Signatures[iIndex] = sig

	return nil
}

func (tx *Tx) Confirm(iIndex uint64, signer *Account) error {
	confHashBytes, err := tx.ConfirmationHash()
	if err != nil {
		return err
	}

	confSigBytes, err := signer.Sign(confHashBytes)
	if err != nil {
		return err
	}

	sig, err := NewSignatureFromBytes(confSigBytes)
	if err != nil {
		return err
	}

	tx.ConfirmationSignatures[iIndex] = sig

	return nil
}

func (tx *Tx) SignerAddresses() ([]common.Address, error) {
	hashBytes, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	return tx.signerAddresses(hashBytes, tx.Signatures)
}

func (tx *Tx) ConfirmationSignerAddresses() ([]common.Address, error) {
	confHashBytes, err := tx.ConfirmationHash()
	if err != nil {
		return nil, err
	}

	return tx.signerAddresses(confHashBytes, tx.ConfirmationSignatures)
}

func (tx *Tx) signerAddresses(b []byte, sigs [TxElementsNum]Signature) ([]common.Address, error) {
	signerAddrs := make([]common.Address, len(sigs))
	for i, sig := range sigs {
		if bytes.Equal(sig.Bytes(), NullSignature.Bytes()) {
			signerAddrs[i] = NullAddress
			continue
		}

		signerAddr, err := sig.SignerAddress(b)
		if err != nil {
			return nil, err
		}

		signerAddrs[i] = signerAddr
	}

	return signerAddrs, nil
}
