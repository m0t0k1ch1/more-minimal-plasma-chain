package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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

func (client *Client) GetTx(ctx context.Context, txHashBytes []byte) (*types.Tx, error) {
	var resp GetTxResponse
	if err := client.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("txes/%s", utils.EncodeToHex(txHashBytes)),
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
