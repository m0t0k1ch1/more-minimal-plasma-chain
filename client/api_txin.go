package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PutTxInResponse struct {
	*ResponseBase
	Result struct{} `json:"result"`
}

func (c *Client) PutTxIn(ctx context.Context, txInPos types.Position, confSig types.Signature) error {
	v := url.Values{}
	v.Set("confsig", confSig.Hex())

	var resp PutTxInResponse
	if err := c.doAPI(
		ctx,
		http.MethodPut,
		fmt.Sprintf("txins/%d", txInPos),
		v,
		&resp,
	); err != nil {
		return err
	}

	return nil
}
