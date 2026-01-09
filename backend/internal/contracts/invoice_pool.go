// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// InvoicePoolInvestment is an auto generated low-level Go binding around an user-defined struct.
type InvoicePoolInvestment struct {
	Investor       common.Address
	Amount         *big.Int
	ExpectedReturn *big.Int
	ActualReturn   *big.Int
	Claimed        bool
	InvestedAt     *big.Int
}

// InvoicePoolPool is an auto generated low-level Go binding around an user-defined struct.
type InvoicePoolPool struct {
	TokenId       *big.Int
	TargetAmount  *big.Int
	FundedAmount  *big.Int
	InvestorCount *big.Int
	InterestRate  *big.Int
	DueDate       *big.Int
	Exporter      common.Address
	Status        uint8
	OpenedAt      *big.Int
	FilledAt      *big.Int
	DisbursedAt   *big.Int
	ClosedAt      *big.Int
}

// InvoicePoolMetaData contains all meta data concerning the InvoicePool contract.
var InvoicePoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_invoiceNFT\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_platformWallet\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DisbursementRecorded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ExcessRepaymentRecorded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expectedReturn\",\"type\":\"uint256\"}],\"name\":\"InvestmentRecorded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InvestorReturnRecorded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"PoolClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"targetAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"PoolDefaulted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"investorCount\",\"type\":\"uint256\"}],\"name\":\"PoolFilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RepaymentRecorded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPERATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"createPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"getInvestorPools\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getPool\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fundedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"investorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"internalType\":\"enumInvoicePool.PoolStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"openedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"filledAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"disbursedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"closedAt\",\"type\":\"uint256\"}],\"internalType\":\"structInvoicePool.Pool\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getPoolInvestments\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedReturn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualReturn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"investedAt\",\"type\":\"uint256\"}],\"internalType\":\"structInvoicePool.Investment[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getRemainingCapacity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"investorPools\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invoiceNFT\",\"outputs\":[{\"internalType\":\"contractInvoiceNFT\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"markDefaulted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformFeeBps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"poolInvestments\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedReturn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualReturn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"investedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pools\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fundedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"investorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"internalType\":\"enumInvoicePool.PoolStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"openedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"filledAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"disbursedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"closedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"recordDisbursement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"recordExcessRepayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"recordInvestment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"investorReturns\",\"type\":\"uint256[]\"}],\"name\":\"recordRepayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFeeBps\",\"type\":\"uint256\"}],\"name\":\"setPlatformFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newWallet\",\"type\":\"address\"}],\"name\":\"setPlatformWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// InvoicePoolABI is the input ABI used to generate the binding from.
// Deprecated: Use InvoicePoolMetaData.ABI instead.
var InvoicePoolABI = InvoicePoolMetaData.ABI

// InvoicePool is an auto generated Go binding around an Ethereum contract.
type InvoicePool struct {
	InvoicePoolCaller     // Read-only binding to the contract
	InvoicePoolTransactor // Write-only binding to the contract
	InvoicePoolFilterer   // Log filterer for contract events
}

// InvoicePoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type InvoicePoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoicePoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InvoicePoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoicePoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InvoicePoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoicePoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InvoicePoolSession struct {
	Contract     *InvoicePool      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InvoicePoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InvoicePoolCallerSession struct {
	Contract *InvoicePoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// InvoicePoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InvoicePoolTransactorSession struct {
	Contract     *InvoicePoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// InvoicePoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type InvoicePoolRaw struct {
	Contract *InvoicePool // Generic contract binding to access the raw methods on
}

// InvoicePoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InvoicePoolCallerRaw struct {
	Contract *InvoicePoolCaller // Generic read-only contract binding to access the raw methods on
}

// InvoicePoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InvoicePoolTransactorRaw struct {
	Contract *InvoicePoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInvoicePool creates a new instance of InvoicePool, bound to a specific deployed contract.
func NewInvoicePool(address common.Address, backend bind.ContractBackend) (*InvoicePool, error) {
	contract, err := bindInvoicePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InvoicePool{InvoicePoolCaller: InvoicePoolCaller{contract: contract}, InvoicePoolTransactor: InvoicePoolTransactor{contract: contract}, InvoicePoolFilterer: InvoicePoolFilterer{contract: contract}}, nil
}

// NewInvoicePoolCaller creates a new read-only instance of InvoicePool, bound to a specific deployed contract.
func NewInvoicePoolCaller(address common.Address, caller bind.ContractCaller) (*InvoicePoolCaller, error) {
	contract, err := bindInvoicePool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolCaller{contract: contract}, nil
}

// NewInvoicePoolTransactor creates a new write-only instance of InvoicePool, bound to a specific deployed contract.
func NewInvoicePoolTransactor(address common.Address, transactor bind.ContractTransactor) (*InvoicePoolTransactor, error) {
	contract, err := bindInvoicePool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolTransactor{contract: contract}, nil
}

// NewInvoicePoolFilterer creates a new log filterer instance of InvoicePool, bound to a specific deployed contract.
func NewInvoicePoolFilterer(address common.Address, filterer bind.ContractFilterer) (*InvoicePoolFilterer, error) {
	contract, err := bindInvoicePool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolFilterer{contract: contract}, nil
}

// bindInvoicePool binds a generic wrapper to an already deployed contract.
func bindInvoicePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InvoicePoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InvoicePool *InvoicePoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InvoicePool.Contract.InvoicePoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InvoicePool *InvoicePoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoicePool.Contract.InvoicePoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InvoicePool *InvoicePoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InvoicePool.Contract.InvoicePoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InvoicePool *InvoicePoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InvoicePool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InvoicePool *InvoicePoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoicePool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InvoicePool *InvoicePoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InvoicePool.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _InvoicePool.Contract.DEFAULTADMINROLE(&_InvoicePool.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _InvoicePool.Contract.DEFAULTADMINROLE(&_InvoicePool.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolCaller) OPERATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "OPERATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolSession) OPERATORROLE() ([32]byte, error) {
	return _InvoicePool.Contract.OPERATORROLE(&_InvoicePool.CallOpts)
}

// OPERATORROLE is a free data retrieval call binding the contract method 0xf5b541a6.
//
// Solidity: function OPERATOR_ROLE() view returns(bytes32)
func (_InvoicePool *InvoicePoolCallerSession) OPERATORROLE() ([32]byte, error) {
	return _InvoicePool.Contract.OPERATORROLE(&_InvoicePool.CallOpts)
}

// GetInvestorPools is a free data retrieval call binding the contract method 0x32ec5788.
//
// Solidity: function getInvestorPools(address investor) view returns(uint256[])
func (_InvoicePool *InvoicePoolCaller) GetInvestorPools(opts *bind.CallOpts, investor common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "getInvestorPools", investor)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetInvestorPools is a free data retrieval call binding the contract method 0x32ec5788.
//
// Solidity: function getInvestorPools(address investor) view returns(uint256[])
func (_InvoicePool *InvoicePoolSession) GetInvestorPools(investor common.Address) ([]*big.Int, error) {
	return _InvoicePool.Contract.GetInvestorPools(&_InvoicePool.CallOpts, investor)
}

