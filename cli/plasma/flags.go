package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	confFlag = cli.StringFlag{Name: "conf", Value: "config.json"}

	// common
	numFlag     = cli.StringFlag{Name: "num", Value: "0"}
	posFlag     = cli.StringFlag{Name: "pos", Value: "0"}
	indexFlag   = cli.StringFlag{Name: "index", Value: "0"}
	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	encodedFlag = cli.BoolFlag{Name: "encoded"}

	// for block
	rootFlag = cli.StringFlag{Name: "root", Value: nullHashStr}

	// for tx
	txFlag     = cli.StringFlag{Name: "tx", Value: ""}
	toFlag     = cli.StringFlag{Name: "to", Value: nullAddressStr}
	amountFlag = cli.StringFlag{Name: "amount", Value: "0"}

	// for challenge
	vsPosFlag = cli.StringFlag{Name: "vspos", Value: "0"}
)

func flags(fs ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		confFlag,
	}, fs...)
}
