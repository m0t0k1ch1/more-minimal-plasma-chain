package app

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type Config struct {
	Port      int                  `json:"port"`
	DB        DBConfig             `json:"db"`
	Operator  OperatorConfig       `json:"operator"`
	RootChain core.RootChainConfig `json:"rootchain"`
	Heartbeat HeartbeatConfig      `json:"heartbeat"`
}

type DBConfig struct {
	Dir string `json:"dir"`
}

type OperatorConfig struct {
	PrivateKeyStr string `json:"privkey"`
}

func (conf OperatorConfig) PrivateKey() (*ecdsa.PrivateKey, error) {
	return utils.HexToPrivateKey(conf.PrivateKeyStr)
}

type HeartbeatConfig struct {
	IsEnabled   bool `json:"enabled"`
	IntervalInt int  `json:"interval"`
}

func (conf HeartbeatConfig) Interval() (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%ds", conf.IntervalInt))
}
