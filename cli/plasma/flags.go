package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	rpcFlag          = cli.StringFlag{Name: "rpc", Value: "http://127.0.0.1:7545"}
	wsFlag           = cli.StringFlag{Name: "ws", Value: "ws://127.0.0.1:7545"}
	contractAddrFlag = cli.StringFlag{Name: "contract", Value: nullAddressStr}

	apiFlag = cli.StringFlag{Name: "api", Value: "http://127.0.0.1:1323"}

	blkHashFlag     = cli.StringFlag{Name: "blkhash", Value: nullHashStr}
	blkRootHashFlag = cli.StringFlag{Name: "blkroot", Value: nullHashStr}
	blkNumFlag      = cli.StringFlag{Name: "blknum", Value: "0"}

	txFlag     = cli.StringFlag{Name: "tx", Value: ""}
	txHashFlag = cli.StringFlag{Name: "txhash", Value: nullHashStr}
	iIndexFlag = cli.StringFlag{Name: "iindex", Value: "0"}
	oIndexFlag = cli.StringFlag{Name: "oindex", Value: "0"}
	toAddrFlag = cli.StringFlag{Name: "to", Value: nullAddressStr}
	amountFlag = cli.StringFlag{Name: "amount", Value: "0"}

	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	encodedFlag = cli.BoolFlag{Name: "encoded"}
)
