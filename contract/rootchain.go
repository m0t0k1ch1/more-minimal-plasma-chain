// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// RootChainABI is the input ABI used to generate the binding from.
const RootChainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"plasmaExits\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"isStarted\",\"type\":\"bool\"},{\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentPlasmaBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"plasmaBlocks\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CHALLENGE_PERIOD\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EXIT_BOND\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"depositBlock\",\"type\":\"uint256\"}],\"name\":\"DepositCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"PlasmaBlockRootCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"utxoPosition\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ExitStarted\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_root\",\"type\":\"bytes32\"}],\"name\":\"commitPlasmaBlockRoot\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_utxoBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_utxoTxIndex\",\"type\":\"uint256\"},{\"name\":\"_utxoOutputIndex\",\"type\":\"uint256\"},{\"name\":\"_encodedTx\",\"type\":\"bytes\"},{\"name\":\"_txInclusionProof\",\"type\":\"bytes\"},{\"name\":\"_txSignatures\",\"type\":\"bytes\"},{\"name\":\"_txConfirmationSignatures\",\"type\":\"bytes\"}],\"name\":\"startExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_exitingUtxoBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_exitingUtxoTxIndex\",\"type\":\"uint256\"},{\"name\":\"_exitingUtxoOutputIndex\",\"type\":\"uint256\"},{\"name\":\"_encodedSpendingTx\",\"type\":\"bytes\"},{\"name\":\"_spendingTxConfirmationSignature\",\"type\":\"bytes\"}],\"name\":\"challengeExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"processExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RootChain is an auto generated Go binding around an Ethereum contract.
type RootChain struct {
	RootChainCaller     // Read-only binding to the contract
	RootChainTransactor // Write-only binding to the contract
	RootChainFilterer   // Log filterer for contract events
}

// RootChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootChainSession struct {
	Contract     *RootChain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootChainCallerSession struct {
	Contract *RootChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RootChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootChainTransactorSession struct {
	Contract     *RootChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RootChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootChainRaw struct {
	Contract *RootChain // Generic contract binding to access the raw methods on
}

// RootChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootChainCallerRaw struct {
	Contract *RootChainCaller // Generic read-only contract binding to access the raw methods on
}

// RootChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootChainTransactorRaw struct {
	Contract *RootChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootChain creates a new instance of RootChain, bound to a specific deployed contract.
func NewRootChain(address common.Address, backend bind.ContractBackend) (*RootChain, error) {
	contract, err := bindRootChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootChain{RootChainCaller: RootChainCaller{contract: contract}, RootChainTransactor: RootChainTransactor{contract: contract}, RootChainFilterer: RootChainFilterer{contract: contract}}, nil
}

// NewRootChainCaller creates a new read-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainCaller(address common.Address, caller bind.ContractCaller) (*RootChainCaller, error) {
	contract, err := bindRootChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainCaller{contract: contract}, nil
}

// NewRootChainTransactor creates a new write-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainTransactor(address common.Address, transactor bind.ContractTransactor) (*RootChainTransactor, error) {
	contract, err := bindRootChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainTransactor{contract: contract}, nil
}

// NewRootChainFilterer creates a new log filterer instance of RootChain, bound to a specific deployed contract.
func NewRootChainFilterer(address common.Address, filterer bind.ContractFilterer) (*RootChainFilterer, error) {
	contract, err := bindRootChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootChainFilterer{contract: contract}, nil
}

// bindRootChain binds a generic wrapper to an already deployed contract.
func bindRootChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.RootChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transact(opts, method, params...)
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_RootChain *RootChainCaller) CHALLENGEPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "CHALLENGE_PERIOD")
	return *ret0, err
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_RootChain *RootChainSession) CHALLENGEPERIOD() (*big.Int, error) {
	return _RootChain.Contract.CHALLENGEPERIOD(&_RootChain.CallOpts)
}

// CHALLENGEPERIOD is a free data retrieval call binding the contract method 0xc3a079ed.
//
// Solidity: function CHALLENGE_PERIOD() constant returns(uint256)
func (_RootChain *RootChainCallerSession) CHALLENGEPERIOD() (*big.Int, error) {
	return _RootChain.Contract.CHALLENGEPERIOD(&_RootChain.CallOpts)
}

