package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type GetAddressUTXOsResponse struct {
	*ResponseBase
	Result struct {
		UTXOs []types.Position `json:"utxos"`
	} `json:"result"`
}

func (c *Client) GetAddressUTXOs(ctx context.Context, addr common.Address) ([]types.Position, error) {
	var resp GetAddressUTXOsResponse
	if err := c.doAPI(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/addresses/%s/utxos", utils.AddressToHex(addr)),
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return resp.Result.UTXOs, nil
}
