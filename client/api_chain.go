package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type GetChainResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockHashStr string `json:"blkhash"`
	} `json:"result"`
}

func (c *Client) GetChain(ctx context.Context, blkNum *big.Int) ([]byte, error) {
	var resp GetChainResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("chain/%s", blkNum.String()),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return utils.DecodeHex(resp.Result.BlockHashStr)
}