// EXITBOND is a free data retrieval call binding the contract method 0xde5a9f09.
//
// Solidity: function EXIT_BOND() constant returns(uint256)
func (_RootChain *RootChainCaller) EXITBOND(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "EXIT_BOND")
	return *ret0, err
}

// EXITBOND is a free data retrieval call binding the contract method 0xde5a9f09.
//
// Solidity: function EXIT_BOND() constant returns(uint256)
func (_RootChain *RootChainSession) EXITBOND() (*big.Int, error) {
	return _RootChain.Contract.EXITBOND(&_RootChain.CallOpts)
}

// EXITBOND is a free data retrieval call binding the contract method 0xde5a9f09.
//
// Solidity: function EXIT_BOND() constant returns(uint256)
func (_RootChain *RootChainCallerSession) EXITBOND() (*big.Int, error) {
	return _RootChain.Contract.EXITBOND(&_RootChain.CallOpts)
}

// CurrentPlasmaBlockNumber is a free data retrieval call binding the contract method 0xa732668e.
//
// Solidity: function currentPlasmaBlockNumber() constant returns(uint256)
func (_RootChain *RootChainCaller) CurrentPlasmaBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "currentPlasmaBlockNumber")
	return *ret0, err
}

// CurrentPlasmaBlockNumber is a free data retrieval call binding the contract method 0xa732668e.
//
// Solidity: function currentPlasmaBlockNumber() constant returns(uint256)
func (_RootChain *RootChainSession) CurrentPlasmaBlockNumber() (*big.Int, error) {
	return _RootChain.Contract.CurrentPlasmaBlockNumber(&_RootChain.CallOpts)
}

// CurrentPlasmaBlockNumber is a free data retrieval call binding the contract method 0xa732668e.
//
// Solidity: function currentPlasmaBlockNumber() constant returns(uint256)
func (_RootChain *RootChainCallerSession) CurrentPlasmaBlockNumber() (*big.Int, error) {
	return _RootChain.Contract.CurrentPlasmaBlockNumber(&_RootChain.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_RootChain *RootChainCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_RootChain *RootChainSession) Operator() (common.Address, error) {
	return _RootChain.Contract.Operator(&_RootChain.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_RootChain *RootChainCallerSession) Operator() (common.Address, error) {
	return _RootChain.Contract.Operator(&_RootChain.CallOpts)
}

// PlasmaBlocks is a free data retrieval call binding the contract method 0xb2153e6b.
//
// Solidity: function plasmaBlocks( uint256) constant returns(root bytes32, timestamp uint256)
func (_RootChain *RootChainCaller) PlasmaBlocks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Root      [32]byte
	Timestamp *big.Int
}, error) {
	ret := new(struct {
		Root      [32]byte
		Timestamp *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "plasmaBlocks", arg0)
	return *ret, err
}

// PlasmaBlocks is a free data retrieval call binding the contract method 0xb2153e6b.
//
// Solidity: function plasmaBlocks( uint256) constant returns(root bytes32, timestamp uint256)
func (_RootChain *RootChainSession) PlasmaBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Timestamp *big.Int
}, error) {
	return _RootChain.Contract.PlasmaBlocks(&_RootChain.CallOpts, arg0)
}

// PlasmaBlocks is a free data retrieval call binding the contract method 0xb2153e6b.
//
// Solidity: function plasmaBlocks( uint256) constant returns(root bytes32, timestamp uint256)
func (_RootChain *RootChainCallerSession) PlasmaBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Timestamp *big.Int
}, error) {
	return _RootChain.Contract.PlasmaBlocks(&_RootChain.CallOpts, arg0)
}

