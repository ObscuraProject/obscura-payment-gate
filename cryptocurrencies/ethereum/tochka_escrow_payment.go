// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TochkaEscrowPaymentMetaData contains all meta data concerning the TochkaEscrowPayment contract.
var TochkaEscrowPaymentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_contributors\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_percents\",\"type\":\"uint256[]\"}],\"name\":\"multitransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104f9806100206000396000f3fe60806040526004361061001e5760003560e01c80633edfd1e114610023575b600080fd5b61003d6004803603810190610038919061024d565b61003f565b005b60003490506000811161005157600080fd5b8484905083839050111561006457600080fd5b6000805b8686905082101561017f576064858584818110610088576100876102ce565b5b90506020020135111561009a57600080fd5b60008585848181106100af576100ae6102ce565b5b9050602002013510156100c157600080fd5b60648585848181106100d6576100d56102ce565b5b90506020020135846100e89190610336565b6100f291906103bf565b9050600081111561016c578686838181106101105761010f6102ce565b5b9050602002016020810190610125919061044e565b73ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f1935050505015801561016a573d6000803e3d6000fd5b505b81806101779061047b565b925050610068565b50505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126101b7576101b6610192565b5b8235905067ffffffffffffffff8111156101d4576101d3610197565b5b6020830191508360208202830111156101f0576101ef61019c565b5b9250929050565b60008083601f84011261020d5761020c610192565b5b8235905067ffffffffffffffff81111561022a57610229610197565b5b6020830191508360208202830111156102465761024561019c565b5b9250929050565b6000806000806040858703121561026757610266610188565b5b600085013567ffffffffffffffff8111156102855761028461018d565b5b610291878288016101a1565b9450945050602085013567ffffffffffffffff8111156102b4576102b361018d565b5b6102c0878288016101f7565b925092505092959194509250565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610341826102fd565b915061034c836102fd565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561038557610384610307565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006103ca826102fd565b91506103d5836102fd565b9250826103e5576103e4610390565b5b828204905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061041b826103f0565b9050919050565b61042b81610410565b811461043657600080fd5b50565b60008135905061044881610422565b92915050565b60006020828403121561046457610463610188565b5b600061047284828501610439565b91505092915050565b6000610486826102fd565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036104b8576104b7610307565b5b60018201905091905056fea2646970667358221220c6fc7b3906ea313ef763156ab583d2206259ef38f4aa5f02e831ab742c31e14964736f6c634300080d0033",
}

// TochkaEscrowPaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use TochkaEscrowPaymentMetaData.ABI instead.
var TochkaEscrowPaymentABI = TochkaEscrowPaymentMetaData.ABI

// TochkaEscrowPaymentBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TochkaEscrowPaymentMetaData.Bin instead.
var TochkaEscrowPaymentBin = TochkaEscrowPaymentMetaData.Bin

// DeployTochkaEscrowPayment deploys a new Ethereum contract, binding an instance of TochkaEscrowPayment to it.
func DeployTochkaEscrowPayment(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TochkaEscrowPayment, error) {
	parsed, err := TochkaEscrowPaymentMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TochkaEscrowPaymentBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TochkaEscrowPayment{TochkaEscrowPaymentCaller: TochkaEscrowPaymentCaller{contract: contract}, TochkaEscrowPaymentTransactor: TochkaEscrowPaymentTransactor{contract: contract}, TochkaEscrowPaymentFilterer: TochkaEscrowPaymentFilterer{contract: contract}}, nil
}

// TochkaEscrowPayment is an auto generated Go binding around an Ethereum contract.
type TochkaEscrowPayment struct {
	TochkaEscrowPaymentCaller     // Read-only binding to the contract
	TochkaEscrowPaymentTransactor // Write-only binding to the contract
	TochkaEscrowPaymentFilterer   // Log filterer for contract events
}

