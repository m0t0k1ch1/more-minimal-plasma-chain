package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
)

type PostTxResponse struct {
	State  string    `json:"state"`
	Result *types.Tx `json:"result"`
}

func (client *Client) PostTx(ctx context.Context, tx *types.Tx) (*types.Tx, error) {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("tx", hexutil.Encode(txBytes))

	var resp PostTxResponse
	if err := client.doAPI(ctx, http.MethodPost, "txes", v, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}