// PlasmaExits is a free data retrieval call binding the contract method 0x7d53e633.
//
// Solidity: function plasmaExits( uint256) constant returns(owner address, amount uint256, isStarted bool, isValid bool)
func (_RootChain *RootChainCaller) PlasmaExits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner     common.Address
	Amount    *big.Int
	IsStarted bool
	IsValid   bool
}, error) {
	ret := new(struct {
		Owner     common.Address
		Amount    *big.Int
		IsStarted bool
		IsValid   bool
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "plasmaExits", arg0)
	return *ret, err
}

// PlasmaExits is a free data retrieval call binding the contract method 0x7d53e633.
//
// Solidity: function plasmaExits( uint256) constant returns(owner address, amount uint256, isStarted bool, isValid bool)
func (_RootChain *RootChainSession) PlasmaExits(arg0 *big.Int) (struct {
	Owner     common.Address
	Amount    *big.Int
	IsStarted bool
	IsValid   bool
}, error) {
	return _RootChain.Contract.PlasmaExits(&_RootChain.CallOpts, arg0)
}

// PlasmaExits is a free data retrieval call binding the contract method 0x7d53e633.
//
// Solidity: function plasmaExits( uint256) constant returns(owner address, amount uint256, isStarted bool, isValid bool)
func (_RootChain *RootChainCallerSession) PlasmaExits(arg0 *big.Int) (struct {
	Owner     common.Address
	Amount    *big.Int
	IsStarted bool
	IsValid   bool
}, error) {
	return _RootChain.Contract.PlasmaExits(&_RootChain.CallOpts, arg0)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0x404f7d66.
//
// Solidity: function challengeExit(_exitingUtxoBlockNumber uint256, _exitingUtxoTxIndex uint256, _exitingUtxoOutputIndex uint256, _encodedSpendingTx bytes, _spendingTxConfirmationSignature bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeExit(opts *bind.TransactOpts, _exitingUtxoBlockNumber *big.Int, _exitingUtxoTxIndex *big.Int, _exitingUtxoOutputIndex *big.Int, _encodedSpendingTx []byte, _spendingTxConfirmationSignature []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeExit", _exitingUtxoBlockNumber, _exitingUtxoTxIndex, _exitingUtxoOutputIndex, _encodedSpendingTx, _spendingTxConfirmationSignature)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0x404f7d66.
//
// Solidity: function challengeExit(_exitingUtxoBlockNumber uint256, _exitingUtxoTxIndex uint256, _exitingUtxoOutputIndex uint256, _encodedSpendingTx bytes, _spendingTxConfirmationSignature bytes) returns()
func (_RootChain *RootChainSession) ChallengeExit(_exitingUtxoBlockNumber *big.Int, _exitingUtxoTxIndex *big.Int, _exitingUtxoOutputIndex *big.Int, _encodedSpendingTx []byte, _spendingTxConfirmationSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeExit(&_RootChain.TransactOpts, _exitingUtxoBlockNumber, _exitingUtxoTxIndex, _exitingUtxoOutputIndex, _encodedSpendingTx, _spendingTxConfirmationSignature)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0x404f7d66.
//
// Solidity: function challengeExit(_exitingUtxoBlockNumber uint256, _exitingUtxoTxIndex uint256, _exitingUtxoOutputIndex uint256, _encodedSpendingTx bytes, _spendingTxConfirmationSignature bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeExit(_exitingUtxoBlockNumber *big.Int, _exitingUtxoTxIndex *big.Int, _exitingUtxoOutputIndex *big.Int, _encodedSpendingTx []byte, _spendingTxConfirmationSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeExit(&_RootChain.TransactOpts, _exitingUtxoBlockNumber, _exitingUtxoTxIndex, _exitingUtxoOutputIndex, _encodedSpendingTx, _spendingTxConfirmationSignature)
}

// CommitPlasmaBlockRoot is a paid mutator transaction binding the contract method 0x7f052e03.
//
// Solidity: function commitPlasmaBlockRoot(_root bytes32) returns()
func (_RootChain *RootChainTransactor) CommitPlasmaBlockRoot(opts *bind.TransactOpts, _root [32]byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "commitPlasmaBlockRoot", _root)
}

// CommitPlasmaBlockRoot is a paid mutator transaction binding the contract method 0x7f052e03.
//
// Solidity: function commitPlasmaBlockRoot(_root bytes32) returns()
func (_RootChain *RootChainSession) CommitPlasmaBlockRoot(_root [32]byte) (*types.Transaction, error) {
	return _RootChain.Contract.CommitPlasmaBlockRoot(&_RootChain.TransactOpts, _root)
}

// CommitPlasmaBlockRoot is a paid mutator transaction binding the contract method 0x7f052e03.
//
// Solidity: function commitPlasmaBlockRoot(_root bytes32) returns()
func (_RootChain *RootChainTransactorSession) CommitPlasmaBlockRoot(_root [32]byte) (*types.Transaction, error) {
	return _RootChain.Contract.CommitPlasmaBlockRoot(&_RootChain.TransactOpts, _root)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainSession) Deposit() (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainTransactorSession) Deposit() (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts)
}

// ProcessExits is a paid mutator transaction binding the contract method 0x58189e03.
//
// Solidity: function processExits() returns()
func (_RootChain *RootChainTransactor) ProcessExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "processExits")
}

