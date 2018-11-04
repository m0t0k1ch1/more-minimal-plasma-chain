package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type GetChainResponse struct {
	State  string `json:"state"`
	Result struct {
		BlockHashStr string `json:"blkhash"`
	} `json:"result"`
}

func (c *Client) GetChain(ctx context.Context, blkNum *big.Int) (common.Hash, error) {
	var resp GetChainResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("chain/%s", blkNum.String()),
		nil,
		&resp,
	); err != nil {
		return types.NullHash, err
	}

	return utils.HexToHash(resp.Result.BlockHashStr)
}
