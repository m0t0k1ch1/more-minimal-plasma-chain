package client

import "github.com/m0t0k1ch1/more-minimal-plasma-chain/app"

type ErrorResponse struct {
	State  string     `json:"state"`
	Result *app.Error `json:"result"`
}
