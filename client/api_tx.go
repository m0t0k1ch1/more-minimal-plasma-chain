package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type PostTxResponse struct {
	State  string `json:"state"`
	Result struct {
		TxHashStr string `json:"txhash"`
	} `json:"result"`
}

func (client *Client) PostTx(ctx context.Context, tx *types.Tx) ([]byte, error) {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("tx", utils.EncodeToHex(txBytes))

	var resp PostTxResponse
	if err := client.doAPI(
		ctx,
		http.MethodPost,
		"txes",
		v,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.TxHashStr)
}

type GetTxResponse struct {
	State  string `json:"state"`
	Result struct {
		TxStr string `json:"tx"`
	} `json:"result"`
}

func (client *Client) GetTx(ctx context.Context, txHash common.Hash) (*types.Tx, error) {
	var resp GetTxResponse
	if err := client.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("txes/%s", utils.EncodeToHex(txHash.Bytes())),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	txBytes, err := utils.DecodeHex(resp.Result.TxStr)
	if err != nil {
		return nil, err
	}

	var tx types.Tx
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

type GetTxProofResponse struct {
	State  string `json:"state"`
	Result struct {
		ProofStr string `json:"proof"`
	} `json:"result"`
}

func (client *Client) GetTxProof(ctx context.Context, txHash common.Hash) ([]byte, error) {
	var resp GetTxProofResponse
	if err := client.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("txes/%s/proof", utils.EncodeToHex(txHash.Bytes())),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.ProofStr)
}

type PutTxResponse struct {
	State  string `json:"state"`
	Result struct {
		TxHashStr string `json:"txhash"`
	} `json:"result"`
}

func (client *Client) PutTx(ctx context.Context, txHash common.Hash, iIndex *big.Int, confSig types.Signature) ([]byte, error) {
	v := url.Values{}
	v.Set("index", iIndex.String())
	v.Set("confsig", utils.EncodeToHex(confSig.Bytes()))

	var resp PutTxResponse
	if err := client.doAPI(
		ctx,
		http.MethodPut,
		fmt.Sprintf("txes/%s", utils.EncodeToHex(txHash.Bytes())),
		v,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.TxHashStr)
}
