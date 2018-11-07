package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type PostTxResponse struct {
	State  string `json:"state"`
	Result struct {
		PosBig *big.Int `json:"pos"`
	} `json:"result"`
}

func (c *Client) PostTx(ctx context.Context, tx *types.Tx) (types.Position, error) {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return types.NullPosition, err
	}

	v := url.Values{}
	v.Set("tx", utils.EncodeToHex(txBytes))

	var resp PostTxResponse
	if err := c.doAPI(
		ctx,
		http.MethodPost,
		"txes",
		v,
		&resp,
	); err != nil {
		return types.NullPosition, err
	}

	return types.NewPosition(resp.Result.PosBig), nil
}

type GetTxResponse struct {
	State  string `json:"state"`
	Result struct {
		TxStr string `json:"tx"`
	} `json:"result"`
}

func (c *Client) GetTx(ctx context.Context, txPos types.Position) (*types.Tx, error) {
	var resp GetTxResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("txes/%s", txPos.String()),
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

func (c *Client) GetTxProof(ctx context.Context, txPos types.Position) ([]byte, error) {
	var resp GetTxProofResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("txes/%s/proof", txPos.String()),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.ProofStr)
}
