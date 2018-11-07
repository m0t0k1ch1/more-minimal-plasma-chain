package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	confFlag = cli.StringFlag{Name: "conf", Value: "config.json"}

	blkRootHashFlag = cli.StringFlag{Name: "blkroot", Value: nullHashStr}
	blkNumFlag      = cli.StringFlag{Name: "blknum", Value: "0"}

	txFlag       = cli.StringFlag{Name: "tx", Value: ""}
	txPosFlag    = cli.StringFlag{Name: "txpos", Value: "0"}
	txOutPosFlag = cli.StringFlag{Name: "txopos", Value: "0"}
	iIndexFlag   = cli.StringFlag{Name: "iindex", Value: "0"}
	oIndexFlag   = cli.StringFlag{Name: "oindex", Value: "0"}
	toAddrFlag   = cli.StringFlag{Name: "to", Value: nullAddressStr}
	amountFlag   = cli.StringFlag{Name: "amount", Value: "0"}

	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	encodedFlag = cli.BoolFlag{Name: "encoded"}
)

func flags(fs ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		confFlag,
	}, fs...)
}
