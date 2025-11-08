// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tochka_escrow_payment

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
)

// TochkaEscrowPaymentMetaData contains all meta data concerning the TochkaEscrowPayment contract.
var TochkaEscrowPaymentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_contributors\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_percents\",\"type\":\"uint256[]\"}],\"name\":\"multitransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// TochkaEscrowPaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use TochkaEscrowPaymentMetaData.ABI instead.
var TochkaEscrowPaymentABI = TochkaEscrowPaymentMetaData.ABI

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
	parsed, err := abi.JSON(strings.NewReader(TochkaEscrowPaymentABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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