// GetInvestorPools is a free data retrieval call binding the contract method 0x32ec5788.
//
// Solidity: function getInvestorPools(address investor) view returns(uint256[])
func (_InvoicePool *InvoicePoolCallerSession) GetInvestorPools(investor common.Address) ([]*big.Int, error) {
	return _InvoicePool.Contract.GetInvestorPools(&_InvoicePool.CallOpts, investor)
}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 tokenId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,address,uint8,uint256,uint256,uint256,uint256))
func (_InvoicePool *InvoicePoolCaller) GetPool(opts *bind.CallOpts, tokenId *big.Int) (InvoicePoolPool, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "getPool", tokenId)

	if err != nil {
		return *new(InvoicePoolPool), err
	}

	out0 := *abi.ConvertType(out[0], new(InvoicePoolPool)).(*InvoicePoolPool)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 tokenId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,address,uint8,uint256,uint256,uint256,uint256))
func (_InvoicePool *InvoicePoolSession) GetPool(tokenId *big.Int) (InvoicePoolPool, error) {
	return _InvoicePool.Contract.GetPool(&_InvoicePool.CallOpts, tokenId)
}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 tokenId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,address,uint8,uint256,uint256,uint256,uint256))
func (_InvoicePool *InvoicePoolCallerSession) GetPool(tokenId *big.Int) (InvoicePoolPool, error) {
	return _InvoicePool.Contract.GetPool(&_InvoicePool.CallOpts, tokenId)
}

// GetPoolInvestments is a free data retrieval call binding the contract method 0xab56d0ce.
//
// Solidity: function getPoolInvestments(uint256 tokenId) view returns((address,uint256,uint256,uint256,bool,uint256)[])
func (_InvoicePool *InvoicePoolCaller) GetPoolInvestments(opts *bind.CallOpts, tokenId *big.Int) ([]InvoicePoolInvestment, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "getPoolInvestments", tokenId)

	if err != nil {
		return *new([]InvoicePoolInvestment), err
	}

	out0 := *abi.ConvertType(out[0], new([]InvoicePoolInvestment)).(*[]InvoicePoolInvestment)

	return out0, err

}

// GetPoolInvestments is a free data retrieval call binding the contract method 0xab56d0ce.
//
// Solidity: function getPoolInvestments(uint256 tokenId) view returns((address,uint256,uint256,uint256,bool,uint256)[])
func (_InvoicePool *InvoicePoolSession) GetPoolInvestments(tokenId *big.Int) ([]InvoicePoolInvestment, error) {
	return _InvoicePool.Contract.GetPoolInvestments(&_InvoicePool.CallOpts, tokenId)
}

// GetPoolInvestments is a free data retrieval call binding the contract method 0xab56d0ce.
//
// Solidity: function getPoolInvestments(uint256 tokenId) view returns((address,uint256,uint256,uint256,bool,uint256)[])
func (_InvoicePool *InvoicePoolCallerSession) GetPoolInvestments(tokenId *big.Int) ([]InvoicePoolInvestment, error) {
	return _InvoicePool.Contract.GetPoolInvestments(&_InvoicePool.CallOpts, tokenId)
}

