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

// InvoiceNFTInvoice is an auto generated low-level Go binding around an user-defined struct.
type InvoiceNFTInvoice struct {
	InvoiceNumber    string
	Amount           *big.Int
	AdvanceAmount    *big.Int
	InterestRate     *big.Int
	IssueDate        *big.Int
	DueDate          *big.Int
	Exporter         common.Address
	BuyerCountry     string
	DocumentHash     string
	Status           uint8
	ShipmentVerified bool
}

// InvoiceNFTMetaData contains all meta data concerning the InvoiceNFT contract.
var InvoiceNFTMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"InvoiceBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"invoiceNumber\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"}],\"name\":\"InvoiceMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumInvoiceNFT.InvoiceStatus\",\"name\":\"oldStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumInvoiceNFT.InvoiceStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"InvoiceStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"ShipmentVerified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"burnInvoice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"exporterInvoices\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"}],\"name\":\"getExporterInvoices\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getInvoice\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"invoiceNumber\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"advanceAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issueDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"buyerCountry\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"documentHash\",\"type\":\"string\"},{\"internalType\":\"enumInvoiceNFT.InvoiceStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"shipmentVerified\",\"type\":\"bool\"}],\"internalType\":\"structInvoiceNFT.Invoice\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"invoiceNumber\",\"type\":\"string\"}],\"name\":\"getTokenIdByInvoiceNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"invoiceNumberToTokenId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"invoices\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"invoiceNumber\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"advanceAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issueDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exporter\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"buyerCountry\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"documentHash\",\"type\":\"string\"},{\"internalType\":\"enumInvoiceNFT.InvoiceStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"shipmentVerified\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isFundable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"invoiceNumber\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"advanceAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issueDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueDate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"buyerCountry\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"documentHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"mintInvoice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalMinted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumInvoiceNFT.InvoiceStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"updateStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"verifyShipment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// InvoiceNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use InvoiceNFTMetaData.ABI instead.
var InvoiceNFTABI = InvoiceNFTMetaData.ABI

// InvoiceNFT is an auto generated Go binding around an Ethereum contract.
type InvoiceNFT struct {
	InvoiceNFTCaller     // Read-only binding to the contract
	InvoiceNFTTransactor // Write-only binding to the contract
	InvoiceNFTFilterer   // Log filterer for contract events
}

// InvoiceNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type InvoiceNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoiceNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InvoiceNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoiceNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InvoiceNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvoiceNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InvoiceNFTSession struct {
	Contract     *InvoiceNFT       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InvoiceNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InvoiceNFTCallerSession struct {
	Contract *InvoiceNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// InvoiceNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InvoiceNFTTransactorSession struct {
	Contract     *InvoiceNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// InvoiceNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type InvoiceNFTRaw struct {
	Contract *InvoiceNFT // Generic contract binding to access the raw methods on
}

// InvoiceNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InvoiceNFTCallerRaw struct {
	Contract *InvoiceNFTCaller // Generic read-only contract binding to access the raw methods on
}

// InvoiceNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InvoiceNFTTransactorRaw struct {
	Contract *InvoiceNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInvoiceNFT creates a new instance of InvoiceNFT, bound to a specific deployed contract.
func NewInvoiceNFT(address common.Address, backend bind.ContractBackend) (*InvoiceNFT, error) {
	contract, err := bindInvoiceNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFT{InvoiceNFTCaller: InvoiceNFTCaller{contract: contract}, InvoiceNFTTransactor: InvoiceNFTTransactor{contract: contract}, InvoiceNFTFilterer: InvoiceNFTFilterer{contract: contract}}, nil
}

// NewInvoiceNFTCaller creates a new read-only instance of InvoiceNFT, bound to a specific deployed contract.
func NewInvoiceNFTCaller(address common.Address, caller bind.ContractCaller) (*InvoiceNFTCaller, error) {
	contract, err := bindInvoiceNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTCaller{contract: contract}, nil
}

// NewInvoiceNFTTransactor creates a new write-only instance of InvoiceNFT, bound to a specific deployed contract.
func NewInvoiceNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*InvoiceNFTTransactor, error) {
	contract, err := bindInvoiceNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTTransactor{contract: contract}, nil
}

// NewInvoiceNFTFilterer creates a new log filterer instance of InvoiceNFT, bound to a specific deployed contract.
func NewInvoiceNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*InvoiceNFTFilterer, error) {
	contract, err := bindInvoiceNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTFilterer{contract: contract}, nil
}

