package services

import (
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

// GradingService handles invoice grading based on risk matrix
type GradingService struct {
	invoiceRepo repository.InvoiceRepositoryInterface
	buyerRepo   repository.BuyerRepositoryInterface
}

func NewGradingService(invoiceRepo repository.InvoiceRepositoryInterface, buyerRepo repository.BuyerRepositoryInterface) *GradingService {
	return &GradingService{
		invoiceRepo: invoiceRepo,
		buyerRepo:   buyerRepo,
	}
}

// CountryRiskTier maps countries to risk tiers
// Tier 1 (Low Risk): Developed countries with strong payment history
// Tier 2 (Medium Risk): Emerging markets with moderate risk
// Tier 3 (High Risk): Countries with high default risk
var CountryRiskTier = map[string]int{
	// Tier 1 - Low Risk
	"USA": 1, "DEU": 1, "JPN": 1, "GBR": 1, "FRA": 1, "CHE": 1, "NLD": 1, "AUS": 1, "CAN": 1, "SGP": 1, "KOR": 1,
	"United States": 1, "Germany": 1, "Japan": 1, "United Kingdom": 1, "France": 1, "Switzerland": 1,
	"Netherlands": 1, "Australia": 1, "Canada": 1, "Singapore": 1, "South Korea": 1,
	// Tier 2 - Medium Risk
	"CHN": 2, "IND": 2, "BRA": 2, "MEX": 2, "THA": 2, "MYS": 2, "VNM": 2, "PHL": 2, "IDN": 2, "TUR": 2, "SAU": 2, "ARE": 2,
	"China": 2, "India": 2, "Brazil": 2, "Mexico": 2, "Thailand": 2, "Malaysia": 2, "Vietnam": 2,
	"Philippines": 2, "Indonesia": 2, "Turkey": 2, "Saudi Arabia": 2, "UAE": 2,
	// Tier 3 - High Risk
	"NGA": 3, "PAK": 3, "BGD": 3, "EGY": 3, "KEN": 3, "ZAF": 3, "ARG": 3,
	"Nigeria": 3, "Pakistan": 3, "Bangladesh": 3, "Egypt": 3, "Kenya": 3, "South Africa": 3, "Argentina": 3,
}

// GradeInvoice calculates grade based on risk matrix
// Returns: grade (A/B/C), score (0-100), country risk (low/medium/high)
func (s *GradingService) GradeInvoice(invoice *models.Invoice, buyer *models.Buyer, exporterInvoiceCount int, documentCount int) (string, int, string) {
	score := 0

	// 1. Buyer Country Risk (40 points max)
	countryRisk := "medium"
	countryTier := 2 // Default to medium risk
	if tier, ok := CountryRiskTier[buyer.Country]; ok {
		countryTier = tier
	}

	switch countryTier {
	case 1:
		score += 40
		countryRisk = "low"
	case 2:
		score += 25
		countryRisk = "medium"
	case 3:
		score += 10
		countryRisk = "high"
	}

	// 2. Exporter History (30 points max)
	// Frequent exporter (5+ invoices) = A
	// Some history (1-4 invoices) = B
	// First time = C
	if exporterInvoiceCount >= 5 {
		score += 30
	} else if exporterInvoiceCount >= 1 {
		score += 20
	} else {
		score += 10
	}

	// 3. Document Completeness (30 points max)
	// Full documents (3+ types) = A
	// Partial (1-2 types) = B
	// No documents = C
	if documentCount >= 3 {
		score += 30
	} else if documentCount >= 1 {
		score += 20
	} else {
		score += 5
	}

	// Determine grade based on score
	var grade string
	if score >= 80 {
		grade = "A"
	} else if score >= 50 {
		grade = "B"
	} else {
		grade = "C"
	}

	return grade, score, countryRisk
}

// CalculateDocumentScore calculates document completeness score
func (s *GradingService) CalculateDocumentScore(documents []models.InvoiceDocument) int {
	if len(documents) == 0 {
		return 0
	}

	// Check for key document types
	hasInvoice := false
	hasBOL := false
	hasPO := false
	hasInsurance := false

	for _, doc := range documents {
		switch doc.DocumentType {
		case models.DocTypeInvoicePDF, models.DocTypeCommercialInvoice:
			hasInvoice = true
		case models.DocTypeBillOfLading:
			hasBOL = true
		case models.DocTypePurchaseOrder:
			hasPO = true
		case models.DocTypeInsurance:
			hasInsurance = true
		}
	}

	score := 0
	if hasInvoice {
		score += 30
	}
	if hasBOL {
		score += 30
	}
	if hasPO {
		score += 20
	}
	if hasInsurance {
		score += 20
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// CheckInsurance checks if invoice has insurance document
func (s *GradingService) CheckInsurance(documents []models.InvoiceDocument) bool {
	for _, doc := range documents {
		if doc.DocumentType == models.DocTypeInsurance {
			return true
		}
	}
	return false
}

// GetExporterInvoiceCount returns the number of completed invoices for an exporter
func (s *GradingService) GetExporterInvoiceCount(exporterID uuid.UUID) (int, error) {
	filter := &models.InvoiceFilter{
		Page:    1,
		PerPage: 1000,
	}
	invoices, _, err := s.invoiceRepo.FindByExporter(exporterID, filter)
	if err != nil {
		return 0, err
	}

	// Count only repaid invoices
	count := 0
	for _, inv := range invoices {
		if inv.Status == models.StatusRepaid {
			count++
		}
	}
	return count, nil
}

// GradeResult represents the grading result
type GradeResult struct {
	Grade             string `json:"grade"`
	Score             int    `json:"score"`
	CountryRisk       string `json:"country_risk"`
	IsInsured         bool   `json:"is_insured"`
	DocumentScore     int    `json:"document_score"`
	ExporterHistory   int    `json:"exporter_history"`
	IsRepeatBuyer     bool   `json:"is_repeat_buyer"`
}

// FullGrade performs full grading analysis on an invoice
func (s *GradingService) FullGrade(invoiceID uuid.UUID) (*GradeResult, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}

	buyer, err := s.buyerRepo.FindByID(invoice.BuyerID)
	if err != nil {
		return nil, err
	}

	documents, err := s.invoiceRepo.FindDocumentsByInvoiceID(invoiceID)
	if err != nil {
		return nil, err
	}

	exporterCount, err := s.GetExporterInvoiceCount(invoice.ExporterID)
	if err != nil {
		exporterCount = 0
	}

	grade, score, countryRisk := s.GradeInvoice(invoice, buyer, exporterCount, len(documents))
	docScore := s.CalculateDocumentScore(documents)
	isInsured := s.CheckInsurance(documents)

	return &GradeResult{
		Grade:           grade,
		Score:           score,
		CountryRisk:     countryRisk,
		IsInsured:       isInsured,
		DocumentScore:   docScore,
		ExporterHistory: exporterCount,
		IsRepeatBuyer:   buyer.TotalInvoices > 0,
	}, nil
}