// TochkaEscrowPaymentCaller is an auto generated read-only Go binding around an Ethereum contract.
type TochkaEscrowPaymentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TochkaEscrowPaymentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TochkaEscrowPaymentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TochkaEscrowPaymentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TochkaEscrowPaymentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TochkaEscrowPaymentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TochkaEscrowPaymentSession struct {
	Contract     *TochkaEscrowPayment // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TochkaEscrowPaymentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TochkaEscrowPaymentCallerSession struct {
	Contract *TochkaEscrowPaymentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TochkaEscrowPaymentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TochkaEscrowPaymentTransactorSession struct {
	Contract     *TochkaEscrowPaymentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TochkaEscrowPaymentRaw is an auto generated low-level Go binding around an Ethereum contract.
type TochkaEscrowPaymentRaw struct {
	Contract *TochkaEscrowPayment // Generic contract binding to access the raw methods on
}

// TochkaEscrowPaymentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TochkaEscrowPaymentCallerRaw struct {
	Contract *TochkaEscrowPaymentCaller // Generic read-only contract binding to access the raw methods on
}

// TochkaEscrowPaymentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TochkaEscrowPaymentTransactorRaw struct {
	Contract *TochkaEscrowPaymentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTochkaEscrowPayment creates a new instance of TochkaEscrowPayment, bound to a specific deployed contract.
func NewTochkaEscrowPayment(address common.Address, backend bind.ContractBackend) (*TochkaEscrowPayment, error) {
	contract, err := bindTochkaEscrowPayment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TochkaEscrowPayment{TochkaEscrowPaymentCaller: TochkaEscrowPaymentCaller{contract: contract}, TochkaEscrowPaymentTransactor: TochkaEscrowPaymentTransactor{contract: contract}, TochkaEscrowPaymentFilterer: TochkaEscrowPaymentFilterer{contract: contract}}, nil
}

// NewTochkaEscrowPaymentCaller creates a new read-only instance of TochkaEscrowPayment, bound to a specific deployed contract.
func NewTochkaEscrowPaymentCaller(address common.Address, caller bind.ContractCaller) (*TochkaEscrowPaymentCaller, error) {
	contract, err := bindTochkaEscrowPayment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TochkaEscrowPaymentCaller{contract: contract}, nil
}

// NewTochkaEscrowPaymentTransactor creates a new write-only instance of TochkaEscrowPayment, bound to a specific deployed contract.
func NewTochkaEscrowPaymentTransactor(address common.Address, transactor bind.ContractTransactor) (*TochkaEscrowPaymentTransactor, error) {
	contract, err := bindTochkaEscrowPayment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TochkaEscrowPaymentTransactor{contract: contract}, nil
}

// NewTochkaEscrowPaymentFilterer creates a new log filterer instance of TochkaEscrowPayment, bound to a specific deployed contract.
func NewTochkaEscrowPaymentFilterer(address common.Address, filterer bind.ContractFilterer) (*TochkaEscrowPaymentFilterer, error) {
	contract, err := bindTochkaEscrowPayment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TochkaEscrowPaymentFilterer{contract: contract}, nil
}

// bindTochkaEscrowPayment binds a generic wrapper to an already deployed contract.
func bindTochkaEscrowPayment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TochkaEscrowPaymentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TochkaEscrowPayment *TochkaEscrowPaymentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TochkaEscrowPayment.Contract.TochkaEscrowPaymentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TochkaEscrowPayment *TochkaEscrowPaymentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.TochkaEscrowPaymentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TochkaEscrowPayment *TochkaEscrowPaymentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.TochkaEscrowPaymentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TochkaEscrowPayment *TochkaEscrowPaymentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TochkaEscrowPayment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TochkaEscrowPayment *TochkaEscrowPaymentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TochkaEscrowPayment *TochkaEscrowPaymentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.contract.Transact(opts, method, params...)
}

// Multitransfer is a paid mutator transaction binding the contract method 0x3edfd1e1.
//
// Solidity: function multitransfer(address[] _contributors, uint256[] _percents) payable returns()
func (_TochkaEscrowPayment *TochkaEscrowPaymentTransactor) Multitransfer(opts *bind.TransactOpts, _contributors []common.Address, _percents []*big.Int) (*types.Transaction, error) {
	return _TochkaEscrowPayment.contract.Transact(opts, "multitransfer", _contributors, _percents)
}

// Multitransfer is a paid mutator transaction binding the contract method 0x3edfd1e1.
//
// Solidity: function multitransfer(address[] _contributors, uint256[] _percents) payable returns()
func (_TochkaEscrowPayment *TochkaEscrowPaymentSession) Multitransfer(_contributors []common.Address, _percents []*big.Int) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.Multitransfer(&_TochkaEscrowPayment.TransactOpts, _contributors, _percents)
}

// Multitransfer is a paid mutator transaction binding the contract method 0x3edfd1e1.
//
// Solidity: function multitransfer(address[] _contributors, uint256[] _percents) payable returns()
func (_TochkaEscrowPayment *TochkaEscrowPaymentTransactorSession) Multitransfer(_contributors []common.Address, _percents []*big.Int) (*types.Transaction, error) {
	return _TochkaEscrowPayment.Contract.Multitransfer(&_TochkaEscrowPayment.TransactOpts, _contributors, _percents)
}