// GetRemainingCapacity is a free data retrieval call binding the contract method 0x1b8ef0bb.
//
// Solidity: function getRemainingCapacity(uint256 tokenId) view returns(uint256)
func (_InvoicePool *InvoicePoolCaller) GetRemainingCapacity(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "getRemainingCapacity", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingCapacity is a free data retrieval call binding the contract method 0x1b8ef0bb.
//
// Solidity: function getRemainingCapacity(uint256 tokenId) view returns(uint256)
func (_InvoicePool *InvoicePoolSession) GetRemainingCapacity(tokenId *big.Int) (*big.Int, error) {
	return _InvoicePool.Contract.GetRemainingCapacity(&_InvoicePool.CallOpts, tokenId)
}

// GetRemainingCapacity is a free data retrieval call binding the contract method 0x1b8ef0bb.
//
// Solidity: function getRemainingCapacity(uint256 tokenId) view returns(uint256)
func (_InvoicePool *InvoicePoolCallerSession) GetRemainingCapacity(tokenId *big.Int) (*big.Int, error) {
	return _InvoicePool.Contract.GetRemainingCapacity(&_InvoicePool.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoicePool *InvoicePoolCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoicePool *InvoicePoolSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _InvoicePool.Contract.GetRoleAdmin(&_InvoicePool.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoicePool *InvoicePoolCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _InvoicePool.Contract.GetRoleAdmin(&_InvoicePool.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoicePool *InvoicePoolCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoicePool *InvoicePoolSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _InvoicePool.Contract.HasRole(&_InvoicePool.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoicePool *InvoicePoolCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _InvoicePool.Contract.HasRole(&_InvoicePool.CallOpts, role, account)
}

// InvestorPools is a free data retrieval call binding the contract method 0xfdf9bc1c.
//
// Solidity: function investorPools(address , uint256 ) view returns(uint256)
func (_InvoicePool *InvoicePoolCaller) InvestorPools(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "investorPools", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InvestorPools is a free data retrieval call binding the contract method 0xfdf9bc1c.
//
// Solidity: function investorPools(address , uint256 ) view returns(uint256)
func (_InvoicePool *InvoicePoolSession) InvestorPools(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _InvoicePool.Contract.InvestorPools(&_InvoicePool.CallOpts, arg0, arg1)
}

// InvestorPools is a free data retrieval call binding the contract method 0xfdf9bc1c.
//
// Solidity: function investorPools(address , uint256 ) view returns(uint256)
func (_InvoicePool *InvoicePoolCallerSession) InvestorPools(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _InvoicePool.Contract.InvestorPools(&_InvoicePool.CallOpts, arg0, arg1)
}

// InvoiceNFT is a free data retrieval call binding the contract method 0x4bb6f06d.
//
// Solidity: function invoiceNFT() view returns(address)
func (_InvoicePool *InvoicePoolCaller) InvoiceNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "invoiceNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// InvoiceNFT is a free data retrieval call binding the contract method 0x4bb6f06d.
//
// Solidity: function invoiceNFT() view returns(address)
func (_InvoicePool *InvoicePoolSession) InvoiceNFT() (common.Address, error) {
	return _InvoicePool.Contract.InvoiceNFT(&_InvoicePool.CallOpts)
}

// InvoiceNFT is a free data retrieval call binding the contract method 0x4bb6f06d.
//
// Solidity: function invoiceNFT() view returns(address)
func (_InvoicePool *InvoicePoolCallerSession) InvoiceNFT() (common.Address, error) {
	return _InvoicePool.Contract.InvoiceNFT(&_InvoicePool.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoicePool *InvoicePoolCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoicePool *InvoicePoolSession) Paused() (bool, error) {
	return _InvoicePool.Contract.Paused(&_InvoicePool.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoicePool *InvoicePoolCallerSession) Paused() (bool, error) {
	return _InvoicePool.Contract.Paused(&_InvoicePool.CallOpts)
}

// PlatformFeeBps is a free data retrieval call binding the contract method 0x22dcd13e.
//
// Solidity: function platformFeeBps() view returns(uint256)
func (_InvoicePool *InvoicePoolCaller) PlatformFeeBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "platformFeeBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PlatformFeeBps is a free data retrieval call binding the contract method 0x22dcd13e.
//
// Solidity: function platformFeeBps() view returns(uint256)
func (_InvoicePool *InvoicePoolSession) PlatformFeeBps() (*big.Int, error) {
	return _InvoicePool.Contract.PlatformFeeBps(&_InvoicePool.CallOpts)
}

// PlatformFeeBps is a free data retrieval call binding the contract method 0x22dcd13e.
//
// Solidity: function platformFeeBps() view returns(uint256)
func (_InvoicePool *InvoicePoolCallerSession) PlatformFeeBps() (*big.Int, error) {
	return _InvoicePool.Contract.PlatformFeeBps(&_InvoicePool.CallOpts)
}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_InvoicePool *InvoicePoolCaller) PlatformWallet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "platformWallet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_InvoicePool *InvoicePoolSession) PlatformWallet() (common.Address, error) {
	return _InvoicePool.Contract.PlatformWallet(&_InvoicePool.CallOpts)
}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_InvoicePool *InvoicePoolCallerSession) PlatformWallet() (common.Address, error) {
	return _InvoicePool.Contract.PlatformWallet(&_InvoicePool.CallOpts)
}

// PoolInvestments is a free data retrieval call binding the contract method 0x84dab984.
//
// Solidity: function poolInvestments(uint256 , uint256 ) view returns(address investor, uint256 amount, uint256 expectedReturn, uint256 actualReturn, bool claimed, uint256 investedAt)
func (_InvoicePool *InvoicePoolCaller) PoolInvestments(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	Investor       common.Address
	Amount         *big.Int
	ExpectedReturn *big.Int
	ActualReturn   *big.Int
	Claimed        bool
	InvestedAt     *big.Int
}, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "poolInvestments", arg0, arg1)

	outstruct := new(struct {
		Investor       common.Address
		Amount         *big.Int
		ExpectedReturn *big.Int
		ActualReturn   *big.Int
		Claimed        bool
		InvestedAt     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Investor = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ExpectedReturn = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.ActualReturn = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Claimed = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.InvestedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PoolInvestments is a free data retrieval call binding the contract method 0x84dab984.
//
// Solidity: function poolInvestments(uint256 , uint256 ) view returns(address investor, uint256 amount, uint256 expectedReturn, uint256 actualReturn, bool claimed, uint256 investedAt)
func (_InvoicePool *InvoicePoolSession) PoolInvestments(arg0 *big.Int, arg1 *big.Int) (struct {
	Investor       common.Address
	Amount         *big.Int
	ExpectedReturn *big.Int
	ActualReturn   *big.Int
	Claimed        bool
	InvestedAt     *big.Int
}, error) {
	return _InvoicePool.Contract.PoolInvestments(&_InvoicePool.CallOpts, arg0, arg1)
}

// PoolInvestments is a free data retrieval call binding the contract method 0x84dab984.
//
// Solidity: function poolInvestments(uint256 , uint256 ) view returns(address investor, uint256 amount, uint256 expectedReturn, uint256 actualReturn, bool claimed, uint256 investedAt)
func (_InvoicePool *InvoicePoolCallerSession) PoolInvestments(arg0 *big.Int, arg1 *big.Int) (struct {
	Investor       common.Address
	Amount         *big.Int
	ExpectedReturn *big.Int
	ActualReturn   *big.Int
	Claimed        bool
	InvestedAt     *big.Int
}, error) {
	return _InvoicePool.Contract.PoolInvestments(&_InvoicePool.CallOpts, arg0, arg1)
}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(uint256 tokenId, uint256 targetAmount, uint256 fundedAmount, uint256 investorCount, uint256 interestRate, uint256 dueDate, address exporter, uint8 status, uint256 openedAt, uint256 filledAt, uint256 disbursedAt, uint256 closedAt)
func (_InvoicePool *InvoicePoolCaller) Pools(opts *bind.CallOpts, arg0 *big.Int) (struct {
	TokenId       *big.Int
	TargetAmount  *big.Int
	FundedAmount  *big.Int
	InvestorCount *big.Int
	InterestRate  *big.Int
	DueDate       *big.Int
	Exporter      common.Address
	Status        uint8
	OpenedAt      *big.Int
	FilledAt      *big.Int
	DisbursedAt   *big.Int
	ClosedAt      *big.Int
}, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "pools", arg0)

	outstruct := new(struct {
		TokenId       *big.Int
		TargetAmount  *big.Int
		FundedAmount  *big.Int
		InvestorCount *big.Int
		InterestRate  *big.Int
		DueDate       *big.Int
		Exporter      common.Address
		Status        uint8
		OpenedAt      *big.Int
		FilledAt      *big.Int
		DisbursedAt   *big.Int
		ClosedAt      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TargetAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FundedAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.InvestorCount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.InterestRate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.DueDate = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Exporter = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.Status = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.OpenedAt = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.FilledAt = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.DisbursedAt = *abi.ConvertType(out[10], new(*big.Int)).(**big.Int)
	outstruct.ClosedAt = *abi.ConvertType(out[11], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(uint256 tokenId, uint256 targetAmount, uint256 fundedAmount, uint256 investorCount, uint256 interestRate, uint256 dueDate, address exporter, uint8 status, uint256 openedAt, uint256 filledAt, uint256 disbursedAt, uint256 closedAt)
func (_InvoicePool *InvoicePoolSession) Pools(arg0 *big.Int) (struct {
	TokenId       *big.Int
	TargetAmount  *big.Int
	FundedAmount  *big.Int
	InvestorCount *big.Int
	InterestRate  *big.Int
	DueDate       *big.Int
	Exporter      common.Address
	Status        uint8
	OpenedAt      *big.Int
	FilledAt      *big.Int
	DisbursedAt   *big.Int
	ClosedAt      *big.Int
}, error) {
	return _InvoicePool.Contract.Pools(&_InvoicePool.CallOpts, arg0)
}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(uint256 tokenId, uint256 targetAmount, uint256 fundedAmount, uint256 investorCount, uint256 interestRate, uint256 dueDate, address exporter, uint8 status, uint256 openedAt, uint256 filledAt, uint256 disbursedAt, uint256 closedAt)
func (_InvoicePool *InvoicePoolCallerSession) Pools(arg0 *big.Int) (struct {
	TokenId       *big.Int
	TargetAmount  *big.Int
	FundedAmount  *big.Int
	InvestorCount *big.Int
	InterestRate  *big.Int
	DueDate       *big.Int
	Exporter      common.Address
	Status        uint8
	OpenedAt      *big.Int
	FilledAt      *big.Int
	DisbursedAt   *big.Int
	ClosedAt      *big.Int
}, error) {
	return _InvoicePool.Contract.Pools(&_InvoicePool.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoicePool *InvoicePoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _InvoicePool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoicePool *InvoicePoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _InvoicePool.Contract.SupportsInterface(&_InvoicePool.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoicePool *InvoicePoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _InvoicePool.Contract.SupportsInterface(&_InvoicePool.CallOpts, interfaceId)
}

// CreatePool is a paid mutator transaction binding the contract method 0x8259e6a0.
//
// Solidity: function createPool(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactor) CreatePool(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "createPool", tokenId)
}

// CreatePool is a paid mutator transaction binding the contract method 0x8259e6a0.
//
// Solidity: function createPool(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolSession) CreatePool(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.CreatePool(&_InvoicePool.TransactOpts, tokenId)
}

// CreatePool is a paid mutator transaction binding the contract method 0x8259e6a0.
//
// Solidity: function createPool(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactorSession) CreatePool(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.CreatePool(&_InvoicePool.TransactOpts, tokenId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.GrantRole(&_InvoicePool.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.GrantRole(&_InvoicePool.TransactOpts, role, account)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactor) MarkDefaulted(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "markDefaulted", tokenId)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolSession) MarkDefaulted(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.MarkDefaulted(&_InvoicePool.TransactOpts, tokenId)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactorSession) MarkDefaulted(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.MarkDefaulted(&_InvoicePool.TransactOpts, tokenId)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoicePool *InvoicePoolTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoicePool *InvoicePoolSession) Pause() (*types.Transaction, error) {
	return _InvoicePool.Contract.Pause(&_InvoicePool.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoicePool *InvoicePoolTransactorSession) Pause() (*types.Transaction, error) {
	return _InvoicePool.Contract.Pause(&_InvoicePool.TransactOpts)
}

// RecordDisbursement is a paid mutator transaction binding the contract method 0x70e44b74.
//
// Solidity: function recordDisbursement(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactor) RecordDisbursement(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "recordDisbursement", tokenId)
}

// RecordDisbursement is a paid mutator transaction binding the contract method 0x70e44b74.
//
// Solidity: function recordDisbursement(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolSession) RecordDisbursement(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordDisbursement(&_InvoicePool.TransactOpts, tokenId)
}

// RecordDisbursement is a paid mutator transaction binding the contract method 0x70e44b74.
//
// Solidity: function recordDisbursement(uint256 tokenId) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RecordDisbursement(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordDisbursement(&_InvoicePool.TransactOpts, tokenId)
}

// RecordExcessRepayment is a paid mutator transaction binding the contract method 0xffdea07d.
//
// Solidity: function recordExcessRepayment(uint256 tokenId, address recipient, uint256 amount) returns()
func (_InvoicePool *InvoicePoolTransactor) RecordExcessRepayment(opts *bind.TransactOpts, tokenId *big.Int, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "recordExcessRepayment", tokenId, recipient, amount)
}

// RecordExcessRepayment is a paid mutator transaction binding the contract method 0xffdea07d.
//
// Solidity: function recordExcessRepayment(uint256 tokenId, address recipient, uint256 amount) returns()
func (_InvoicePool *InvoicePoolSession) RecordExcessRepayment(tokenId *big.Int, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordExcessRepayment(&_InvoicePool.TransactOpts, tokenId, recipient, amount)
}

// RecordExcessRepayment is a paid mutator transaction binding the contract method 0xffdea07d.
//
// Solidity: function recordExcessRepayment(uint256 tokenId, address recipient, uint256 amount) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RecordExcessRepayment(tokenId *big.Int, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordExcessRepayment(&_InvoicePool.TransactOpts, tokenId, recipient, amount)
}

// RecordInvestment is a paid mutator transaction binding the contract method 0x336faacf.
//
// Solidity: function recordInvestment(uint256 tokenId, address investor, uint256 amount) returns()
func (_InvoicePool *InvoicePoolTransactor) RecordInvestment(opts *bind.TransactOpts, tokenId *big.Int, investor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "recordInvestment", tokenId, investor, amount)
}

// RecordInvestment is a paid mutator transaction binding the contract method 0x336faacf.
//
// Solidity: function recordInvestment(uint256 tokenId, address investor, uint256 amount) returns()
func (_InvoicePool *InvoicePoolSession) RecordInvestment(tokenId *big.Int, investor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordInvestment(&_InvoicePool.TransactOpts, tokenId, investor, amount)
}

// RecordInvestment is a paid mutator transaction binding the contract method 0x336faacf.
//
// Solidity: function recordInvestment(uint256 tokenId, address investor, uint256 amount) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RecordInvestment(tokenId *big.Int, investor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordInvestment(&_InvoicePool.TransactOpts, tokenId, investor, amount)
}

// RecordRepayment is a paid mutator transaction binding the contract method 0x5c6e786d.
//
// Solidity: function recordRepayment(uint256 tokenId, uint256 totalAmount, uint256[] investorReturns) returns()
func (_InvoicePool *InvoicePoolTransactor) RecordRepayment(opts *bind.TransactOpts, tokenId *big.Int, totalAmount *big.Int, investorReturns []*big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "recordRepayment", tokenId, totalAmount, investorReturns)
}

// RecordRepayment is a paid mutator transaction binding the contract method 0x5c6e786d.
//
// Solidity: function recordRepayment(uint256 tokenId, uint256 totalAmount, uint256[] investorReturns) returns()
func (_InvoicePool *InvoicePoolSession) RecordRepayment(tokenId *big.Int, totalAmount *big.Int, investorReturns []*big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordRepayment(&_InvoicePool.TransactOpts, tokenId, totalAmount, investorReturns)
}

// RecordRepayment is a paid mutator transaction binding the contract method 0x5c6e786d.
//
// Solidity: function recordRepayment(uint256 tokenId, uint256 totalAmount, uint256[] investorReturns) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RecordRepayment(tokenId *big.Int, totalAmount *big.Int, investorReturns []*big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.RecordRepayment(&_InvoicePool.TransactOpts, tokenId, totalAmount, investorReturns)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoicePool *InvoicePoolTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoicePool *InvoicePoolSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.RenounceRole(&_InvoicePool.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.RenounceRole(&_InvoicePool.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.RevokeRole(&_InvoicePool.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoicePool *InvoicePoolTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.RevokeRole(&_InvoicePool.TransactOpts, role, account)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 newFeeBps) returns()
func (_InvoicePool *InvoicePoolTransactor) SetPlatformFee(opts *bind.TransactOpts, newFeeBps *big.Int) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "setPlatformFee", newFeeBps)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 newFeeBps) returns()
func (_InvoicePool *InvoicePoolSession) SetPlatformFee(newFeeBps *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.SetPlatformFee(&_InvoicePool.TransactOpts, newFeeBps)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 newFeeBps) returns()
func (_InvoicePool *InvoicePoolTransactorSession) SetPlatformFee(newFeeBps *big.Int) (*types.Transaction, error) {
	return _InvoicePool.Contract.SetPlatformFee(&_InvoicePool.TransactOpts, newFeeBps)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address newWallet) returns()
