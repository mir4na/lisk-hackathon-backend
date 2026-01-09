package services

import (
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

// AuthServiceInterface defines the contract for authentication operations
type AuthServiceInterface interface {
	Register(req *models.RegisterRequest) (*models.LoginResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(refreshToken string) (*models.LoginResponse, error)
}

// InvoiceServiceInterface defines the contract for invoice operations
type InvoiceServiceInterface interface {
	Create(exporterID uuid.UUID, req *models.CreateInvoiceRequest) (*models.Invoice, error)
	GetByID(id uuid.UUID) (*models.Invoice, error)
	GetByExporter(exporterID uuid.UUID, filter *models.InvoiceFilter) (*models.InvoiceListResponse, error)
	GetFundable(page, perPage int) (*models.InvoiceListResponse, error)
	Update(id, exporterID uuid.UUID, req *models.UpdateInvoiceRequest) (*models.Invoice, error)
	Delete(id, exporterID uuid.UUID) error
	Submit(id, exporterID uuid.UUID) error
	Approve(id uuid.UUID, interestRate float64) error
	Reject(id uuid.UUID, reason string) error
	UploadDocument(invoiceID, exporterID uuid.UUID, docType models.DocumentType, fileData []byte, fileName string) (*models.InvoiceDocument, error)
	GetDocuments(invoiceID uuid.UUID) ([]models.InvoiceDocument, error)
	DeleteDocument(docID, exporterID uuid.UUID) error
}

// FundingServiceInterface defines the contract for funding operations
type FundingServiceInterface interface {
	CreatePool(invoiceID uuid.UUID) (*models.FundingPool, error)
	GetPool(poolID uuid.UUID) (*models.FundingPoolResponse, error)
	GetOpenPools(page, perPage int) (*models.PoolListResponse, error)
	Invest(investorID uuid.UUID, req *models.InvestRequest) (*models.Investment, error)
	GetInvestmentsByInvestor(investorID uuid.UUID, page, perPage int) (*models.InvestmentListResponse, error)
	DisburseToExporter(poolID uuid.UUID) error
	ProcessRepayment(invoiceID uuid.UUID, amount float64) error
	ClosePoolAndNotifyExporter(poolID uuid.UUID) (*models.ExporterPaymentNotificationData, error)
}

// BlockchainServiceInterface defines the contract for blockchain operations
type BlockchainServiceInterface interface {
	TokenizeInvoice(invoiceID uuid.UUID, ownerAddress string) (*models.InvoiceNFT, error)
	BurnNFT(invoiceID uuid.UUID) error
	GetNFTByInvoice(invoiceID uuid.UUID) (*models.InvoiceNFT, error)
	TransferNFT(nftID uuid.UUID, newOwner string) error
}

// PinataServiceInterface defines the contract for IPFS operations
type PinataServiceInterface interface {
	UploadFile(fileData []byte, fileName string, metadata map[string]string) (*PinataResponse, string, error)
	UploadJSON(data interface{}, name string) (*PinataResponse, error)
	GetIPFSURL(hash string) string
	UploadNFTMetadata(metadata *NFTMetadata) (string, error)
}

// Ensure implementations satisfy interfaces
var _ AuthServiceInterface = (*AuthService)(nil)
var _ InvoiceServiceInterface = (*InvoiceService)(nil)
var _ FundingServiceInterface = (*FundingService)(nil)
var _ PinataServiceInterface = (*PinataService)(nil)
