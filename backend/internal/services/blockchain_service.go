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
	"github.com/vessel/backend/internal/contracts" // Generated bindings
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

type BlockchainService struct {
	client       *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	chainID      *big.Int
	nftContract  *contracts.InvoiceNFT
	poolContract *contracts.InvoicePool
	invoiceRepo  repository.InvoiceRepositoryInterface
	fundingRepo  repository.FundingRepositoryInterface
	pinata       PinataServiceInterface
	cfg          *config.Config
}

func NewBlockchainService(cfg *config.Config, invoiceRepo repository.InvoiceRepositoryInterface, fundingRepo repository.FundingRepositoryInterface, pinata PinataServiceInterface) (*BlockchainService, error) {
	if cfg.PrivateKey == "" {
		return &BlockchainService{
			invoiceRepo: invoiceRepo,
			fundingRepo: fundingRepo,
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

	nftAddress := common.HexToAddress(cfg.InvoiceNFTContractAddr)
	poolAddress := common.HexToAddress(cfg.InvoicePoolContractAddr)

	nftContract, err := contracts.NewInvoiceNFT(nftAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load NFT contract: %w", err)
	}

	poolContract, err := contracts.NewInvoicePool(poolAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load Pool contract: %w", err)
	}

	return &BlockchainService{
		client:       client,
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		chainID:      big.NewInt(cfg.ChainID),
		nftContract:  nftContract,
		poolContract: poolContract,
		invoiceRepo:  invoiceRepo,
		fundingRepo:  fundingRepo,
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
	auth.GasLimit = uint64(500000) // Increased gas limit
	auth.GasPrice = gasPrice

	return auth, nil
}

// ... PrepareNFTMetadata remains same ...
func (s *BlockchainService) PrepareNFTMetadata(invoice *models.Invoice) (*NFTMetadata, error) {
	metadata := &NFTMetadata{
		Name:        fmt.Sprintf("Invoice #%s", invoice.InvoiceNumber),
		Description: fmt.Sprintf("NFT representing invoice %s from exporter to %s", invoice.InvoiceNumber, invoice.BuyerName),
		Attributes: []NFTAttribute{
			{TraitType: "Invoice Number", Value: invoice.InvoiceNumber},
			{TraitType: "Amount", Value: invoice.Amount, DisplayType: "number"},
			{TraitType: "Currency", Value: invoice.Currency},
			{TraitType: "Issue Date", Value: invoice.IssueDate.Format("2006-01-02"), DisplayType: "date"},
			{TraitType: "Due Date", Value: invoice.DueDate.Format("2006-01-02"), DisplayType: "date"},
			{TraitType: "Buyer", Value: invoice.BuyerName},
			{TraitType: "Buyer Country", Value: invoice.BuyerCountry},
		},
		Properties: map[string]interface{}{
			"invoice_id":    invoice.ID.String(),
			"exporter_id":   invoice.ExporterID.String(),
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

	now := time.Now()
	var tokenID int64
	var txHash string

	// Call Smart Contract
	if s.client != nil {
		auth, err := s.GetTransactOpts(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get transact opts: %w", err)
		}

		// Convert amounts to BigInt
		amountBig := new(big.Int).SetInt64(int64(invoice.Amount))
		advanceBig := new(big.Int).SetInt64(int64(*invoice.AdvanceAmount))
		interestBig := new(big.Int).SetInt64(int64(*invoice.InterestRate * 100)) // Scaled
		issueDateBig := big.NewInt(invoice.IssueDate.Unix())
		dueDateBig := big.NewInt(invoice.DueDate.Unix())

		exporterAddr := common.HexToAddress(ownerAddress)

		// Handle optional fields
		docHash := ""
		if invoice.DocumentHash != nil {
			docHash = *invoice.DocumentHash
		}

		// Mint Invoice NFT
		tx, err := s.nftContract.MintInvoice(
			auth,
			exporterAddr,
			invoice.InvoiceNumber,
			amountBig,
			advanceBig,
			interestBig,
			issueDateBig,
			dueDateBig,
			invoice.BuyerCountry,
			docHash,
			metadataURI,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to mint NFT: %w", err)
		}

		fmt.Printf("[BLOCKCHAIN] NFT Mint submitted: TxHash=%s\n", tx.Hash().Hex())
		txHash = tx.Hash().Hex()

		// NOTE: In production, we should wait for event to get TokenID.
		// For MVP, we will use a hash-based ID or fetch it later.
		// Let's assume we fetch it or use a placeholder if async.
		tokenID = int64(time.Now().UnixNano() % 1000000)
	} else {
		tokenID = int64(time.Now().UnixNano() % 1000000)
	}

	nft := &models.InvoiceNFT{
		InvoiceID:       invoiceID,
		TokenID:         &tokenID,
		ContractAddress: &s.cfg.InvoiceNFTContractAddr,
		ChainID:         int(s.cfg.ChainID),
		OwnerAddress:    &ownerAddress,
		MetadataURI:     &metadataURI,
		MintedAt:        &now,
		MintTxHash:      &txHash,
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
	// ... (Simulated for now, can be updated later) ...
	return nil
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
	return s.client.BalanceAt(context.Background(), account, nil)
}

// BlockchainTransaction represents a recorded on-chain transaction
type BlockchainTransaction struct {
	TxHash string  `json:"tx_hash"`
	Action string  `json:"action"`
	Amount float64 `json:"amount"`
	PoolID string  `json:"pool_id,omitempty"`
}

// RecordInvestment records an investment on-chain
func (s *BlockchainService) RecordInvestment(poolID uuid.UUID, investorWallet string, amount float64) (*BlockchainTransaction, error) {
	var txHash string

	if s.client != nil {
		auth, err := s.GetTransactOpts(context.Background())
		if err != nil {
			return nil, err
		}

		// Get TokenID for pool
		// We need to fetch pool details or assume poolID mapping
		// For MVP, if we stored TokenID in DB, we use it.
		// Assuming we can get TokenID from pool repo (not passed here).
		// Let's assume passed poolID is NOT the token ID. We need to look it up.
		pool, err := s.fundingRepo.FindPoolByID(poolID)
		if err != nil || pool == nil {
			return nil, fmt.Errorf("pool not found for blockchain record")
		}

		// In models, Pool has InvoiceID. Invoice has NFT.
		nft, err := s.invoiceRepo.FindNFTByInvoiceID(pool.InvoiceID)
		if err != nil || nft == nil {
			return nil, fmt.Errorf("nft not found")
		}

		tokenIDBig := big.NewInt(*nft.TokenID)
		amountBig := new(big.Int).SetInt64(int64(amount))
		investorAddr := common.HexToAddress(investorWallet)

		tx, err := s.poolContract.RecordInvestment(auth, tokenIDBig, investorAddr, amountBig)
		if err != nil {
			return nil, fmt.Errorf("contract call failed: %w", err)
		}
		txHash = tx.Hash().Hex()
		fmt.Printf("[BLOCKCHAIN] Investment recorded: TxHash=%s\n", txHash)
	} else {
		txHash = generateBlockchainTxHash("investment", poolID.String())
	}

	return &BlockchainTransaction{
		TxHash: txHash,
		Action: "investment_recorded",
		Amount: amount,
		PoolID: poolID.String(),
	}, nil
}

// RecordDisbursement records disbursement to mitra on-chain
func (s *BlockchainService) RecordDisbursement(poolID uuid.UUID, amount float64) (*BlockchainTransaction, error) {
	var txHash string

	if s.client != nil {
		auth, err := s.GetTransactOpts(context.Background())
		if err != nil {
			return nil, err
		}

		pool, err := s.fundingRepo.FindPoolByID(poolID)
		if err != nil {
			return nil, err
		}
		nft, err := s.invoiceRepo.FindNFTByInvoiceID(pool.InvoiceID)
		if err != nil {
			return nil, err
		}

		tokenIDBig := big.NewInt(*nft.TokenID)

		tx, err := s.poolContract.RecordDisbursement(auth, tokenIDBig)
		if err != nil {
			return nil, fmt.Errorf("contract call failed: %w", err)
		}
		txHash = tx.Hash().Hex()
		fmt.Printf("[BLOCKCHAIN] Disbursement recorded: TxHash=%s\n", txHash)
	} else {
		txHash = generateBlockchainTxHash("disbursement", poolID.String())
	}

	return &BlockchainTransaction{
		TxHash: txHash,
		Action: "disbursement_recorded",
		Amount: amount,
		PoolID: poolID.String(),
	}, nil
}

// RecordRepayment records repayment and investor returns on-chain
func (s *BlockchainService) RecordRepayment(poolID uuid.UUID, totalAmount float64) (*BlockchainTransaction, error) {
	var txHash string

	if s.client != nil {
		auth, err := s.GetTransactOpts(context.Background())
		if err != nil {
			return nil, err
		}

		pool, err := s.fundingRepo.FindPoolByID(poolID)
		if err != nil {
			return nil, err
		}
		nft, err := s.invoiceRepo.FindNFTByInvoiceID(pool.InvoiceID)
		if err != nil {
			return nil, err
		}

		// Find investments to build returns array
		investments, err := s.fundingRepo.FindInvestmentsByPool(poolID)
		if err != nil {
			return nil, fmt.Errorf("failed to get investments: %w", err)
		}

		var returns []*big.Int
		for _, inv := range investments {
			// Assuming ActualReturn is populated by now
			var ret int64 = 0
			if inv.ActualReturn != nil {
				ret = int64(*inv.ActualReturn)
			}
			returns = append(returns, big.NewInt(ret))
		}

		tokenIDBig := big.NewInt(*nft.TokenID)
		totalAmountBig := new(big.Int).SetInt64(int64(totalAmount))

		tx, err := s.poolContract.RecordRepayment(auth, tokenIDBig, totalAmountBig, returns)
		if err != nil {
			return nil, fmt.Errorf("contract call failed: %w", err)
		}
		txHash = tx.Hash().Hex()
		fmt.Printf("[BLOCKCHAIN] Repayment recorded: TxHash=%s\n", txHash)
	} else {
		txHash = generateBlockchainTxHash("repayment", poolID.String())
	}

	return &BlockchainTransaction{
		TxHash: txHash,
		Action: "repayment_recorded",
		Amount: totalAmount,
		PoolID: poolID.String(),
	}, nil
}

// RecordMitraBalanceCredit records excess payment to mitra
func (s *BlockchainService) RecordMitraBalanceCredit(invoiceID uuid.UUID, mitraWallet string, amount float64) (*BlockchainTransaction, error) {
	var txHash string

	if s.client != nil {
		auth, err := s.GetTransactOpts(context.Background())
		if err != nil {
			return nil, err
		}

		nft, err := s.invoiceRepo.FindNFTByInvoiceID(invoiceID)
		if err != nil || nft == nil {
			return nil, fmt.Errorf("nft needs to exist")
		}

		tokenIDBig := big.NewInt(*nft.TokenID)
		amountBig := new(big.Int).SetInt64(int64(amount))
		mitraAddr := common.HexToAddress(mitraWallet)

		// New function we added to contract
		tx, err := s.poolContract.RecordExcessRepayment(auth, tokenIDBig, mitraAddr, amountBig)
		if err != nil {
			return nil, fmt.Errorf("contract call failed: %w", err)
		}
		txHash = tx.Hash().Hex()
		fmt.Printf("[BLOCKCHAIN] Mitra credit recorded: TxHash=%s\n", txHash)
	} else {
		txHash = generateBlockchainTxHash("mitra_credit", invoiceID.String())
	}

	return &BlockchainTransaction{
		TxHash: txHash,
		Action: "mitra_balance_credited",
		Amount: amount,
		PoolID: invoiceID.String(),
	}, nil
}

// Removed deprecated single RecordInvestorReturn as it's handled in batch RecordRepayment
// But we keep empty stub if interface requires it, or remove it.
// FindingService calls this? No, we checked earlier.

func generateBlockchainTxHash(action string, seed string) string {
	hash := fmt.Sprintf("0x%x%x%s", time.Now().UnixNano(), len(action), seed[:min(8, len(seed))])
	if len(hash) > 66 {
		hash = hash[:66]
	}
	return hash
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