// bindInvoiceNFT binds a generic wrapper to an already deployed contract.
func bindInvoiceNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InvoiceNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InvoiceNFT *InvoiceNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InvoiceNFT.Contract.InvoiceNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InvoiceNFT *InvoiceNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.InvoiceNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InvoiceNFT *InvoiceNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.InvoiceNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InvoiceNFT *InvoiceNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InvoiceNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InvoiceNFT *InvoiceNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InvoiceNFT *InvoiceNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.DEFAULTADMINROLE(&_InvoiceNFT.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.DEFAULTADMINROLE(&_InvoiceNFT.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTSession) MINTERROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.MINTERROLE(&_InvoiceNFT.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCallerSession) MINTERROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.MINTERROLE(&_InvoiceNFT.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTSession) ORACLEROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.ORACLEROLE(&_InvoiceNFT.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCallerSession) ORACLEROLE() ([32]byte, error) {
	return _InvoiceNFT.Contract.ORACLEROLE(&_InvoiceNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _InvoiceNFT.Contract.BalanceOf(&_InvoiceNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _InvoiceNFT.Contract.BalanceOf(&_InvoiceNFT.CallOpts, owner)
}

// ExporterInvoices is a free data retrieval call binding the contract method 0x655924ba.
//
// Solidity: function exporterInvoices(address , uint256 ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCaller) ExporterInvoices(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "exporterInvoices", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExporterInvoices is a free data retrieval call binding the contract method 0x655924ba.
//
// Solidity: function exporterInvoices(address , uint256 ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) ExporterInvoices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _InvoiceNFT.Contract.ExporterInvoices(&_InvoiceNFT.CallOpts, arg0, arg1)
}

// ExporterInvoices is a free data retrieval call binding the contract method 0x655924ba.
//
// Solidity: function exporterInvoices(address , uint256 ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCallerSession) ExporterInvoices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _InvoiceNFT.Contract.ExporterInvoices(&_InvoiceNFT.CallOpts, arg0, arg1)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _InvoiceNFT.Contract.GetApproved(&_InvoiceNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _InvoiceNFT.Contract.GetApproved(&_InvoiceNFT.CallOpts, tokenId)
}

// GetExporterInvoices is a free data retrieval call binding the contract method 0x2c6a3891.
//
// Solidity: function getExporterInvoices(address exporter) view returns(uint256[])
func (_InvoiceNFT *InvoiceNFTCaller) GetExporterInvoices(opts *bind.CallOpts, exporter common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "getExporterInvoices", exporter)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetExporterInvoices is a free data retrieval call binding the contract method 0x2c6a3891.
//
// Solidity: function getExporterInvoices(address exporter) view returns(uint256[])
func (_InvoiceNFT *InvoiceNFTSession) GetExporterInvoices(exporter common.Address) ([]*big.Int, error) {
	return _InvoiceNFT.Contract.GetExporterInvoices(&_InvoiceNFT.CallOpts, exporter)
}

// GetExporterInvoices is a free data retrieval call binding the contract method 0x2c6a3891.
//
// Solidity: function getExporterInvoices(address exporter) view returns(uint256[])
func (_InvoiceNFT *InvoiceNFTCallerSession) GetExporterInvoices(exporter common.Address) ([]*big.Int, error) {
	return _InvoiceNFT.Contract.GetExporterInvoices(&_InvoiceNFT.CallOpts, exporter)
}

// GetInvoice is a free data retrieval call binding the contract method 0x3a23cc0a.
//
// Solidity: function getInvoice(uint256 tokenId) view returns((string,uint256,uint256,uint256,uint256,uint256,address,string,string,uint8,bool))
func (_InvoiceNFT *InvoiceNFTCaller) GetInvoice(opts *bind.CallOpts, tokenId *big.Int) (InvoiceNFTInvoice, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "getInvoice", tokenId)

	if err != nil {
		return *new(InvoiceNFTInvoice), err
	}

	out0 := *abi.ConvertType(out[0], new(InvoiceNFTInvoice)).(*InvoiceNFTInvoice)

	return out0, err

}

// GetInvoice is a free data retrieval call binding the contract method 0x3a23cc0a.
//
// Solidity: function getInvoice(uint256 tokenId) view returns((string,uint256,uint256,uint256,uint256,uint256,address,string,string,uint8,bool))
func (_InvoiceNFT *InvoiceNFTSession) GetInvoice(tokenId *big.Int) (InvoiceNFTInvoice, error) {
	return _InvoiceNFT.Contract.GetInvoice(&_InvoiceNFT.CallOpts, tokenId)
}

// GetInvoice is a free data retrieval call binding the contract method 0x3a23cc0a.
//
// Solidity: function getInvoice(uint256 tokenId) view returns((string,uint256,uint256,uint256,uint256,uint256,address,string,string,uint8,bool))
func (_InvoiceNFT *InvoiceNFTCallerSession) GetInvoice(tokenId *big.Int) (InvoiceNFTInvoice, error) {
	return _InvoiceNFT.Contract.GetInvoice(&_InvoiceNFT.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _InvoiceNFT.Contract.GetRoleAdmin(&_InvoiceNFT.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_InvoiceNFT *InvoiceNFTCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _InvoiceNFT.Contract.GetRoleAdmin(&_InvoiceNFT.CallOpts, role)
}

// GetTokenIdByInvoiceNumber is a free data retrieval call binding the contract method 0xe65e36c8.
//
// Solidity: function getTokenIdByInvoiceNumber(string invoiceNumber) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCaller) GetTokenIdByInvoiceNumber(opts *bind.CallOpts, invoiceNumber string) (*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "getTokenIdByInvoiceNumber", invoiceNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenIdByInvoiceNumber is a free data retrieval call binding the contract method 0xe65e36c8.
//
// Solidity: function getTokenIdByInvoiceNumber(string invoiceNumber) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) GetTokenIdByInvoiceNumber(invoiceNumber string) (*big.Int, error) {
	return _InvoiceNFT.Contract.GetTokenIdByInvoiceNumber(&_InvoiceNFT.CallOpts, invoiceNumber)
}

// GetTokenIdByInvoiceNumber is a free data retrieval call binding the contract method 0xe65e36c8.
//
// Solidity: function getTokenIdByInvoiceNumber(string invoiceNumber) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCallerSession) GetTokenIdByInvoiceNumber(invoiceNumber string) (*big.Int, error) {
	return _InvoiceNFT.Contract.GetTokenIdByInvoiceNumber(&_InvoiceNFT.CallOpts, invoiceNumber)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoiceNFT *InvoiceNFTSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _InvoiceNFT.Contract.HasRole(&_InvoiceNFT.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _InvoiceNFT.Contract.HasRole(&_InvoiceNFT.CallOpts, role, account)
}

// InvoiceNumberToTokenId is a free data retrieval call binding the contract method 0x6ea88abe.
//
// Solidity: function invoiceNumberToTokenId(string ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCaller) InvoiceNumberToTokenId(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "invoiceNumberToTokenId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InvoiceNumberToTokenId is a free data retrieval call binding the contract method 0x6ea88abe.
//
// Solidity: function invoiceNumberToTokenId(string ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) InvoiceNumberToTokenId(arg0 string) (*big.Int, error) {
	return _InvoiceNFT.Contract.InvoiceNumberToTokenId(&_InvoiceNFT.CallOpts, arg0)
}

