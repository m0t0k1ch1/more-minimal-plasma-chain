package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/contract"
)

type RootChain struct {
	config     *RootChainConfig
	address    common.Address
	abi        abi.ABI
	httpClient *ethclient.Client
	wsClient   *rpc.Client
	contract   *bind.BoundContract
}

func NewRootChain(conf *RootChainConfig) (*RootChain, error) {
	if ok := common.IsHexAddress(conf.Address); !ok {
		return nil, fmt.Errorf("invalid root chain address")
	}
	addr := common.HexToAddress(conf.Address)

	abi, err := abi.JSON(strings.NewReader(contract.RootChainABI))
	if err != nil {
		return nil, err
	}

	httpClient, err := ethclient.Dial(conf.RPC.HTTP)
	if err != nil {
		return nil, err
	}

	wsClient, err := rpc.Dial(conf.RPC.WS)
	if err != nil {
		return nil, err
	}

	c := bind.NewBoundContract(addr, abi, httpClient, httpClient, httpClient)

	return &RootChain{
		config:     conf,
		address:    addr,
		abi:        abi,
		httpClient: httpClient,
		wsClient:   wsClient,
		contract:   c,
	}, nil
}

// NOTICE:
// By right, we should use rc.contract.WatchDepositCreated() instead of this func,
// but we use this func because ganache cannot parse web3.eth.subscribe request created by rc.contract.
func (rc *RootChain) WatchDepositCreated(ctx context.Context, sink chan<- *contract.RootChainDepositCreated) (event.Subscription, error) {
	logs := make(chan types.Log)
	arg := map[string]interface{}{
		"fromBlock": "0x0",
		"toBlock":   "latest",
		"address":   rc.config.Address,
		"topics": []interface{}{
			rc.abi.Events["DepositCreated"].Id().Hex(),
			nil,
		},
	}

	sub, err := rc.wsClient.EthSubscribe(ctx, logs, "logs", arg)
	if err != nil {
		return nil, err
	}

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				event := new(contract.RootChainDepositCreated)
				if err := rc.contract.UnpackLog(event, "DepositCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}

			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
