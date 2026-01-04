package repository

import (
	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/models"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	Create(user *models.User, profile *models.UserProfile) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindProfileByUserID(userID uuid.UUID) (*models.UserProfile, error)
	UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error
	UpdateWallet(userID uuid.UUID, walletAddress string) error
	SetVerified(userID uuid.UUID, verified bool) error
	EmailExists(email string) (bool, error)
	WalletExists(wallet string) (bool, error)
}

// KYCRepositoryInterface defines the contract for KYC data operations
type KYCRepositoryInterface interface {
	Create(kyc *models.KYCVerification) error
	FindByID(id uuid.UUID) (*models.KYCVerification, error)
	FindByUserID(userID uuid.UUID) (*models.KYCVerification, error)
	FindPending(page, perPage int) ([]models.KYCVerification, int, error)
	Approve(kycID, adminID uuid.UUID) error
	Reject(kycID, adminID uuid.UUID, reason string) error
	UpdateDocumentURL(kycID uuid.UUID, docURL string) error
	UpdateSelfieURL(kycID uuid.UUID, selfieURL string) error
}

// BuyerRepositoryInterface defines the contract for buyer data operations
type BuyerRepositoryInterface interface {
	Create(buyer *models.Buyer) error
	FindByID(id uuid.UUID) (*models.Buyer, error)
	FindByExporter(exporterID uuid.UUID, page, perPage int) ([]models.Buyer, int, error)
	Update(buyer *models.Buyer) error
	Delete(id uuid.UUID) error
	UpdateCreditScore(id uuid.UUID, score int) error
	IncrementInvoiceStats(id uuid.UUID, amount float64, isPaid bool) error
	SetVerified(id uuid.UUID, verified bool) error
}

// InvoiceRepositoryInterface defines the contract for invoice data operations
type InvoiceRepositoryInterface interface {
	Create(invoice *models.Invoice) error
	FindByID(id uuid.UUID) (*models.Invoice, error)
	FindByExporter(exporterID uuid.UUID, filter *models.InvoiceFilter) ([]models.Invoice, int, error)
	FindFundable(page, perPage int) ([]models.Invoice, int, error)
	Update(invoice *models.Invoice) error
	UpdateStatus(id uuid.UUID, status models.InvoiceStatus) error
	SetInterestRate(id uuid.UUID, rate float64) error
	SetAdvanceAmount(id uuid.UUID, amount float64) error
	SetDocumentHash(id uuid.UUID, hash string) error
	Delete(id uuid.UUID) error

	// Document methods
	CreateDocument(doc *models.InvoiceDocument) error
	FindDocumentsByInvoiceID(invoiceID uuid.UUID) ([]models.InvoiceDocument, error)
	DeleteDocument(id uuid.UUID) error

	// NFT methods
	CreateNFT(nft *models.InvoiceNFT) error
	FindNFTByInvoiceID(invoiceID uuid.UUID) (*models.InvoiceNFT, error)
	UpdateNFTOwner(id uuid.UUID, owner string) error
	BurnNFT(id uuid.UUID, txHash string) error

	// Transaction methods
	ApproveWithTransaction(id uuid.UUID, interestRate, advanceAmount float64) error
}

// FundingRepositoryInterface defines the contract for funding data operations
type FundingRepositoryInterface interface {
	// Pool methods
	CreatePool(pool *models.FundingPool) error
	FindPoolByID(id uuid.UUID) (*models.FundingPool, error)
	FindPoolByInvoiceID(invoiceID uuid.UUID) (*models.FundingPool, error)
	FindOpenPools(page, perPage int) ([]models.FundingPool, int, error)
	UpdatePoolFunding(id uuid.UUID, amount float64) error
	UpdatePoolStatus(id uuid.UUID, status models.PoolStatus) error

	// Investment methods
	CreateInvestment(inv *models.Investment) error
	FindInvestmentByID(id uuid.UUID) (*models.Investment, error)
	FindInvestmentsByInvestor(investorID uuid.UUID, page, perPage int) ([]models.Investment, int, error)
	FindInvestmentsByPool(poolID uuid.UUID) ([]models.Investment, error)
	UpdateInvestmentStatus(id uuid.UUID, status models.InvestmentStatus, actualReturn *float64) error
}

// TransactionRepositoryInterface defines the contract for transaction data operations
type TransactionRepositoryInterface interface {
	Create(tx *models.Transaction) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindByTxHash(txHash string) (*models.Transaction, error)
	FindByUser(userID uuid.UUID, page, perPage int) ([]models.Transaction, int, error)
	FindByInvoice(invoiceID uuid.UUID) ([]models.Transaction, error)
	UpdateStatus(id uuid.UUID, status models.TransactionStatus) error
	UpdateBlockInfo(id uuid.UUID, blockNumber, gasUsed int64) error
}

// Ensure implementations satisfy interfaces
var _ UserRepositoryInterface = (*UserRepository)(nil)
var _ KYCRepositoryInterface = (*KYCRepository)(nil)
var _ BuyerRepositoryInterface = (*BuyerRepository)(nil)
var _ InvoiceRepositoryInterface = (*InvoiceRepository)(nil)
var _ FundingRepositoryInterface = (*FundingRepository)(nil)
var _ TransactionRepositoryInterface = (*TransactionRepository)(nil)