// InvoiceNumberToTokenId is a free data retrieval call binding the contract method 0x6ea88abe.
//
// Solidity: function invoiceNumberToTokenId(string ) view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCallerSession) InvoiceNumberToTokenId(arg0 string) (*big.Int, error) {
	return _InvoiceNFT.Contract.InvoiceNumberToTokenId(&_InvoiceNFT.CallOpts, arg0)
}

// Invoices is a free data retrieval call binding the contract method 0x4e6d1405.
//
// Solidity: function invoices(uint256 ) view returns(string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, address exporter, string buyerCountry, string documentHash, uint8 status, bool shipmentVerified)
func (_InvoiceNFT *InvoiceNFTCaller) Invoices(opts *bind.CallOpts, arg0 *big.Int) (struct {
	InvoiceNumber    string
	Amount           *big.Int
	AdvanceAmount    *big.Int
	InterestRate     *big.Int
	IssueDate        *big.Int
	DueDate          *big.Int
	Exporter         common.Address
	BuyerCountry     string
	DocumentHash     string
	Status           uint8
	ShipmentVerified bool
}, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "invoices", arg0)

	outstruct := new(struct {
		InvoiceNumber    string
		Amount           *big.Int
		AdvanceAmount    *big.Int
		InterestRate     *big.Int
		IssueDate        *big.Int
		DueDate          *big.Int
		Exporter         common.Address
		BuyerCountry     string
		DocumentHash     string
		Status           uint8
		ShipmentVerified bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.InvoiceNumber = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AdvanceAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.InterestRate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IssueDate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.DueDate = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Exporter = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.BuyerCountry = *abi.ConvertType(out[7], new(string)).(*string)
	outstruct.DocumentHash = *abi.ConvertType(out[8], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[9], new(uint8)).(*uint8)
	outstruct.ShipmentVerified = *abi.ConvertType(out[10], new(bool)).(*bool)

	return *outstruct, err

}

// Invoices is a free data retrieval call binding the contract method 0x4e6d1405.
//
// Solidity: function invoices(uint256 ) view returns(string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, address exporter, string buyerCountry, string documentHash, uint8 status, bool shipmentVerified)
func (_InvoiceNFT *InvoiceNFTSession) Invoices(arg0 *big.Int) (struct {
	InvoiceNumber    string
	Amount           *big.Int
	AdvanceAmount    *big.Int
	InterestRate     *big.Int
	IssueDate        *big.Int
	DueDate          *big.Int
	Exporter         common.Address
	BuyerCountry     string
	DocumentHash     string
	Status           uint8
	ShipmentVerified bool
}, error) {
	return _InvoiceNFT.Contract.Invoices(&_InvoiceNFT.CallOpts, arg0)
}

// Invoices is a free data retrieval call binding the contract method 0x4e6d1405.
//
// Solidity: function invoices(uint256 ) view returns(string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, address exporter, string buyerCountry, string documentHash, uint8 status, bool shipmentVerified)
func (_InvoiceNFT *InvoiceNFTCallerSession) Invoices(arg0 *big.Int) (struct {
	InvoiceNumber    string
	Amount           *big.Int
	AdvanceAmount    *big.Int
	InterestRate     *big.Int
	IssueDate        *big.Int
	DueDate          *big.Int
	Exporter         common.Address
	BuyerCountry     string
	DocumentHash     string
	Status           uint8
	ShipmentVerified bool
}, error) {
	return _InvoiceNFT.Contract.Invoices(&_InvoiceNFT.CallOpts, arg0)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_InvoiceNFT *InvoiceNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _InvoiceNFT.Contract.IsApprovedForAll(&_InvoiceNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _InvoiceNFT.Contract.IsApprovedForAll(&_InvoiceNFT.CallOpts, owner, operator)
}

// IsFundable is a free data retrieval call binding the contract method 0xa92e2886.
//
// Solidity: function isFundable(uint256 tokenId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCaller) IsFundable(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "isFundable", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFundable is a free data retrieval call binding the contract method 0xa92e2886.
//
// Solidity: function isFundable(uint256 tokenId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTSession) IsFundable(tokenId *big.Int) (bool, error) {
	return _InvoiceNFT.Contract.IsFundable(&_InvoiceNFT.CallOpts, tokenId)
}

// IsFundable is a free data retrieval call binding the contract method 0xa92e2886.
//
// Solidity: function isFundable(uint256 tokenId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCallerSession) IsFundable(tokenId *big.Int) (bool, error) {
	return _InvoiceNFT.Contract.IsFundable(&_InvoiceNFT.CallOpts, tokenId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InvoiceNFT *InvoiceNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InvoiceNFT *InvoiceNFTSession) Name() (string, error) {
	return _InvoiceNFT.Contract.Name(&_InvoiceNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_InvoiceNFT *InvoiceNFTCallerSession) Name() (string, error) {
	return _InvoiceNFT.Contract.Name(&_InvoiceNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _InvoiceNFT.Contract.OwnerOf(&_InvoiceNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_InvoiceNFT *InvoiceNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _InvoiceNFT.Contract.OwnerOf(&_InvoiceNFT.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoiceNFT *InvoiceNFTCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoiceNFT *InvoiceNFTSession) Paused() (bool, error) {
	return _InvoiceNFT.Contract.Paused(&_InvoiceNFT.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_InvoiceNFT *InvoiceNFTCallerSession) Paused() (bool, error) {
	return _InvoiceNFT.Contract.Paused(&_InvoiceNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _InvoiceNFT.Contract.SupportsInterface(&_InvoiceNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_InvoiceNFT *InvoiceNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _InvoiceNFT.Contract.SupportsInterface(&_InvoiceNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_InvoiceNFT *InvoiceNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_InvoiceNFT *InvoiceNFTSession) Symbol() (string, error) {
	return _InvoiceNFT.Contract.Symbol(&_InvoiceNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_InvoiceNFT *InvoiceNFTCallerSession) Symbol() (string, error) {
	return _InvoiceNFT.Contract.Symbol(&_InvoiceNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_InvoiceNFT *InvoiceNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_InvoiceNFT *InvoiceNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _InvoiceNFT.Contract.TokenURI(&_InvoiceNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_InvoiceNFT *InvoiceNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _InvoiceNFT.Contract.TokenURI(&_InvoiceNFT.CallOpts, tokenId)
}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCaller) TotalMinted(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InvoiceNFT.contract.Call(opts, &out, "totalMinted")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) TotalMinted() (*big.Int, error) {
	return _InvoiceNFT.Contract.TotalMinted(&_InvoiceNFT.CallOpts)
}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_InvoiceNFT *InvoiceNFTCallerSession) TotalMinted() (*big.Int, error) {
	return _InvoiceNFT.Contract.TotalMinted(&_InvoiceNFT.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Approve(&_InvoiceNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Approve(&_InvoiceNFT.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Burn(&_InvoiceNFT.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Burn(&_InvoiceNFT.TransactOpts, tokenId)
}

// BurnInvoice is a paid mutator transaction binding the contract method 0x0192d67d.
//
// Solidity: function burnInvoice(uint256 tokenId, string reason) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) BurnInvoice(opts *bind.TransactOpts, tokenId *big.Int, reason string) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "burnInvoice", tokenId, reason)
}

// BurnInvoice is a paid mutator transaction binding the contract method 0x0192d67d.
//
// Solidity: function burnInvoice(uint256 tokenId, string reason) returns()
func (_InvoiceNFT *InvoiceNFTSession) BurnInvoice(tokenId *big.Int, reason string) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.BurnInvoice(&_InvoiceNFT.TransactOpts, tokenId, reason)
}

// BurnInvoice is a paid mutator transaction binding the contract method 0x0192d67d.
//
// Solidity: function burnInvoice(uint256 tokenId, string reason) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) BurnInvoice(tokenId *big.Int, reason string) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.BurnInvoice(&_InvoiceNFT.TransactOpts, tokenId, reason)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.GrantRole(&_InvoiceNFT.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.GrantRole(&_InvoiceNFT.TransactOpts, role, account)
}

// MintInvoice is a paid mutator transaction binding the contract method 0x1c45e7a1.
//
// Solidity: function mintInvoice(address to, string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, string buyerCountry, string documentHash, string uri) returns(uint256)
func (_InvoiceNFT *InvoiceNFTTransactor) MintInvoice(opts *bind.TransactOpts, to common.Address, invoiceNumber string, amount *big.Int, advanceAmount *big.Int, interestRate *big.Int, issueDate *big.Int, dueDate *big.Int, buyerCountry string, documentHash string, uri string) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "mintInvoice", to, invoiceNumber, amount, advanceAmount, interestRate, issueDate, dueDate, buyerCountry, documentHash, uri)
}

// MintInvoice is a paid mutator transaction binding the contract method 0x1c45e7a1.
//
// Solidity: function mintInvoice(address to, string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, string buyerCountry, string documentHash, string uri) returns(uint256)
func (_InvoiceNFT *InvoiceNFTSession) MintInvoice(to common.Address, invoiceNumber string, amount *big.Int, advanceAmount *big.Int, interestRate *big.Int, issueDate *big.Int, dueDate *big.Int, buyerCountry string, documentHash string, uri string) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.MintInvoice(&_InvoiceNFT.TransactOpts, to, invoiceNumber, amount, advanceAmount, interestRate, issueDate, dueDate, buyerCountry, documentHash, uri)
}

// MintInvoice is a paid mutator transaction binding the contract method 0x1c45e7a1.
//
// Solidity: function mintInvoice(address to, string invoiceNumber, uint256 amount, uint256 advanceAmount, uint256 interestRate, uint256 issueDate, uint256 dueDate, string buyerCountry, string documentHash, string uri) returns(uint256)
func (_InvoiceNFT *InvoiceNFTTransactorSession) MintInvoice(to common.Address, invoiceNumber string, amount *big.Int, advanceAmount *big.Int, interestRate *big.Int, issueDate *big.Int, dueDate *big.Int, buyerCountry string, documentHash string, uri string) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.MintInvoice(&_InvoiceNFT.TransactOpts, to, invoiceNumber, amount, advanceAmount, interestRate, issueDate, dueDate, buyerCountry, documentHash, uri)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoiceNFT *InvoiceNFTTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoiceNFT *InvoiceNFTSession) Pause() (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Pause(&_InvoiceNFT.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) Pause() (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Pause(&_InvoiceNFT.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoiceNFT *InvoiceNFTSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.RenounceRole(&_InvoiceNFT.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.RenounceRole(&_InvoiceNFT.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.RevokeRole(&_InvoiceNFT.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.RevokeRole(&_InvoiceNFT.TransactOpts, role, account)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SafeTransferFrom(&_InvoiceNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SafeTransferFrom(&_InvoiceNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_InvoiceNFT *InvoiceNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SafeTransferFrom0(&_InvoiceNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SafeTransferFrom0(&_InvoiceNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_InvoiceNFT *InvoiceNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SetApprovalForAll(&_InvoiceNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.SetApprovalForAll(&_InvoiceNFT.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.TransferFrom(&_InvoiceNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.TransferFrom(&_InvoiceNFT.TransactOpts, from, to, tokenId)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoiceNFT *InvoiceNFTTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoiceNFT *InvoiceNFTSession) Unpause() (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Unpause(&_InvoiceNFT.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) Unpause() (*types.Transaction, error) {
	return _InvoiceNFT.Contract.Unpause(&_InvoiceNFT.TransactOpts)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x3a1b3d31.
//
// Solidity: function updateStatus(uint256 tokenId, uint8 newStatus) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) UpdateStatus(opts *bind.TransactOpts, tokenId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "updateStatus", tokenId, newStatus)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x3a1b3d31.
//
// Solidity: function updateStatus(uint256 tokenId, uint8 newStatus) returns()
func (_InvoiceNFT *InvoiceNFTSession) UpdateStatus(tokenId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.UpdateStatus(&_InvoiceNFT.TransactOpts, tokenId, newStatus)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x3a1b3d31.
//
// Solidity: function updateStatus(uint256 tokenId, uint8 newStatus) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) UpdateStatus(tokenId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.UpdateStatus(&_InvoiceNFT.TransactOpts, tokenId, newStatus)
}

// VerifyShipment is a paid mutator transaction binding the contract method 0x92d15fc4.
//
// Solidity: function verifyShipment(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactor) VerifyShipment(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.contract.Transact(opts, "verifyShipment", tokenId)
}

// VerifyShipment is a paid mutator transaction binding the contract method 0x92d15fc4.
//
// Solidity: function verifyShipment(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTSession) VerifyShipment(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.VerifyShipment(&_InvoiceNFT.TransactOpts, tokenId)
}

// VerifyShipment is a paid mutator transaction binding the contract method 0x92d15fc4.
//
// Solidity: function verifyShipment(uint256 tokenId) returns()
func (_InvoiceNFT *InvoiceNFTTransactorSession) VerifyShipment(tokenId *big.Int) (*types.Transaction, error) {
	return _InvoiceNFT.Contract.VerifyShipment(&_InvoiceNFT.TransactOpts, tokenId)
}

// InvoiceNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the InvoiceNFT contract.
type InvoiceNFTApprovalIterator struct {
	Event *InvoiceNFTApproval // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTApproval)
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
		it.Event = new(InvoiceNFTApproval)
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
func (it *InvoiceNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTApproval represents a Approval event raised by the InvoiceNFT contract.
type InvoiceNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*InvoiceNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTApprovalIterator{contract: _InvoiceNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *InvoiceNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTApproval)
				if err := _InvoiceNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseApproval(log types.Log) (*InvoiceNFTApproval, error) {
	event := new(InvoiceNFTApproval)
	if err := _InvoiceNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the InvoiceNFT contract.
type InvoiceNFTApprovalForAllIterator struct {
	Event *InvoiceNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTApprovalForAll)
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
		it.Event = new(InvoiceNFTApprovalForAll)
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
func (it *InvoiceNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTApprovalForAll represents a ApprovalForAll event raised by the InvoiceNFT contract.
type InvoiceNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*InvoiceNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTApprovalForAllIterator{contract: _InvoiceNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *InvoiceNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTApprovalForAll)
				if err := _InvoiceNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseApprovalForAll(log types.Log) (*InvoiceNFTApprovalForAll, error) {
	event := new(InvoiceNFTApprovalForAll)
	if err := _InvoiceNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the InvoiceNFT contract.
type InvoiceNFTBatchMetadataUpdateIterator struct {
	Event *InvoiceNFTBatchMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTBatchMetadataUpdate)
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
		it.Event = new(InvoiceNFTBatchMetadataUpdate)
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
func (it *InvoiceNFTBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the InvoiceNFT contract.
type InvoiceNFTBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*InvoiceNFTBatchMetadataUpdateIterator, error) {

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTBatchMetadataUpdateIterator{contract: _InvoiceNFT.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *InvoiceNFTBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTBatchMetadataUpdate)
				if err := _InvoiceNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
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

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseBatchMetadataUpdate(log types.Log) (*InvoiceNFTBatchMetadataUpdate, error) {
	event := new(InvoiceNFTBatchMetadataUpdate)
	if err := _InvoiceNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTInvoiceBurnedIterator is returned from FilterInvoiceBurned and is used to iterate over the raw logs and unpacked data for InvoiceBurned events raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceBurnedIterator struct {
	Event *InvoiceNFTInvoiceBurned // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTInvoiceBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTInvoiceBurned)
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
		it.Event = new(InvoiceNFTInvoiceBurned)
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
func (it *InvoiceNFTInvoiceBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTInvoiceBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTInvoiceBurned represents a InvoiceBurned event raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceBurned struct {
	TokenId *big.Int
	Reason  string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInvoiceBurned is a free log retrieval operation binding the contract event 0xf87c3983a8e51190c8bf32fad4356a45da34d97b6a933423e8f46af9a432a5c3.
//
// Solidity: event InvoiceBurned(uint256 indexed tokenId, string reason)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterInvoiceBurned(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoiceNFTInvoiceBurnedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "InvoiceBurned", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTInvoiceBurnedIterator{contract: _InvoiceNFT.contract, event: "InvoiceBurned", logs: logs, sub: sub}, nil
}

// WatchInvoiceBurned is a free log subscription operation binding the contract event 0xf87c3983a8e51190c8bf32fad4356a45da34d97b6a933423e8f46af9a432a5c3.
//
// Solidity: event InvoiceBurned(uint256 indexed tokenId, string reason)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchInvoiceBurned(opts *bind.WatchOpts, sink chan<- *InvoiceNFTInvoiceBurned, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "InvoiceBurned", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTInvoiceBurned)
				if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceBurned", log); err != nil {
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

// ParseInvoiceBurned is a log parse operation binding the contract event 0xf87c3983a8e51190c8bf32fad4356a45da34d97b6a933423e8f46af9a432a5c3.
//
// Solidity: event InvoiceBurned(uint256 indexed tokenId, string reason)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseInvoiceBurned(log types.Log) (*InvoiceNFTInvoiceBurned, error) {
	event := new(InvoiceNFTInvoiceBurned)
	if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTInvoiceMintedIterator is returned from FilterInvoiceMinted and is used to iterate over the raw logs and unpacked data for InvoiceMinted events raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceMintedIterator struct {
	Event *InvoiceNFTInvoiceMinted // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTInvoiceMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTInvoiceMinted)
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
		it.Event = new(InvoiceNFTInvoiceMinted)
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
func (it *InvoiceNFTInvoiceMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTInvoiceMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTInvoiceMinted represents a InvoiceMinted event raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceMinted struct {
	TokenId       *big.Int
	Exporter      common.Address
	InvoiceNumber string
	Amount        *big.Int
	DueDate       *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterInvoiceMinted is a free log retrieval operation binding the contract event 0xf6618f6e0aa9f042b2c68fbf59c637f6d4dfd66267f8b99901527f860d736724.
//
// Solidity: event InvoiceMinted(uint256 indexed tokenId, address indexed exporter, string invoiceNumber, uint256 amount, uint256 dueDate)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterInvoiceMinted(opts *bind.FilterOpts, tokenId []*big.Int, exporter []common.Address) (*InvoiceNFTInvoiceMintedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var exporterRule []interface{}
	for _, exporterItem := range exporter {
		exporterRule = append(exporterRule, exporterItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "InvoiceMinted", tokenIdRule, exporterRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTInvoiceMintedIterator{contract: _InvoiceNFT.contract, event: "InvoiceMinted", logs: logs, sub: sub}, nil
}

// WatchInvoiceMinted is a free log subscription operation binding the contract event 0xf6618f6e0aa9f042b2c68fbf59c637f6d4dfd66267f8b99901527f860d736724.
//
// Solidity: event InvoiceMinted(uint256 indexed tokenId, address indexed exporter, string invoiceNumber, uint256 amount, uint256 dueDate)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchInvoiceMinted(opts *bind.WatchOpts, sink chan<- *InvoiceNFTInvoiceMinted, tokenId []*big.Int, exporter []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var exporterRule []interface{}
	for _, exporterItem := range exporter {
		exporterRule = append(exporterRule, exporterItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "InvoiceMinted", tokenIdRule, exporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTInvoiceMinted)
				if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceMinted", log); err != nil {
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

// ParseInvoiceMinted is a log parse operation binding the contract event 0xf6618f6e0aa9f042b2c68fbf59c637f6d4dfd66267f8b99901527f860d736724.
//
// Solidity: event InvoiceMinted(uint256 indexed tokenId, address indexed exporter, string invoiceNumber, uint256 amount, uint256 dueDate)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseInvoiceMinted(log types.Log) (*InvoiceNFTInvoiceMinted, error) {
	event := new(InvoiceNFTInvoiceMinted)
	if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTInvoiceStatusChangedIterator is returned from FilterInvoiceStatusChanged and is used to iterate over the raw logs and unpacked data for InvoiceStatusChanged events raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceStatusChangedIterator struct {
	Event *InvoiceNFTInvoiceStatusChanged // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTInvoiceStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTInvoiceStatusChanged)
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
		it.Event = new(InvoiceNFTInvoiceStatusChanged)
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
func (it *InvoiceNFTInvoiceStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTInvoiceStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTInvoiceStatusChanged represents a InvoiceStatusChanged event raised by the InvoiceNFT contract.
type InvoiceNFTInvoiceStatusChanged struct {
	TokenId   *big.Int
	OldStatus uint8
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvoiceStatusChanged is a free log retrieval operation binding the contract event 0x1bd9ff6dd5bac6bcdb836cff31ce376515b7ed6eebade1982ce3ab2d428d4948.
//
// Solidity: event InvoiceStatusChanged(uint256 indexed tokenId, uint8 oldStatus, uint8 newStatus)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterInvoiceStatusChanged(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoiceNFTInvoiceStatusChangedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "InvoiceStatusChanged", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTInvoiceStatusChangedIterator{contract: _InvoiceNFT.contract, event: "InvoiceStatusChanged", logs: logs, sub: sub}, nil
}

// WatchInvoiceStatusChanged is a free log subscription operation binding the contract event 0x1bd9ff6dd5bac6bcdb836cff31ce376515b7ed6eebade1982ce3ab2d428d4948.
//
// Solidity: event InvoiceStatusChanged(uint256 indexed tokenId, uint8 oldStatus, uint8 newStatus)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchInvoiceStatusChanged(opts *bind.WatchOpts, sink chan<- *InvoiceNFTInvoiceStatusChanged, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "InvoiceStatusChanged", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTInvoiceStatusChanged)
				if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceStatusChanged", log); err != nil {
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

// ParseInvoiceStatusChanged is a log parse operation binding the contract event 0x1bd9ff6dd5bac6bcdb836cff31ce376515b7ed6eebade1982ce3ab2d428d4948.
//
// Solidity: event InvoiceStatusChanged(uint256 indexed tokenId, uint8 oldStatus, uint8 newStatus)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseInvoiceStatusChanged(log types.Log) (*InvoiceNFTInvoiceStatusChanged, error) {
	event := new(InvoiceNFTInvoiceStatusChanged)
	if err := _InvoiceNFT.contract.UnpackLog(event, "InvoiceStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the InvoiceNFT contract.
type InvoiceNFTMetadataUpdateIterator struct {
	Event *InvoiceNFTMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTMetadataUpdate)
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
		it.Event = new(InvoiceNFTMetadataUpdate)
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
func (it *InvoiceNFTMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTMetadataUpdate represents a MetadataUpdate event raised by the InvoiceNFT contract.
type InvoiceNFTMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*InvoiceNFTMetadataUpdateIterator, error) {

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTMetadataUpdateIterator{contract: _InvoiceNFT.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *InvoiceNFTMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTMetadataUpdate)
				if err := _InvoiceNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseMetadataUpdate(log types.Log) (*InvoiceNFTMetadataUpdate, error) {
	event := new(InvoiceNFTMetadataUpdate)
	if err := _InvoiceNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the InvoiceNFT contract.
type InvoiceNFTPausedIterator struct {
	Event *InvoiceNFTPaused // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTPaused)
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
		it.Event = new(InvoiceNFTPaused)
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
func (it *InvoiceNFTPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTPaused represents a Paused event raised by the InvoiceNFT contract.
type InvoiceNFTPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterPaused(opts *bind.FilterOpts) (*InvoiceNFTPausedIterator, error) {

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTPausedIterator{contract: _InvoiceNFT.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *InvoiceNFTPaused) (event.Subscription, error) {

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTPaused)
				if err := _InvoiceNFT.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_InvoiceNFT *InvoiceNFTFilterer) ParsePaused(log types.Log) (*InvoiceNFTPaused, error) {
	event := new(InvoiceNFTPaused)
	if err := _InvoiceNFT.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the InvoiceNFT contract.
type InvoiceNFTRoleAdminChangedIterator struct {
	Event *InvoiceNFTRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTRoleAdminChanged)
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
		it.Event = new(InvoiceNFTRoleAdminChanged)
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
func (it *InvoiceNFTRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTRoleAdminChanged represents a RoleAdminChanged event raised by the InvoiceNFT contract.
type InvoiceNFTRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*InvoiceNFTRoleAdminChangedIterator, error) {

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

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTRoleAdminChangedIterator{contract: _InvoiceNFT.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *InvoiceNFTRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTRoleAdminChanged)
				if err := _InvoiceNFT.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_InvoiceNFT *InvoiceNFTFilterer) ParseRoleAdminChanged(log types.Log) (*InvoiceNFTRoleAdminChanged, error) {
	event := new(InvoiceNFTRoleAdminChanged)
	if err := _InvoiceNFT.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the InvoiceNFT contract.
type InvoiceNFTRoleGrantedIterator struct {
	Event *InvoiceNFTRoleGranted // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTRoleGranted)
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
		it.Event = new(InvoiceNFTRoleGranted)
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
func (it *InvoiceNFTRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTRoleGranted represents a RoleGranted event raised by the InvoiceNFT contract.
type InvoiceNFTRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*InvoiceNFTRoleGrantedIterator, error) {

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

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTRoleGrantedIterator{contract: _InvoiceNFT.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *InvoiceNFTRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTRoleGranted)
				if err := _InvoiceNFT.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_InvoiceNFT *InvoiceNFTFilterer) ParseRoleGranted(log types.Log) (*InvoiceNFTRoleGranted, error) {
	event := new(InvoiceNFTRoleGranted)
	if err := _InvoiceNFT.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the InvoiceNFT contract.
type InvoiceNFTRoleRevokedIterator struct {
	Event *InvoiceNFTRoleRevoked // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTRoleRevoked)
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
		it.Event = new(InvoiceNFTRoleRevoked)
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
func (it *InvoiceNFTRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTRoleRevoked represents a RoleRevoked event raised by the InvoiceNFT contract.
type InvoiceNFTRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*InvoiceNFTRoleRevokedIterator, error) {

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

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTRoleRevokedIterator{contract: _InvoiceNFT.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *InvoiceNFTRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTRoleRevoked)
				if err := _InvoiceNFT.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_InvoiceNFT *InvoiceNFTFilterer) ParseRoleRevoked(log types.Log) (*InvoiceNFTRoleRevoked, error) {
	event := new(InvoiceNFTRoleRevoked)
	if err := _InvoiceNFT.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTShipmentVerifiedIterator is returned from FilterShipmentVerified and is used to iterate over the raw logs and unpacked data for ShipmentVerified events raised by the InvoiceNFT contract.
type InvoiceNFTShipmentVerifiedIterator struct {
	Event *InvoiceNFTShipmentVerified // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTShipmentVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTShipmentVerified)
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
		it.Event = new(InvoiceNFTShipmentVerified)
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
func (it *InvoiceNFTShipmentVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTShipmentVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTShipmentVerified represents a ShipmentVerified event raised by the InvoiceNFT contract.
type InvoiceNFTShipmentVerified struct {
	TokenId  *big.Int
	Verifier common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterShipmentVerified is a free log retrieval operation binding the contract event 0x36c86a23ee3c04e849e7c64cd50dcc0bd810d35dab9792986864808eb1b8282c.
//
// Solidity: event ShipmentVerified(uint256 indexed tokenId, address verifier)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterShipmentVerified(opts *bind.FilterOpts, tokenId []*big.Int) (*InvoiceNFTShipmentVerifiedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "ShipmentVerified", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTShipmentVerifiedIterator{contract: _InvoiceNFT.contract, event: "ShipmentVerified", logs: logs, sub: sub}, nil
}

// WatchShipmentVerified is a free log subscription operation binding the contract event 0x36c86a23ee3c04e849e7c64cd50dcc0bd810d35dab9792986864808eb1b8282c.
//
// Solidity: event ShipmentVerified(uint256 indexed tokenId, address verifier)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchShipmentVerified(opts *bind.WatchOpts, sink chan<- *InvoiceNFTShipmentVerified, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "ShipmentVerified", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTShipmentVerified)
				if err := _InvoiceNFT.contract.UnpackLog(event, "ShipmentVerified", log); err != nil {
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

// ParseShipmentVerified is a log parse operation binding the contract event 0x36c86a23ee3c04e849e7c64cd50dcc0bd810d35dab9792986864808eb1b8282c.
//
// Solidity: event ShipmentVerified(uint256 indexed tokenId, address verifier)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseShipmentVerified(log types.Log) (*InvoiceNFTShipmentVerified, error) {
	event := new(InvoiceNFTShipmentVerified)
	if err := _InvoiceNFT.contract.UnpackLog(event, "ShipmentVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the InvoiceNFT contract.
type InvoiceNFTTransferIterator struct {
	Event *InvoiceNFTTransfer // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTTransfer)
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
		it.Event = new(InvoiceNFTTransfer)
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
func (it *InvoiceNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTTransfer represents a Transfer event raised by the InvoiceNFT contract.
type InvoiceNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*InvoiceNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTTransferIterator{contract: _InvoiceNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *InvoiceNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTTransfer)
				if err := _InvoiceNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_InvoiceNFT *InvoiceNFTFilterer) ParseTransfer(log types.Log) (*InvoiceNFTTransfer, error) {
	event := new(InvoiceNFTTransfer)
	if err := _InvoiceNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InvoiceNFTUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the InvoiceNFT contract.
type InvoiceNFTUnpausedIterator struct {
	Event *InvoiceNFTUnpaused // Event containing the contract specifics and raw log

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
func (it *InvoiceNFTUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvoiceNFTUnpaused)
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
		it.Event = new(InvoiceNFTUnpaused)
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
func (it *InvoiceNFTUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvoiceNFTUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvoiceNFTUnpaused represents a Unpaused event raised by the InvoiceNFT contract.
type InvoiceNFTUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_InvoiceNFT *InvoiceNFTFilterer) FilterUnpaused(opts *bind.FilterOpts) (*InvoiceNFTUnpausedIterator, error) {

	logs, sub, err := _InvoiceNFT.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &InvoiceNFTUnpausedIterator{contract: _InvoiceNFT.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_InvoiceNFT *InvoiceNFTFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *InvoiceNFTUnpaused) (event.Subscription, error) {

	logs, sub, err := _InvoiceNFT.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvoiceNFTUnpaused)
				if err := _InvoiceNFT.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_InvoiceNFT *InvoiceNFTFilterer) ParseUnpaused(log types.Log) (*InvoiceNFTUnpaused, error) {
	event := new(InvoiceNFTUnpaused)
	if err := _InvoiceNFT.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
