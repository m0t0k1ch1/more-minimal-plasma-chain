package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	rpcFlag      = cli.StringFlag{Name: "rpc", Value: "http://127.0.0.1:7545"}
	wsFlag       = cli.StringFlag{Name: "ws", Value: "ws://127.0.0.1:7545"}
	contractFlag = cli.StringFlag{Name: "contract", Value: nullAddressStr}

	apiFlag = cli.StringFlag{Name: "api", Value: "http://127.0.0.1:1323"}

	rootFlag    = cli.StringFlag{Name: "root", Value: nullHashStr}
	txFlag      = cli.StringFlag{Name: "tx", Value: ""}
	blkNumFlag  = cli.StringFlag{Name: "blknum", Value: "0"}
	txIndexFlag = cli.StringFlag{Name: "txindex", Value: "0"}
	oIndexFlag  = cli.StringFlag{Name: "oindex", Value: "0"}
	ownerFlag   = cli.StringFlag{Name: "owner", Value: nullAddressStr}
	amountFlag  = cli.StringFlag{Name: "amount", Value: "0"}

	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	hashFlag    = cli.StringFlag{Name: "hash", Value: nullHashStr}
	encodedFlag = cli.BoolFlag{Name: "encoded"}
)
