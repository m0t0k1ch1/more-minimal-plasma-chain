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

const RootChainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"plasmaExits\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"isStarted\",\"type\":\"bool\"},{\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentPlasmaBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"plasmaBlocks\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CHALLENGE_PERIOD\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EXIT_BOND\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"depositBlock\",\"type\":\"uint256\"}],\"name\":\"DepositCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"PlasmaBlockRootCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"utxoPosition\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ExitStarted\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_root\",\"type\":\"bytes32\"}],\"name\":\"commitPlasmaBlockRoot\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_utxoBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_utxoTxIndex\",\"type\":\"uint256\"},{\"name\":\"_utxoOutputIndex\",\"type\":\"uint256\"},{\"name\":\"_encodedTx\",\"type\":\"bytes\"},{\"name\":\"_txInclusionProof\",\"type\":\"bytes\"},{\"name\":\"_txSignatures\",\"type\":\"bytes\"},{\"name\":\"_txConfirmationSignatures\",\"type\":\"bytes\"}],\"name\":\"startExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_exitingUtxoBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_exitingUtxoTxIndex\",\"type\":\"uint256\"},{\"name\":\"_exitingUtxoOutputIndex\",\"type\":\"uint256\"},{\"name\":\"_encodedSpendingTx\",\"type\":\"bytes\"},{\"name\":\"_spendingTxConfirmationSignature\",\"type\":\"bytes\"}],\"name\":\"challengeExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"processExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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
	blkNum := new(uint64)
	if err := rc.contract.Call(nil, blkNum, "currentPlasmaBlockNumber"); err != nil {
		return 0, err
	}

	return *blkNum, nil
}

func (rc *RootChain) PlasmaExits(txOutPos types.Position) (types.Exit, error) {
	exit := new(types.Exit)
	if err := rc.contract.Call(nil, exit, "plasmaExits", txOutPos.Uint64()); err != nil {
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

	return rc.contract.Transact(opts, "startExit", blkNum, txIndex, outIndex, encodedTxBytes, txProofBytes, sigsBytes, confSigsBytes)
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

	return rc.contract.Transact(a.TransactOpts(), "challengeExit", blkNum, txIndex, outIndex, encodedSpendingTxBytes, spendingTxIn.ConfirmationSignature.Bytes())
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
