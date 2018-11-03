package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type PostBlockResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockHashStr string `json:"blkhash"`
	} `json:"result"`
}

func (c *Client) PostBlock(ctx context.Context) ([]byte, error) {
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

	return utils.DecodeHex(resp.Result.BlockHashStr)
}

type GetBlockResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockStr string `json:"blk"`
	} `json:"result"`
}

func (c *Client) GetBlock(ctx context.Context, blkHash common.Hash) (*types.Block, error) {
	var resp GetBlockResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("blocks/%s", utils.EncodeToHex(blkHash.Bytes())),
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
