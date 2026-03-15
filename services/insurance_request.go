package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// ─── Constants ───

const adminFee = 10000.0

// InsuranceRequestService handles insurance request business logic.
type InsuranceRequestService struct{}

// NewInsuranceRequestService creates a new InsuranceRequestService.
func NewInsuranceRequestService() *InsuranceRequestService {
	return &InsuranceRequestService{}
}

// ─── Customer: Create Request ───

// Create creates a new insurance request with premium calculation.
func (s *InsuranceRequestService) Create(userID string, req dto.CreateInsuranceRequestRequest) (*models.InsuranceRequest, error) {
	db := database.GetDB()

	// Validate occupation type exists
	var occupationType models.OccupationType
	if result := db.First(&occupationType, "id = ?", req.OccupationTypeID); result.Error != nil {
		return nil, errors.New("Tipe okupasi tidak ditemukan")
	}

	// Calculate premium: buildingPrice × premiumRate / 1000 × duration
	basicPremium := (req.BuildingPrice * occupationType.PremiumRate) / 1000 * float64(req.Duration)
	totalAmount := basicPremium + adminFee

	// Generate invoice number
	invoiceNumber, err := s.generateInvoiceNumber()
	if err != nil {
		return nil, errors.New("Gagal membuat nomor invoice")
	}

	// Determine earthquake value
	earthquake := false
	if req.Earthquake != nil {
		earthquake = *req.Earthquake
	}

	request := models.InsuranceRequest{
		UserID:            userID,
		InvoiceNumber:     invoiceNumber,
		OccupationTypeID:  req.OccupationTypeID,
		BuildingPrice:     req.BuildingPrice,
		Duration:          req.Duration,
		ConstructionClass: models.ConstructionClass(req.ConstructionClass),
		Address:           req.Address,
		Province:          req.Province,
		City:              req.City,
		District:          req.District,
		Area:              req.Area,
		Earthquake:        earthquake,
		BasicPremium:      basicPremium,
		AdminFee:          adminFee,
		TotalAmount:       totalAmount,
	}

	if result := db.Create(&request); result.Error != nil {
		return nil, errors.New("Gagal membuat request asuransi")
	}

	// Reload with relation
	db.Preload("OccupationType").First(&request, "id = ?", request.ID)

	return &request, nil
}

// ─── Customer: My Requests ───

// FindByUser returns all insurance requests for a specific user.
func (s *InsuranceRequestService) FindByUser(userID string) ([]models.InsuranceRequest, error) {
	db := database.GetDB()

	var requests []models.InsuranceRequest
	if result := db.Preload("OccupationType").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&requests); result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}

// ─── Admin: List All Requests ───

// FindAll returns all insurance requests with user and occupation type info.
func (s *InsuranceRequestService) FindAll() ([]models.InsuranceRequest, error) {
	db := database.GetDB()

	var requests []models.InsuranceRequest
	if result := db.Preload("User").Preload("OccupationType").
		Order("created_at desc").
		Find(&requests); result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}

// ─── Get Single Request ───

// FindByID returns a single insurance request by ID.
func (s *InsuranceRequestService) FindByID(id string) (*models.InsuranceRequest, error) {
	db := database.GetDB()

	var request models.InsuranceRequest
	if result := db.Preload("User").Preload("OccupationType").
		First(&request, "id = ?", id); result.Error != nil {
		return nil, errors.New("Request tidak ditemukan")
	}

	return &request, nil
}

// FindByInvoiceNumber returns a single insurance request by its invoice number.
func (s *InsuranceRequestService) FindByInvoiceNumber(invoiceNumber string) (*models.InsuranceRequest, error) {
	db := database.GetDB()

	var request models.InsuranceRequest
	if result := db.Preload("User").Preload("OccupationType").
		First(&request, "invoice_number = ?", invoiceNumber); result.Error != nil {
		return nil, errors.New("Request dengan nomor invoice tersebut tidak ditemukan")
	}

	return &request, nil
}

// ─── Admin: Approve Request ───

// Approve changes the request status to APPROVED and generates a policy number.
func (s *InsuranceRequestService) Approve(id string) (*models.InsuranceRequest, error) {
	db := database.GetDB()

	request, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	if request.Status != models.StatusPending {
		return nil, errors.New("Request sudah diproses sebelumnya")
	}

	// Generate policy number
	policyNumber, err := s.generatePolicyNumber()
	if err != nil {
		return nil, errors.New("Gagal membuat nomor polis")
	}

	db.Model(request).Updates(map[string]interface{}{
		"status":        models.StatusApproved,
		"policy_number": policyNumber,
	})

	// Reload with relations
	return s.FindByID(id)
}

// ─── Admin: Reject Request ───

// Reject changes the request status to REJECTED.
func (s *InsuranceRequestService) Reject(id string) (*models.InsuranceRequest, error) {
	db := database.GetDB()

	request, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	if request.Status != models.StatusPending {
		return nil, errors.New("Request sudah diproses sebelumnya")
	}

	db.Model(request).Update("status", models.StatusRejected)

	// Reload with relations
	return s.FindByID(id)
}

// ─── Helpers ───

// generateInvoiceNumber generates a new invoice number: K.001.XXXXX
func (s *InsuranceRequestService) generateInvoiceNumber() (string, error) {
	db := database.GetDB()

	var lastRequest models.InsuranceRequest
	result := db.Order("invoice_number desc").
		Select("invoice_number").
		First(&lastRequest)

	nextNumber := 1

	if result.Error == nil && lastRequest.InvoiceNumber != "" {
		parts := strings.Split(lastRequest.InvoiceNumber, ".")
		if len(parts) > 0 {
			lastNum, err := strconv.Atoi(parts[len(parts)-1])
			if err == nil {
				nextNumber = lastNum + 1
			}
		}
	}

	return fmt.Sprintf("K.001.%05d", nextNumber), nil
}

// generatePolicyNumber generates a new policy number: K.01.001.XXXXX
func (s *InsuranceRequestService) generatePolicyNumber() (string, error) {
	db := database.GetDB()

	var lastApproved models.InsuranceRequest
	result := db.Where("status = ? AND policy_number IS NOT NULL", models.StatusApproved).
		Order("policy_number desc").
		Select("policy_number").
		First(&lastApproved)

	nextNumber := 1

	if result.Error == nil && lastApproved.PolicyNumber != nil {
		parts := strings.Split(*lastApproved.PolicyNumber, ".")
		if len(parts) > 0 {
			lastNum, err := strconv.Atoi(parts[len(parts)-1])
			if err == nil {
				nextNumber = lastNum + 1
			}
		}
	}

	return fmt.Sprintf("K.01.001.%05d", nextNumber), nil
}
