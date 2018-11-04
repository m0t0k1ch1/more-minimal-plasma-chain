package core

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/contract"
	mmpctypes "github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

type RootChainConfig struct {
	RPC     string `json:"rpc"`
	WS      string `json:"ws"`
	Address string `json:"address"`
}

type RootChain struct {
	config    *RootChainConfig
	address   common.Address
	abi       abi.ABI
	rpcClient *ethclient.Client
	wsClient  *rpc.Client
	contract  *bind.BoundContract
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
	if err := rc.initRPCClient(); err != nil {
		return nil, err
	}
	if err := rc.initWSClient(); err != nil {
		return nil, err
	}
	rc.initContract()

	return rc, nil
}

func (rc *RootChain) initAddress() error {
	if ok := utils.IsHexAddress(rc.config.Address); !ok {
		return fmt.Errorf("invalid root chain address")
	}
	rc.address = utils.HexToAddress(rc.config.Address)
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

func (rc *RootChain) initRPCClient() error {
	rpcClient, err := ethclient.Dial(rc.config.RPC)
	if err != nil {
		return err
	}
	rc.rpcClient = rpcClient
	return nil
}

func (rc *RootChain) initWSClient() error {
	wsClient, err := rpc.Dial(rc.config.WS)
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
		rc.rpcClient,
		rc.rpcClient,
		rc.rpcClient,
	)
}

func (rc *RootChain) CurrentPlasmaBlockNumber() (*big.Int, error) {
	blkNum := new(*big.Int)
	if err := rc.contract.Call(nil, blkNum, "currentPlasmaBlockNumber"); err != nil {
		return nil, err
	}
	return *blkNum, nil
}

func (rc *RootChain) CommitPlasmaBlockRoot(a *mmpctypes.Account, rootHash common.Hash) (*gethtypes.Transaction, error) {
	return rc.contract.Transact(a.TransactOpts(), "commitPlasmaBlockRoot", rootHash)
}

// NOTICE:
// By right, we should use contract.RootChain.WatchDepositCreated() instead of this func,
// but we use this func because ganache cannot parse web3.eth.subscribe request created by contract.RootChain.
func (rc *RootChain) WatchDepositCreated(ctx context.Context, sink chan<- *contract.RootChainDepositCreated) (event.Subscription, error) {
	logs := make(chan gethtypes.Log)
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