// ProcessExits is a paid mutator transaction binding the contract method 0x58189e03.
//
// Solidity: function processExits() returns()
func (_RootChain *RootChainSession) ProcessExits() (*types.Transaction, error) {
	return _RootChain.Contract.ProcessExits(&_RootChain.TransactOpts)
}

// ProcessExits is a paid mutator transaction binding the contract method 0x58189e03.
//
// Solidity: function processExits() returns()
func (_RootChain *RootChainTransactorSession) ProcessExits() (*types.Transaction, error) {
	return _RootChain.Contract.ProcessExits(&_RootChain.TransactOpts)
}

// StartExit is a paid mutator transaction binding the contract method 0xf2f76549.
//
// Solidity: function startExit(_utxoBlockNumber uint256, _utxoTxIndex uint256, _utxoOutputIndex uint256, _encodedTx bytes, _txInclusionProof bytes, _txSignatures bytes, _txConfirmationSignatures bytes) returns()
func (_RootChain *RootChainTransactor) StartExit(opts *bind.TransactOpts, _utxoBlockNumber *big.Int, _utxoTxIndex *big.Int, _utxoOutputIndex *big.Int, _encodedTx []byte, _txInclusionProof []byte, _txSignatures []byte, _txConfirmationSignatures []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startExit", _utxoBlockNumber, _utxoTxIndex, _utxoOutputIndex, _encodedTx, _txInclusionProof, _txSignatures, _txConfirmationSignatures)
}

// StartExit is a paid mutator transaction binding the contract method 0xf2f76549.
//
// Solidity: function startExit(_utxoBlockNumber uint256, _utxoTxIndex uint256, _utxoOutputIndex uint256, _encodedTx bytes, _txInclusionProof bytes, _txSignatures bytes, _txConfirmationSignatures bytes) returns()
func (_RootChain *RootChainSession) StartExit(_utxoBlockNumber *big.Int, _utxoTxIndex *big.Int, _utxoOutputIndex *big.Int, _encodedTx []byte, _txInclusionProof []byte, _txSignatures []byte, _txConfirmationSignatures []byte) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, _utxoBlockNumber, _utxoTxIndex, _utxoOutputIndex, _encodedTx, _txInclusionProof, _txSignatures, _txConfirmationSignatures)
}

// StartExit is a paid mutator transaction binding the contract method 0xf2f76549.
//
// Solidity: function startExit(_utxoBlockNumber uint256, _utxoTxIndex uint256, _utxoOutputIndex uint256, _encodedTx bytes, _txInclusionProof bytes, _txSignatures bytes, _txConfirmationSignatures bytes) returns()
func (_RootChain *RootChainTransactorSession) StartExit(_utxoBlockNumber *big.Int, _utxoTxIndex *big.Int, _utxoOutputIndex *big.Int, _encodedTx []byte, _txInclusionProof []byte, _txSignatures []byte, _txConfirmationSignatures []byte) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, _utxoBlockNumber, _utxoTxIndex, _utxoOutputIndex, _encodedTx, _txInclusionProof, _txSignatures, _txConfirmationSignatures)
}

