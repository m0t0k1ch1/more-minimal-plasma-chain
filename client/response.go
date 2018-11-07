package client

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/app"

type ResponseBase struct {
	State string `json:"state"`
}

type ErrorResponse struct {
	*ResponseBase
	Result *app.Error `json:"result"`
}
