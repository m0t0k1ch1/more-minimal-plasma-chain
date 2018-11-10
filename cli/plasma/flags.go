package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	confFlag = cli.StringFlag{Name: "conf", Value: "config.json"}

	addressFlag = cli.StringFlag{Name: "address", Value: nullAddressStr}
	amountFlag  = cli.StringFlag{Name: "amount", Value: "0"}
	encodedFlag = cli.BoolFlag{Name: "encoded"}
	indexFlag   = cli.StringFlag{Name: "index", Value: "0"}
	numFlag     = cli.StringFlag{Name: "num", Value: "0"}
	posFlag     = cli.StringFlag{Name: "pos", Value: "0"}
	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	txFlag      = cli.StringFlag{Name: "tx", Value: ""}
	vsPosFlag   = cli.StringFlag{Name: "vspos", Value: "0"}
)

func flags(fs ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		confFlag,
	}, fs...)
}
