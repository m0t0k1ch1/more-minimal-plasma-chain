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
	rc := &RootChain{
		config: conf,
	}

	if err := rc.initAddress(); err != nil {
		return nil, err
	}
	if err := rc.initABI(); err != nil {
		return nil, err
	}
	if err := rc.initHTTPClient(); err != nil {
		return nil, err
	}
	if err := rc.initWSClient(); err != nil {
		return nil, err
	}
	rc.initContract()

	return rc, nil
}

func (rc *RootChain) initAddress() error {
	if ok := common.IsHexAddress(rc.config.Address); !ok {
		return fmt.Errorf("invalid root chain address")
	}
	rc.address = common.HexToAddress(rc.config.Address)
	return nil
}

func (rc *RootChain) initABI() error {
	abi, err := abi.JSON(strings.NewReader(contract.RootChainABI))
	if err != nil {
		return err
	}
	rc.abi = abi
	return nil
}

func (rc *RootChain) initHTTPClient() error {
	httpClient, err := ethclient.Dial(rc.config.RPC.HTTP)
	if err != nil {
		return err
	}
	rc.httpClient = httpClient
	return nil
}

func (rc *RootChain) initWSClient() error {
	wsClient, err := rpc.Dial(rc.config.RPC.WS)
	if err != nil {
		return err
	}
	rc.wsClient = wsClient
	return nil
}

func (rc *RootChain) initContract() {
	rc.contract = bind.NewBoundContract(
		rc.address,
		rc.abi,
		rc.httpClient,
		rc.httpClient,
		rc.httpClient,
	)
}

// NOTICE:
// By right, we should use contract.RootChain.WatchDepositCreated() instead of this func,
// but we use this func because ganache cannot parse web3.eth.subscribe request created by contract.RootChain.
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
