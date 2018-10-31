package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PostBlockResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockNumber uint64 `json:"blknum"`
	} `json:"result"`
}

func (client *Client) PostBlock(ctx context.Context) (uint64, error) {
	var resp PostBlockResponse
	if err := client.doAPI(
		ctx,
		http.MethodPost,
		"blocks",
		nil,
		&resp); err != nil {
		return 0, err
	}

	return resp.Result.BlockNumber, nil
}

type GetBlockResponse struct {
	State  string `json:"state"`
	Result struct {
		Hex string `json:"hex"`
	} `json:"result"`
}

func (client *Client) GetBlock(ctx context.Context, blkNum uint64) (*types.Block, error) {
	var resp GetBlockResponse
	if err := client.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("blocks/%d", blkNum),
		nil,
		&resp); err != nil {
		return nil, err
	}

	blkBytes, err := hexutil.Decode(resp.Result.Hex)
	if err != nil {
		return nil, err
	}

	var blk types.Block
	if err := rlp.DecodeBytes(blkBytes, &blk); err != nil {
		return nil, err
	}

	return &blk, nil
}
