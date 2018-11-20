package client

import (
	"context"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type PostDepositResponse struct {
	*ResponseBase
	Result struct {
		BlockNumber uint64 `json:"blknum"`
	} `json:"result"`
}

func (c *Client) PostDeposit(ctx context.Context, ownerAddr common.Address, amount *big.Int) (uint64, error) {
	v := url.Values{}
	v.Set("owner", utils.AddressToHex(ownerAddr))
	v.Set("amount", amount.String())

	var resp PostDepositResponse
	if err := c.doAPI(
		ctx,
		http.MethodPost,
		"deposits",
		v,
		&resp,
	); err != nil {
		return 0, err
	}

	return resp.Result.BlockNumber, nil
}
