package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

type BlockchainService struct {
	client       *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	chainID      *big.Int
	nftContract  common.Address
	poolContract common.Address
	invoiceRepo  repository.InvoiceRepositoryInterface
	pinata       PinataServiceInterface
	cfg          *config.Config
}

func NewBlockchainService(cfg *config.Config, invoiceRepo repository.InvoiceRepositoryInterface, pinata PinataServiceInterface) (*BlockchainService, error) {
	if cfg.PrivateKey == "" {
		return &BlockchainService{
			invoiceRepo: invoiceRepo,
			pinata:      pinata,
			cfg:         cfg,
		}, nil
	}

	client, err := ethclient.Dial(cfg.BlockchainRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &BlockchainService{
		client:       client,
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		chainID:      big.NewInt(cfg.ChainID),
		nftContract:  common.HexToAddress(cfg.InvoiceNFTContractAddr),
		poolContract: common.HexToAddress(cfg.InvoicePoolContractAddr),
		invoiceRepo:  invoiceRepo,
		pinata:       pinata,
		cfg:          cfg,
	}, nil
}

func (s *BlockchainService) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	if s.client == nil {
		return nil, errors.New("blockchain client not initialized")
	}

	nonce, err := s.client.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (s *BlockchainService) PrepareNFTMetadata(invoice *models.Invoice) (*NFTMetadata, error) {
	metadata := &NFTMetadata{
		Name:        fmt.Sprintf("Invoice #%s", invoice.InvoiceNumber),
		Description: fmt.Sprintf("NFT representing invoice %s from exporter to %s", invoice.InvoiceNumber, invoice.Buyer.CompanyName),
		Attributes: []NFTAttribute{
			{TraitType: "Invoice Number", Value: invoice.InvoiceNumber},
			{TraitType: "Amount", Value: invoice.Amount, DisplayType: "number"},
			{TraitType: "Currency", Value: invoice.Currency},
			{TraitType: "Issue Date", Value: invoice.IssueDate.Format("2006-01-02"), DisplayType: "date"},
			{TraitType: "Due Date", Value: invoice.DueDate.Format("2006-01-02"), DisplayType: "date"},
			{TraitType: "Buyer", Value: invoice.Buyer.CompanyName},
			{TraitType: "Buyer Country", Value: invoice.Buyer.Country},
			{TraitType: "Buyer Credit Score", Value: invoice.Buyer.CreditScore, DisplayType: "number"},
		},
		Properties: map[string]interface{}{
			"invoice_id":    invoice.ID.String(),
			"exporter_id":   invoice.ExporterID.String(),
			"buyer_id":      invoice.BuyerID.String(),
			"document_hash": invoice.DocumentHash,
		},
	}

	if invoice.InterestRate != nil {
		metadata.Attributes = append(metadata.Attributes, NFTAttribute{
			TraitType:   "Interest Rate",
			Value:       *invoice.InterestRate,
			DisplayType: "number",
		})
	}

	if invoice.AdvanceAmount != nil {
		metadata.Attributes = append(metadata.Attributes, NFTAttribute{
			TraitType:   "Advance Amount",
			Value:       *invoice.AdvanceAmount,
			DisplayType: "number",
		})
	}

	return metadata, nil
}

func (s *BlockchainService) TokenizeInvoice(invoiceID uuid.UUID, ownerAddress string) (*models.InvoiceNFT, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.Status != models.StatusApproved {
		return nil, errors.New("invoice must be approved before tokenization")
	}

	// Prepare and upload NFT metadata
	metadata, err := s.PrepareNFTMetadata(invoice)
	if err != nil {
		return nil, err
	}

	metadataURI, err := s.pinata.UploadNFTMetadata(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to upload metadata: %w", err)
	}

	// For now, we'll simulate the minting process
	// In production, this would call the smart contract
	now := time.Now()
	tokenID := int64(time.Now().UnixNano() % 1000000) // Simulated token ID

	nft := &models.InvoiceNFT{
		InvoiceID:       invoiceID,
		TokenID:         &tokenID,
		ContractAddress: &s.cfg.InvoiceNFTContractAddr,
		ChainID:         int(s.cfg.ChainID),
		OwnerAddress:    &ownerAddress,
		MetadataURI:     &metadataURI,
		MintedAt:        &now,
	}

	// If blockchain client is available, mint on-chain
	if s.client != nil {
		// In a real implementation, we would:
		// 1. Call the smart contract's mint function
		// 2. Wait for transaction confirmation
		// 3. Get the actual token ID from the event logs
		// For now, we'll just record the pending state
		nft.MintTxHash = nil // Would be set after actual minting
	}

	if err := s.invoiceRepo.CreateNFT(nft); err != nil {
		return nil, err
	}

	// Update invoice status
	if err := s.invoiceRepo.UpdateStatus(invoiceID, models.StatusTokenized); err != nil {
		return nil, err
	}

	return nft, nil
}

func (s *BlockchainService) BurnNFT(invoiceID uuid.UUID) error {
	nft, err := s.invoiceRepo.FindNFTByInvoiceID(invoiceID)
	if err != nil {
		return err
	}
	if nft == nil {
		return errors.New("NFT not found")
	}
	if nft.BurnedAt != nil {
		return errors.New("NFT already burned")
	}

	// In production, call the smart contract's burn function
	// For now, just update the database
	return s.invoiceRepo.BurnNFT(nft.ID, "0x...simulated_burn_tx_hash")
}

func (s *BlockchainService) GetNFTByInvoice(invoiceID uuid.UUID) (*models.InvoiceNFT, error) {
	return s.invoiceRepo.FindNFTByInvoiceID(invoiceID)
}

func (s *BlockchainService) TransferNFT(nftID uuid.UUID, newOwner string) error {
	return s.invoiceRepo.UpdateNFTOwner(nftID, newOwner)
}

func (s *BlockchainService) GetBalance(address string) (*big.Int, error) {
	if s.client == nil {
		return big.NewInt(0), nil
	}

	account := common.HexToAddress(address)
	balance, err := s.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
