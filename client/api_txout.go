package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PutTxOutResponse struct {
	*ResponseBase
	Result struct{} `json:"result"`
}

func (c *Client) PutTxOut(ctx context.Context, txOutPos types.Position, isExited bool) error {
	v := url.Values{}
	v.Set("exited", fmt.Sprintf("%t", isExited))

	var resp PutTxOutResponse
	if err := c.doAPI(
		ctx,
		http.MethodPut,
		fmt.Sprintf("txouts/%d", txOutPos),
		v,
		&resp,
	); err != nil {
		return err
	}

	return nil
}
