package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type PostBlockResponse struct {
	*ResponseBase
	Result struct {
		BlockNumber *big.Int `json:"blknum"`
	} `json:"result"`
}

func (c *Client) PostBlock(ctx context.Context) (*big.Int, error) {
	var resp PostBlockResponse
	if err := c.doAPI(
		ctx,
		http.MethodPost,
		"blocks",
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return resp.Result.BlockNumber, nil
}

type GetBlockResponse struct {
	*ResponseBase
	Result struct {
		BlockStr string `json:"blk"`
	} `json:"result"`
}

func (c *Client) GetBlock(ctx context.Context, blkNum *big.Int) (*types.Block, error) {
	var resp GetBlockResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("blocks/%d", blkNum),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	blkBytes, err := utils.DecodeHex(resp.Result.BlockStr)
	if err != nil {
		return nil, err
	}

	var blk types.Block
	if err := rlp.DecodeBytes(blkBytes, &blk); err != nil {
		return nil, err
	}

	return &blk, nil
}
