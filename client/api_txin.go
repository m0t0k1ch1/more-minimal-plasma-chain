package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PutTxInResponse struct {
	State  string `json:"state"`
	Result struct {
		Pos *big.Int `json:"pos"`
	} `json:"result"`
}

func (c *Client) PutTxIn(ctx context.Context, txInPos types.Position, confSig types.Signature) (types.Position, error) {
	v := url.Values{}
	v.Set("confsig", confSig.Hex())

	var resp PutTxInResponse
	if err := c.doAPI(
		ctx,
		http.MethodPut,
		fmt.Sprintf("txins/%s", txInPos.String()),
		v,
		&resp,
	); err != nil {
		return types.NullPosition, err
	}

	return types.NewPosition(resp.Result.Pos), nil
}