// RootChainDepositCreatedIterator is returned from FilterDepositCreated and is used to iterate over the raw logs and unpacked data for DepositCreated events raised by the RootChain contract.
type RootChainDepositCreatedIterator struct {
	Event *RootChainDepositCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RootChainDepositCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainDepositCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RootChainDepositCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RootChainDepositCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainDepositCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainDepositCreated represents a DepositCreated event raised by the RootChain contract.
type RootChainDepositCreated struct {
	Owner        common.Address
	Amount       *big.Int
	DepositBlock *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDepositCreated is a free log retrieval operation binding the contract event 0x72235260b946c311b2455bc3db7b029180adc32d36d49946ccb668abcd1a1a19.
//
// Solidity: event DepositCreated(owner indexed address, amount uint256, depositBlock uint256)
func (_RootChain *RootChainFilterer) FilterDepositCreated(opts *bind.FilterOpts, owner []common.Address) (*RootChainDepositCreatedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "DepositCreated", ownerRule)
	if err != nil {
		return nil, err
	}
	return &RootChainDepositCreatedIterator{contract: _RootChain.contract, event: "DepositCreated", logs: logs, sub: sub}, nil
}

// WatchDepositCreated is a free log subscription operation binding the contract event 0x72235260b946c311b2455bc3db7b029180adc32d36d49946ccb668abcd1a1a19.
//
// Solidity: event DepositCreated(owner indexed address, amount uint256, depositBlock uint256)
func (_RootChain *RootChainFilterer) WatchDepositCreated(opts *bind.WatchOpts, sink chan<- *RootChainDepositCreated, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "DepositCreated", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainDepositCreated)
				if err := _RootChain.contract.UnpackLog(event, "DepositCreated", log); err != nil {
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

// RootChainExitStartedIterator is returned from FilterExitStarted and is used to iterate over the raw logs and unpacked data for ExitStarted events raised by the RootChain contract.
type RootChainExitStartedIterator struct {
	Event *RootChainExitStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RootChainExitStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainExitStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RootChainExitStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RootChainExitStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainExitStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainExitStarted represents a ExitStarted event raised by the RootChain contract.
type RootChainExitStarted struct {
	Owner        common.Address
	UtxoPosition *big.Int
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterExitStarted is a free log retrieval operation binding the contract event 0xd62f48e64bfb0b95c15737946c829e7fe6378cd679a268c1ca5993dc235d900b.
//
// Solidity: event ExitStarted(owner indexed address, utxoPosition uint256, amount uint256)
func (_RootChain *RootChainFilterer) FilterExitStarted(opts *bind.FilterOpts, owner []common.Address) (*RootChainExitStartedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ExitStarted", ownerRule)
	if err != nil {
		return nil, err
	}
	return &RootChainExitStartedIterator{contract: _RootChain.contract, event: "ExitStarted", logs: logs, sub: sub}, nil
}

// WatchExitStarted is a free log subscription operation binding the contract event 0xd62f48e64bfb0b95c15737946c829e7fe6378cd679a268c1ca5993dc235d900b.
//
// Solidity: event ExitStarted(owner indexed address, utxoPosition uint256, amount uint256)
func (_RootChain *RootChainFilterer) WatchExitStarted(opts *bind.WatchOpts, sink chan<- *RootChainExitStarted, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ExitStarted", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainExitStarted)
				if err := _RootChain.contract.UnpackLog(event, "ExitStarted", log); err != nil {
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

// RootChainPlasmaBlockRootCommittedIterator is returned from FilterPlasmaBlockRootCommitted and is used to iterate over the raw logs and unpacked data for PlasmaBlockRootCommitted events raised by the RootChain contract.
type RootChainPlasmaBlockRootCommittedIterator struct {
	Event *RootChainPlasmaBlockRootCommitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RootChainPlasmaBlockRootCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainPlasmaBlockRootCommitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RootChainPlasmaBlockRootCommitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RootChainPlasmaBlockRootCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainPlasmaBlockRootCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainPlasmaBlockRootCommitted represents a PlasmaBlockRootCommitted event raised by the RootChain contract.
type RootChainPlasmaBlockRootCommitted struct {
	BlockNumber *big.Int
	Root        [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPlasmaBlockRootCommitted is a free log retrieval operation binding the contract event 0xa787dc997c5916e41f62289b9230d631cdaa09737895a5bbe1ddffe1fd596567.
//
// Solidity: event PlasmaBlockRootCommitted(blockNumber uint256, root bytes32)
func (_RootChain *RootChainFilterer) FilterPlasmaBlockRootCommitted(opts *bind.FilterOpts) (*RootChainPlasmaBlockRootCommittedIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "PlasmaBlockRootCommitted")
	if err != nil {
		return nil, err
	}
	return &RootChainPlasmaBlockRootCommittedIterator{contract: _RootChain.contract, event: "PlasmaBlockRootCommitted", logs: logs, sub: sub}, nil
}

// WatchPlasmaBlockRootCommitted is a free log subscription operation binding the contract event 0xa787dc997c5916e41f62289b9230d631cdaa09737895a5bbe1ddffe1fd596567.
//
// Solidity: event PlasmaBlockRootCommitted(blockNumber uint256, root bytes32)
func (_RootChain *RootChainFilterer) WatchPlasmaBlockRootCommitted(opts *bind.WatchOpts, sink chan<- *RootChainPlasmaBlockRootCommitted) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "PlasmaBlockRootCommitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainPlasmaBlockRootCommitted)
				if err := _RootChain.contract.UnpackLog(event, "PlasmaBlockRootCommitted", log); err != nil {
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
