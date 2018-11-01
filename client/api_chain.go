package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type GetChainResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockHashStr string `json:"blkhash"`
	} `json:"result"`
}

func (client *Client) GetChain(ctx context.Context, blkNum uint64) ([]byte, error) {
	var resp GetChainResponse
	if err := client.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("chain/%d", blkNum),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.BlockHashStr)
}
