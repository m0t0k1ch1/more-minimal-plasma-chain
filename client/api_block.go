package client

import (
	"context"
	"net/http"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PostBlockResponse struct {
	State  string              `json:"state"`
	Result *types.BlockSummary `json:"result"`
}

func (client *Client) PostBlock(ctx context.Context) (*types.BlockSummary, error) {
	var resp PostBlockResponse
	if err := client.doAPI(ctx, http.MethodPost, "blocks", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}