func (_InvoicePool *InvoicePoolTransactor) SetPlatformWallet(opts *bind.TransactOpts, newWallet common.Address) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "setPlatformWallet", newWallet)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address newWallet) returns()
func (_InvoicePool *InvoicePoolSession) SetPlatformWallet(newWallet common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.SetPlatformWallet(&_InvoicePool.TransactOpts, newWallet)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address newWallet) returns()
func (_InvoicePool *InvoicePoolTransactorSession) SetPlatformWallet(newWallet common.Address) (*types.Transaction, error) {
	return _InvoicePool.Contract.SetPlatformWallet(&_InvoicePool.TransactOpts, newWallet)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoicePool *InvoicePoolTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoicePool.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoicePool *InvoicePoolSession) Unpause() (*types.Transaction, error) {
	return _InvoicePool.Contract.Unpause(&_InvoicePool.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoicePool *InvoicePoolTransactorSession) Unpause() (*types.Transaction, error) {
	return _InvoicePool.Contract.Unpause(&_InvoicePool.TransactOpts)
}

// InvoicePoolDisbursementRecordedIterator is returned from FilterDisbursementRecorded and is used to iterate over the raw logs and unpacked data for DisbursementRecorded events raised by the InvoicePool contract.
type InvoicePoolDisbursementRecordedIterator struct {
	Event *InvoicePoolDisbursementRecorded // Event containing the contract specifics and raw log

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
func (it *InvoicePoolDisbursementRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolDisbursementRecorded)
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
		it.Event = new(InvoicePoolDisbursementRecorded)
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
func (it *InvoicePoolDisbursementRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolDisbursementRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolDisbursementRecorded represents a DisbursementRecorded event raised by the InvoicePool contract.
type InvoicePoolDisbursementRecorded struct {
	TokenId  *big.Int
	Exporter common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDisbursementRecorded is a free log retrieval operation binding the contract event 0xa057bd151be659ccb58b76ac146344667076c4c3d3959ba95106bbe227c095e9.
//
// Solidity: event DisbursementRecorded(uint256 indexed tokenId, address indexed exporter, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) FilterDisbursementRecorded(opts *bind.FilterOpts, tokenId []*big.Int, exporter []common.Address) (*InvoicePoolDisbursementRecordedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var exporterRule []interface{}
	for _, exporterItem := range exporter {
		exporterRule = append(exporterRule, exporterItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "DisbursementRecorded", tokenIdRule, exporterRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolDisbursementRecordedIterator{contract: _InvoicePool.contract, event: "DisbursementRecorded", logs: logs, sub: sub}, nil
}

// WatchDisbursementRecorded is a free log subscription operation binding the contract event 0xa057bd151be659ccb58b76ac146344667076c4c3d3959ba95106bbe227c095e9.
//
// Solidity: event DisbursementRecorded(uint256 indexed tokenId, address indexed exporter, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) WatchDisbursementRecorded(opts *bind.WatchOpts, sink chan<- *InvoicePoolDisbursementRecorded, tokenId []*big.Int, exporter []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var exporterRule []interface{}
	for _, exporterItem := range exporter {
		exporterRule = append(exporterRule, exporterItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "DisbursementRecorded", tokenIdRule, exporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolDisbursementRecorded)
				if err := _InvoicePool.contract.UnpackLog(event, "DisbursementRecorded", log); err != nil {
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

// ParseDisbursementRecorded is a log parse operation binding the contract event 0xa057bd151be659ccb58b76ac146344667076c4c3d3959ba95106bbe227c095e9.
//
// Solidity: event DisbursementRecorded(uint256 indexed tokenId, address indexed exporter, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) ParseDisbursementRecorded(log types.Log) (*InvoicePoolDisbursementRecorded, error) {
	event := new(InvoicePoolDisbursementRecorded)
	if err := _InvoicePool.contract.UnpackLog(event, "DisbursementRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolExcessRepaymentRecordedIterator is returned from FilterExcessRepaymentRecorded and is used to iterate over the raw logs and unpacked data for ExcessRepaymentRecorded events raised by the InvoicePool contract.
type InvoicePoolExcessRepaymentRecordedIterator struct {
	Event *InvoicePoolExcessRepaymentRecorded // Event containing the contract specifics and raw log

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
func (it *InvoicePoolExcessRepaymentRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolExcessRepaymentRecorded)
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
		it.Event = new(InvoicePoolExcessRepaymentRecorded)
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
func (it *InvoicePoolExcessRepaymentRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolExcessRepaymentRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolExcessRepaymentRecorded represents a ExcessRepaymentRecorded event raised by the InvoicePool contract.
type InvoicePoolExcessRepaymentRecorded struct {
	TokenId   *big.Int
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterExcessRepaymentRecorded is a free log retrieval operation binding the contract event 0xb3765d358a4d74dfb6bb4959cc097beac5deb6a122e26154f6f9e3bbf2b89851.
//
// Solidity: event ExcessRepaymentRecorded(uint256 indexed tokenId, address indexed recipient, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) FilterExcessRepaymentRecorded(opts *bind.FilterOpts, tokenId []*big.Int, recipient []common.Address) (*InvoicePoolExcessRepaymentRecordedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "ExcessRepaymentRecorded", tokenIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolExcessRepaymentRecordedIterator{contract: _InvoicePool.contract, event: "ExcessRepaymentRecorded", logs: logs, sub: sub}, nil
}

// WatchExcessRepaymentRecorded is a free log subscription operation binding the contract event 0xb3765d358a4d74dfb6bb4959cc097beac5deb6a122e26154f6f9e3bbf2b89851.
//
// Solidity: event ExcessRepaymentRecorded(uint256 indexed tokenId, address indexed recipient, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) WatchExcessRepaymentRecorded(opts *bind.WatchOpts, sink chan<- *InvoicePoolExcessRepaymentRecorded, tokenId []*big.Int, recipient []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "ExcessRepaymentRecorded", tokenIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolExcessRepaymentRecorded)
				if err := _InvoicePool.contract.UnpackLog(event, "ExcessRepaymentRecorded", log); err != nil {
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

// ParseExcessRepaymentRecorded is a log parse operation binding the contract event 0xb3765d358a4d74dfb6bb4959cc097beac5deb6a122e26154f6f9e3bbf2b89851.
//
// Solidity: event ExcessRepaymentRecorded(uint256 indexed tokenId, address indexed recipient, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) ParseExcessRepaymentRecorded(log types.Log) (*InvoicePoolExcessRepaymentRecorded, error) {
	event := new(InvoicePoolExcessRepaymentRecorded)
	if err := _InvoicePool.contract.UnpackLog(event, "ExcessRepaymentRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolInvestmentRecordedIterator is returned from FilterInvestmentRecorded and is used to iterate over the raw logs and unpacked data for InvestmentRecorded events raised by the InvoicePool contract.
type InvoicePoolInvestmentRecordedIterator struct {
	Event *InvoicePoolInvestmentRecorded // Event containing the contract specifics and raw log

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
func (it *InvoicePoolInvestmentRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolInvestmentRecorded)
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
		it.Event = new(InvoicePoolInvestmentRecorded)
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
func (it *InvoicePoolInvestmentRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolInvestmentRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolInvestmentRecorded represents a InvestmentRecorded event raised by the InvoicePool contract.
type InvoicePoolInvestmentRecorded struct {
	TokenId        *big.Int
	Investor       common.Address
	Amount         *big.Int
	ExpectedReturn *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterInvestmentRecorded is a free log retrieval operation binding the contract event 0x4350cce8440dabd6bb58e4f0c9d02cfc7912ad433db358d824624522733754f2.
//
// Solidity: event InvestmentRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount, uint256 expectedReturn)
func (_InvoicePool *InvoicePoolFilterer) FilterInvestmentRecorded(opts *bind.FilterOpts, tokenId []*big.Int, investor []common.Address) (*InvoicePoolInvestmentRecordedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "InvestmentRecorded", tokenIdRule, investorRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolInvestmentRecordedIterator{contract: _InvoicePool.contract, event: "InvestmentRecorded", logs: logs, sub: sub}, nil
}

// WatchInvestmentRecorded is a free log subscription operation binding the contract event 0x4350cce8440dabd6bb58e4f0c9d02cfc7912ad433db358d824624522733754f2.
//
// Solidity: event InvestmentRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount, uint256 expectedReturn)
func (_InvoicePool *InvoicePoolFilterer) WatchInvestmentRecorded(opts *bind.WatchOpts, sink chan<- *InvoicePoolInvestmentRecorded, tokenId []*big.Int, investor []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "InvestmentRecorded", tokenIdRule, investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolInvestmentRecorded)
				if err := _InvoicePool.contract.UnpackLog(event, "InvestmentRecorded", log); err != nil {
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

// ParseInvestmentRecorded is a log parse operation binding the contract event 0x4350cce8440dabd6bb58e4f0c9d02cfc7912ad433db358d824624522733754f2.
//
// Solidity: event InvestmentRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount, uint256 expectedReturn)
func (_InvoicePool *InvoicePoolFilterer) ParseInvestmentRecorded(log types.Log) (*InvoicePoolInvestmentRecorded, error) {
	event := new(InvoicePoolInvestmentRecorded)
	if err := _InvoicePool.contract.UnpackLog(event, "InvestmentRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolInvestorReturnRecordedIterator is returned from FilterInvestorReturnRecorded and is used to iterate over the raw logs and unpacked data for InvestorReturnRecorded events raised by the InvoicePool contract.
type InvoicePoolInvestorReturnRecordedIterator struct {
	Event *InvoicePoolInvestorReturnRecorded // Event containing the contract specifics and raw log

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
func (it *InvoicePoolInvestorReturnRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolInvestorReturnRecorded)
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
		it.Event = new(InvoicePoolInvestorReturnRecorded)
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
func (it *InvoicePoolInvestorReturnRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolInvestorReturnRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolInvestorReturnRecorded represents a InvestorReturnRecorded event raised by the InvoicePool contract.
type InvoicePoolInvestorReturnRecorded struct {
	TokenId  *big.Int
	Investor common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInvestorReturnRecorded is a free log retrieval operation binding the contract event 0x235bb080582355af3f9131b865d54eabe7caf5b5ee4131760ad5ab02b74fc494.
//
// Solidity: event InvestorReturnRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) FilterInvestorReturnRecorded(opts *bind.FilterOpts, tokenId []*big.Int, investor []common.Address) (*InvoicePoolInvestorReturnRecordedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "InvestorReturnRecorded", tokenIdRule, investorRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolInvestorReturnRecordedIterator{contract: _InvoicePool.contract, event: "InvestorReturnRecorded", logs: logs, sub: sub}, nil
}

// WatchInvestorReturnRecorded is a free log subscription operation binding the contract event 0x235bb080582355af3f9131b865d54eabe7caf5b5ee4131760ad5ab02b74fc494.
//
// Solidity: event InvestorReturnRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) WatchInvestorReturnRecorded(opts *bind.WatchOpts, sink chan<- *InvoicePoolInvestorReturnRecorded, tokenId []*big.Int, investor []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "InvestorReturnRecorded", tokenIdRule, investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolInvestorReturnRecorded)
				if err := _InvoicePool.contract.UnpackLog(event, "InvestorReturnRecorded", log); err != nil {
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

// ParseInvestorReturnRecorded is a log parse operation binding the contract event 0x235bb080582355af3f9131b865d54eabe7caf5b5ee4131760ad5ab02b74fc494.
//
// Solidity: event InvestorReturnRecorded(uint256 indexed tokenId, address indexed investor, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) ParseInvestorReturnRecorded(log types.Log) (*InvoicePoolInvestorReturnRecorded, error) {
	event := new(InvoicePoolInvestorReturnRecorded)
	if err := _InvoicePool.contract.UnpackLog(event, "InvestorReturnRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the InvoicePool contract.
type InvoicePoolPausedIterator struct {
	Event *InvoicePoolPaused // Event containing the contract specifics and raw log

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
func (it *InvoicePoolPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolPaused)
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
		it.Event = new(InvoicePoolPaused)
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
func (it *InvoicePoolPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolPaused represents a Paused event raised by the InvoicePool contract.
type InvoicePoolPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_InvoicePool *InvoicePoolFilterer) FilterPaused(opts *bind.FilterOpts) (*InvoicePoolPausedIterator, error) {

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &InvoicePoolPausedIterator{contract: _InvoicePool.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_InvoicePool *InvoicePoolFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *InvoicePoolPaused) (event.Subscription, error) {

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolPaused)
				if err := _InvoicePool.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_InvoicePool *InvoicePoolFilterer) ParsePaused(log types.Log) (*InvoicePoolPaused, error) {
	event := new(InvoicePoolPaused)
	if err := _InvoicePool.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolPoolClosedIterator is returned from FilterPoolClosed and is used to iterate over the raw logs and unpacked data for PoolClosed events raised by the InvoicePool contract.
type InvoicePoolPoolClosedIterator struct {
	Event *InvoicePoolPoolClosed // Event containing the contract specifics and raw log

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
func (it *InvoicePoolPoolClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolPoolClosed)
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
		it.Event = new(InvoicePoolPoolClosed)
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
func (it *InvoicePoolPoolClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolPoolClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolPoolClosed represents a PoolClosed event raised by the InvoicePool contract.
type InvoicePoolPoolClosed struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPoolClosed is a free log retrieval operation binding the contract event 0x925a19753e677c9dc36a80e0fc824ca0c5b1afde494872b43daccab9ffeaffd4.
//
// Solidity: event PoolClosed(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) FilterPoolClosed(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoicePoolPoolClosedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "PoolClosed", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolPoolClosedIterator{contract: _InvoicePool.contract, event: "PoolClosed", logs: logs, sub: sub}, nil
}

// WatchPoolClosed is a free log subscription operation binding the contract event 0x925a19753e677c9dc36a80e0fc824ca0c5b1afde494872b43daccab9ffeaffd4.
//
// Solidity: event PoolClosed(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) WatchPoolClosed(opts *bind.WatchOpts, sink chan<- *InvoicePoolPoolClosed, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "PoolClosed", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolPoolClosed)
				if err := _InvoicePool.contract.UnpackLog(event, "PoolClosed", log); err != nil {
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

// ParsePoolClosed is a log parse operation binding the contract event 0x925a19753e677c9dc36a80e0fc824ca0c5b1afde494872b43daccab9ffeaffd4.
//
// Solidity: event PoolClosed(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) ParsePoolClosed(log types.Log) (*InvoicePoolPoolClosed, error) {
	event := new(InvoicePoolPoolClosed)
	if err := _InvoicePool.contract.UnpackLog(event, "PoolClosed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the InvoicePool contract.
type InvoicePoolPoolCreatedIterator struct {
	Event *InvoicePoolPoolCreated // Event containing the contract specifics and raw log

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
func (it *InvoicePoolPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolPoolCreated)
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
		it.Event = new(InvoicePoolPoolCreated)
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
func (it *InvoicePoolPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolPoolCreated represents a PoolCreated event raised by the InvoicePool contract.
type InvoicePoolPoolCreated struct {
	TokenId      *big.Int
	TargetAmount *big.Int
	InterestRate *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0x6a0c7fbf44f6331867816b75328f586816c7ff60b5f3b71d7ccd1da786a93898.
//
// Solidity: event PoolCreated(uint256 indexed tokenId, uint256 targetAmount, uint256 interestRate)
func (_InvoicePool *InvoicePoolFilterer) FilterPoolCreated(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoicePoolPoolCreatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "PoolCreated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolPoolCreatedIterator{contract: _InvoicePool.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0x6a0c7fbf44f6331867816b75328f586816c7ff60b5f3b71d7ccd1da786a93898.
//
// Solidity: event PoolCreated(uint256 indexed tokenId, uint256 targetAmount, uint256 interestRate)
func (_InvoicePool *InvoicePoolFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *InvoicePoolPoolCreated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "PoolCreated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolPoolCreated)
				if err := _InvoicePool.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0x6a0c7fbf44f6331867816b75328f586816c7ff60b5f3b71d7ccd1da786a93898.
//
// Solidity: event PoolCreated(uint256 indexed tokenId, uint256 targetAmount, uint256 interestRate)
func (_InvoicePool *InvoicePoolFilterer) ParsePoolCreated(log types.Log) (*InvoicePoolPoolCreated, error) {
	event := new(InvoicePoolPoolCreated)
	if err := _InvoicePool.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolPoolDefaultedIterator is returned from FilterPoolDefaulted and is used to iterate over the raw logs and unpacked data for PoolDefaulted events raised by the InvoicePool contract.
type InvoicePoolPoolDefaultedIterator struct {
	Event *InvoicePoolPoolDefaulted // Event containing the contract specifics and raw log

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
func (it *InvoicePoolPoolDefaultedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolPoolDefaulted)
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
		it.Event = new(InvoicePoolPoolDefaulted)
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
func (it *InvoicePoolPoolDefaultedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolPoolDefaultedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolPoolDefaulted represents a PoolDefaulted event raised by the InvoicePool contract.
type InvoicePoolPoolDefaulted struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPoolDefaulted is a free log retrieval operation binding the contract event 0x308716355dc28c7f34a1242f45560112cd36ac144d797726ee0ae38ae86e2ff1.
//
// Solidity: event PoolDefaulted(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) FilterPoolDefaulted(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoicePoolPoolDefaultedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "PoolDefaulted", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolPoolDefaultedIterator{contract: _InvoicePool.contract, event: "PoolDefaulted", logs: logs, sub: sub}, nil
}

// WatchPoolDefaulted is a free log subscription operation binding the contract event 0x308716355dc28c7f34a1242f45560112cd36ac144d797726ee0ae38ae86e2ff1.
//
// Solidity: event PoolDefaulted(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) WatchPoolDefaulted(opts *bind.WatchOpts, sink chan<- *InvoicePoolPoolDefaulted, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "PoolDefaulted", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolPoolDefaulted)
				if err := _InvoicePool.contract.UnpackLog(event, "PoolDefaulted", log); err != nil {
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

// ParsePoolDefaulted is a log parse operation binding the contract event 0x308716355dc28c7f34a1242f45560112cd36ac144d797726ee0ae38ae86e2ff1.
//
// Solidity: event PoolDefaulted(uint256 indexed tokenId)
func (_InvoicePool *InvoicePoolFilterer) ParsePoolDefaulted(log types.Log) (*InvoicePoolPoolDefaulted, error) {
	event := new(InvoicePoolPoolDefaulted)
	if err := _InvoicePool.contract.UnpackLog(event, "PoolDefaulted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolPoolFilledIterator is returned from FilterPoolFilled and is used to iterate over the raw logs and unpacked data for PoolFilled events raised by the InvoicePool contract.
type InvoicePoolPoolFilledIterator struct {
	Event *InvoicePoolPoolFilled // Event containing the contract specifics and raw log

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
func (it *InvoicePoolPoolFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolPoolFilled)
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
		it.Event = new(InvoicePoolPoolFilled)
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
func (it *InvoicePoolPoolFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolPoolFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolPoolFilled represents a PoolFilled event raised by the InvoicePool contract.
type InvoicePoolPoolFilled struct {
	TokenId       *big.Int
	TotalAmount   *big.Int
	InvestorCount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPoolFilled is a free log retrieval operation binding the contract event 0x2f175f06e9e98a031b67651e940b202eb5594b8c635aeb92bdac6fce470183d4.
//
// Solidity: event PoolFilled(uint256 indexed tokenId, uint256 totalAmount, uint256 investorCount)
func (_InvoicePool *InvoicePoolFilterer) FilterPoolFilled(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoicePoolPoolFilledIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "PoolFilled", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolPoolFilledIterator{contract: _InvoicePool.contract, event: "PoolFilled", logs: logs, sub: sub}, nil
}

// WatchPoolFilled is a free log subscription operation binding the contract event 0x2f175f06e9e98a031b67651e940b202eb5594b8c635aeb92bdac6fce470183d4.
//
// Solidity: event PoolFilled(uint256 indexed tokenId, uint256 totalAmount, uint256 investorCount)
func (_InvoicePool *InvoicePoolFilterer) WatchPoolFilled(opts *bind.WatchOpts, sink chan<- *InvoicePoolPoolFilled, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "PoolFilled", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolPoolFilled)
				if err := _InvoicePool.contract.UnpackLog(event, "PoolFilled", log); err != nil {
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

// ParsePoolFilled is a log parse operation binding the contract event 0x2f175f06e9e98a031b67651e940b202eb5594b8c635aeb92bdac6fce470183d4.
//
// Solidity: event PoolFilled(uint256 indexed tokenId, uint256 totalAmount, uint256 investorCount)
func (_InvoicePool *InvoicePoolFilterer) ParsePoolFilled(log types.Log) (*InvoicePoolPoolFilled, error) {
	event := new(InvoicePoolPoolFilled)
	if err := _InvoicePool.contract.UnpackLog(event, "PoolFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolRepaymentRecordedIterator is returned from FilterRepaymentRecorded and is used to iterate over the raw logs and unpacked data for RepaymentRecorded events raised by the InvoicePool contract.
type InvoicePoolRepaymentRecordedIterator struct {
	Event *InvoicePoolRepaymentRecorded // Event containing the contract specifics and raw log

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
func (it *InvoicePoolRepaymentRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolRepaymentRecorded)
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
		it.Event = new(InvoicePoolRepaymentRecorded)
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
func (it *InvoicePoolRepaymentRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolRepaymentRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolRepaymentRecorded represents a RepaymentRecorded event raised by the InvoicePool contract.
type InvoicePoolRepaymentRecorded struct {
	TokenId *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRepaymentRecorded is a free log retrieval operation binding the contract event 0x5dcbf9c927d185b37be08114b8c07f3c61be903c4f8bee4a6cc9494d3d105739.
//
// Solidity: event RepaymentRecorded(uint256 indexed tokenId, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) FilterRepaymentRecorded(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoicePoolRepaymentRecordedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "RepaymentRecorded", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolRepaymentRecordedIterator{contract: _InvoicePool.contract, event: "RepaymentRecorded", logs: logs, sub: sub}, nil
}

// WatchRepaymentRecorded is a free log subscription operation binding the contract event 0x5dcbf9c927d185b37be08114b8c07f3c61be903c4f8bee4a6cc9494d3d105739.
//
// Solidity: event RepaymentRecorded(uint256 indexed tokenId, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) WatchRepaymentRecorded(opts *bind.WatchOpts, sink chan<- *InvoicePoolRepaymentRecorded, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "RepaymentRecorded", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolRepaymentRecorded)
				if err := _InvoicePool.contract.UnpackLog(event, "RepaymentRecorded", log); err != nil {
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

// ParseRepaymentRecorded is a log parse operation binding the contract event 0x5dcbf9c927d185b37be08114b8c07f3c61be903c4f8bee4a6cc9494d3d105739.
//
// Solidity: event RepaymentRecorded(uint256 indexed tokenId, uint256 amount)
func (_InvoicePool *InvoicePoolFilterer) ParseRepaymentRecorded(log types.Log) (*InvoicePoolRepaymentRecorded, error) {
	event := new(InvoicePoolRepaymentRecorded)
	if err := _InvoicePool.contract.UnpackLog(event, "RepaymentRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the InvoicePool contract.
type InvoicePoolRoleAdminChangedIterator struct {
	Event *InvoicePoolRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *InvoicePoolRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolRoleAdminChanged)
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
		it.Event = new(InvoicePoolRoleAdminChanged)
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
func (it *InvoicePoolRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolRoleAdminChanged represents a RoleAdminChanged event raised by the InvoicePool contract.
type InvoicePoolRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_InvoicePool *InvoicePoolFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*InvoicePoolRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolRoleAdminChangedIterator{contract: _InvoicePool.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_InvoicePool *InvoicePoolFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *InvoicePoolRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolRoleAdminChanged)
				if err := _InvoicePool.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_InvoicePool *InvoicePoolFilterer) ParseRoleAdminChanged(log types.Log) (*InvoicePoolRoleAdminChanged, error) {
	event := new(InvoicePoolRoleAdminChanged)
	if err := _InvoicePool.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the InvoicePool contract.
type InvoicePoolRoleGrantedIterator struct {
	Event *InvoicePoolRoleGranted // Event containing the contract specifics and raw log

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
func (it *InvoicePoolRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolRoleGranted)
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
		it.Event = new(InvoicePoolRoleGranted)
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
func (it *InvoicePoolRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolRoleGranted represents a RoleGranted event raised by the InvoicePool contract.
type InvoicePoolRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*InvoicePoolRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolRoleGrantedIterator{contract: _InvoicePool.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *InvoicePoolRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolRoleGranted)
				if err := _InvoicePool.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) ParseRoleGranted(log types.Log) (*InvoicePoolRoleGranted, error) {
	event := new(InvoicePoolRoleGranted)
	if err := _InvoicePool.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the InvoicePool contract.
type InvoicePoolRoleRevokedIterator struct {
	Event *InvoicePoolRoleRevoked // Event containing the contract specifics and raw log

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
func (it *InvoicePoolRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolRoleRevoked)
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
		it.Event = new(InvoicePoolRoleRevoked)
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
func (it *InvoicePoolRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolRoleRevoked represents a RoleRevoked event raised by the InvoicePool contract.
type InvoicePoolRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*InvoicePoolRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &InvoicePoolRoleRevokedIterator{contract: _InvoicePool.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *InvoicePoolRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolRoleRevoked)
				if err := _InvoicePool.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoicePool *InvoicePoolFilterer) ParseRoleRevoked(log types.Log) (*InvoicePoolRoleRevoked, error) {
	event := new(InvoicePoolRoleRevoked)
	if err := _InvoicePool.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoicePoolUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the InvoicePool contract.
type InvoicePoolUnpausedIterator struct {
	Event *InvoicePoolUnpaused // Event containing the contract specifics and raw log

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
func (it *InvoicePoolUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoicePoolUnpaused)
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
		it.Event = new(InvoicePoolUnpaused)
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
func (it *InvoicePoolUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoicePoolUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoicePoolUnpaused represents a Unpaused event raised by the InvoicePool contract.
type InvoicePoolUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_InvoicePool *InvoicePoolFilterer) FilterUnpaused(opts *bind.FilterOpts) (*InvoicePoolUnpausedIterator, error) {

	logs, sub, err := _InvoicePool.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &InvoicePoolUnpausedIterator{contract: _InvoicePool.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_InvoicePool *InvoicePoolFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *InvoicePoolUnpaused) (event.Subscription, error) {

	logs, sub, err := _InvoicePool.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoicePoolUnpaused)
				if err := _InvoicePool.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_InvoicePool *InvoicePoolFilterer) ParseUnpaused(log types.Log) (*InvoicePoolUnpaused, error) {
	event := new(InvoicePoolUnpaused)
	if err := _InvoicePool.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
