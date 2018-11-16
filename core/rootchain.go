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
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/core/types"
	"github.com/m0t0k1ch1/more-minimal-plasma-chain/utils"
)

const (
	DefaultExitBondAmount = 123456789
)

type RootChainConfig struct {
	RPC        string `json:"rpc"`
	WS         string `json:"ws"`
	AddressStr string `json:"address"`
}

func (conf RootChainConfig) Address() (common.Address, error) {
	if ok := utils.IsHexAddress(conf.AddressStr); !ok {
		return types.NullAddress, fmt.Errorf("invalid root chain address")
	}

	return utils.HexToAddress(conf.AddressStr), nil
}

type RootChain struct {
	config    RootChainConfig
	address   common.Address
	abi       abi.ABI
	rpcClient *ethclient.Client
	wsClient  *rpc.Client
	contract  *bind.BoundContract
}

func NewRootChain(conf RootChainConfig) (*RootChain, error) {
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
	addr, err := rc.config.Address()
	if err != nil {
		return err
	}
	rc.address = addr
	return nil
}

func (rc *RootChain) initABI() error {
	abi, err := abi.JSON(strings.NewReader(RootChainABI))
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

func (rc *RootChain) CurrentPlasmaBlockNumber() (uint64, error) {
	blkNum := new(*big.Int)
	if err := rc.contract.Call(nil, blkNum, "currentPlasmaBlockNumber"); err != nil {
		return 0, err
	}

	return (*blkNum).Uint64(), nil
}

func (rc *RootChain) PlasmaExits(txOutPos types.Position) (types.Exit, error) {
	exit := new(types.Exit)
	if err := rc.contract.Call(nil, exit, "plasmaExits", new(big.Int).SetUint64(txOutPos.Uint64())); err != nil {
		return types.Exit{}, err
	}

	return *exit, nil
}

func (rc *RootChain) CommitPlasmaBlockRoot(a *types.Account, rootHash common.Hash) (*gethtypes.Transaction, error) {
	return rc.contract.Transact(a.TransactOpts(), "commitPlasmaBlockRoot", rootHash)
}

func (rc *RootChain) Deposit(a *types.Account, amount uint64) (*gethtypes.Transaction, error) {
	opts := a.TransactOpts()
	opts.Value = new(big.Int).SetUint64(amount)

	return rc.contract.Transact(opts, "deposit")
}

func (rc *RootChain) StartExit(a *types.Account, txOutPos types.Position, tx *types.Tx, txProofBytes []byte) (*gethtypes.Transaction, error) {
	blkNum, txIndex, outIndex := types.ParseTxOutPosition(txOutPos)

	encodedTxBytes, err := tx.Encode()
	if err != nil {
		return nil, err
	}

	sigsBytes, err := tx.SignaturesBytes()
	if err != nil {
		return nil, err
	}

	confSigsBytes, err := tx.ConfirmationSignaturesBytes()
	if err != nil {
		return nil, err
	}

	opts := a.TransactOpts()
	opts.Value = big.NewInt(DefaultExitBondAmount)

	return rc.contract.Transact(
		opts,
		"startExit",
		new(big.Int).SetUint64(blkNum), new(big.Int).SetUint64(txIndex), new(big.Int).SetUint64(outIndex),
		encodedTxBytes,
		txProofBytes,
		sigsBytes, confSigsBytes,
	)
}

func (rc *RootChain) ChallengeExit(a *types.Account, txOutPos types.Position, spendingTx *types.Tx, spendingInIndex uint64) (*gethtypes.Transaction, error) {
	blkNum, txIndex, outIndex := types.ParseTxOutPosition(txOutPos)

	encodedSpendingTxBytes, err := spendingTx.Encode()
	if err != nil {
		return nil, err
	}

	spendingTxIn := spendingTx.GetInput(spendingInIndex)
	if spendingTxIn == nil {
		return nil, ErrTxInNotFound
	}
	if spendingTxIn.ConfirmationSignature.IsNull() {
		return nil, ErrNullConfirmationSignature
	}

	return rc.contract.Transact(
		a.TransactOpts(),
		"challengeExit",
		new(big.Int).SetUint64(blkNum), new(big.Int).SetUint64(txIndex), new(big.Int).SetUint64(outIndex),
		encodedSpendingTxBytes,
		spendingTxIn.ConfirmationSignature.Bytes(),
	)
}

func (rc *RootChain) ProcessExits(a *types.Account) (*gethtypes.Transaction, error) {
	return rc.contract.Transact(a.TransactOpts(), "processExits")
}

type RootChainDepositCreated struct {
	Owner        common.Address
	Amount       *big.Int
	DepositBlock *big.Int
	Raw          gethtypes.Log
}

func (rc *RootChain) WatchDepositCreated(ctx context.Context, sink chan<- *RootChainDepositCreated) (event.Subscription, error) {
	logs := make(chan gethtypes.Log)
	arg := map[string]interface{}{
		"fromBlock": "0x0",
		"toBlock":   "latest",
		"address":   rc.config.AddressStr,
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
				event := new(RootChainDepositCreated)
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

type RootChainExitStarted struct {
	Owner        common.Address
	UtxoPosition *big.Int
	Amount       *big.Int
	Raw          gethtypes.Log
}

func (rc *RootChain) WatchExitStarted(ctx context.Context, sink chan<- *RootChainExitStarted) (event.Subscription, error) {
	logs := make(chan gethtypes.Log)
	arg := map[string]interface{}{
		"fromBlock": "0x0",
		"toBlock":   "latest",
		"address":   rc.config.AddressStr,
		"topics": []interface{}{
			rc.abi.Events["ExitStarted"].Id().Hex(),
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
				event := new(RootChainExitStarted)
				if err := rc.contract.UnpackLog(event, "ExitStarted", log); err != nil {
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

func (rc *RootChain) Ping() error {
	if err := rc.wsClient.Call(nil, "eth_getFilterLogs", "0x0"); err != nil {
		if err.Error() == "filter not found" {
			return nil
		}
		return err
	}
	return nil
}
