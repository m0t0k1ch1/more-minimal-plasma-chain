package main

import "github.com/urfave/cli"

const (
	nullAddressStr = "0x0000000000000000000000000000000000000000"
	nullHashStr    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	rpcFlag = cli.StringFlag{Name: "rpc", Value: "http://127.0.0.1:7545"}
	wsFlag  = cli.StringFlag{Name: "ws", Value: "ws://127.0.0.1:7545"}
	apiFlag = cli.StringFlag{Name: "api", Value: "http://127.0.0.1:1323"}

	// TODO: use apiFlag instead of these flags
	hostFlag = cli.StringFlag{Name: "host", Value: "127.0.0.1"}
	portFlag = cli.Uint64Flag{Name: "port", Value: 1323}

	hashFlag    = cli.StringFlag{Name: "hash", Value: ""}
	addrFlag    = cli.StringFlag{Name: "addr", Value: nullAddressStr}
	privKeyFlag = cli.StringFlag{Name: "privkey", Value: ""}
	rootFlag    = cli.StringFlag{Name: "root", Value: nullHashStr}
	txFlag      = cli.StringFlag{Name: "tx", Value: ""}
	blkNumFlag  = cli.Uint64Flag{Name: "blknum", Value: 0}
	txIndexFlag = cli.Uint64Flag{Name: "txindex", Value: 0}
	oIndexFlag  = cli.Uint64Flag{Name: "oindex", Value: 0}
	ownerFlag   = cli.StringFlag{Name: "owner", Value: nullAddressStr}
	amountFlag  = cli.Uint64Flag{Name: "amount", Value: 0}
	encodedFlag = cli.BoolFlag{Name: "encoded"}
)
