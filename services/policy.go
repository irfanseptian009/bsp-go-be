package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// PolicyService handles policy business logic.
type PolicyService struct{}

// NewPolicyService creates a new PolicyService.
func NewPolicyService() *PolicyService {
	return &PolicyService{}
}

// ─── Create Policy ───

// Create creates a new policy with auto-generated numbers and premium calculation.
func (s *PolicyService) Create(req dto.CreatePolicyRequest) (*models.Policy, error) {
	db := database.GetDB()

	// Validate occupation type
	var occupationType models.OccupationType
	if result := db.First(&occupationType, "id = ?", req.OccupationTypeID); result.Error != nil {
		return nil, errors.New("Tipe okupasi tidak ditemukan")
	}

	// Validate branch
	var branch models.Branch
	if result := db.First(&branch, "id = ?", req.BranchID); result.Error != nil {
		return nil, errors.New("Cabang tidak ditemukan")
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		birthDate, err = time.Parse(time.RFC3339, req.BirthDate)
		if err != nil {
			return nil, errors.New("Format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
		}
	}

	// Calculate premium: buildingPrice × premiumRate / 1000 × duration + 2500
	premium := (req.BuildingPrice*occupationType.PremiumRate)/1000*float64(req.Duration) + 2500

	// Generate numbers
	policyNumber, err := s.generatePolicyNumber()
	if err != nil {
		return nil, errors.New("Gagal membuat nomor polis")
	}

	applicationNumber, err := s.generateApplicationNumber()
	if err != nil {
		return nil, errors.New("Gagal membuat nomor aplikasi")
	}

	policy := models.Policy{
		PolicyNumber:      policyNumber,
		ApplicationNumber: applicationNumber,
		Name:              req.Name,
		BranchID:          req.BranchID,
		BirthDate:         birthDate,
		Duration:          req.Duration,
		BuildingPrice:     req.BuildingPrice,
		OccupationTypeID:  req.OccupationTypeID,
		Premium:           premium,
	}

	if result := db.Create(&policy); result.Error != nil {
		return nil, errors.New("Gagal membuat polis")
	}

	// Reload with relations
	db.Preload("Branch").Preload("OccupationType").First(&policy, "id = ?", policy.ID)

	return &policy, nil
}

// ─── Find All / Search ───

// FindAll returns all policies, optionally filtered by name, branch, or occupation type.
func (s *PolicyService) FindAll(query dto.SearchPolicyQuery) ([]models.Policy, error) {
	db := database.GetDB()

	q := db.Preload("Branch").Preload("OccupationType")

	if query.Name != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?)", "%"+query.Name+"%")
	}
	if query.BranchID != "" {
		q = q.Where("branch_id = ?", query.BranchID)
	}
	if query.OccupationTypeID != "" {
		q = q.Where("occupation_type_id = ?", query.OccupationTypeID)
	}

	var policies []models.Policy
	if result := q.Order("created_at desc").Find(&policies); result.Error != nil {
		return nil, result.Error
	}

	return policies, nil
}

// ─── Find By ID ───

// FindByID returns a single policy by ID.
func (s *PolicyService) FindByID(id string) (*models.Policy, error) {
	db := database.GetDB()

	var policy models.Policy
	if result := db.Preload("Branch").Preload("OccupationType").
		First(&policy, "id = ?", id); result.Error != nil {
		return nil, errors.New("Polis tidak ditemukan")
	}

	return &policy, nil
}

// ─── Update ───

// Update updates an existing policy, recalculating premium if needed.
func (s *PolicyService) Update(id string, req dto.UpdatePolicyRequest) (*models.Policy, error) {
	db := database.GetDB()

	policy, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.BranchID != nil {
		updates["branch_id"] = *req.BranchID
	}
	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			birthDate, err = time.Parse(time.RFC3339, *req.BirthDate)
			if err != nil {
				return nil, errors.New("Format tanggal lahir tidak valid")
			}
		}
		updates["birth_date"] = birthDate
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.BuildingPrice != nil {
		updates["building_price"] = *req.BuildingPrice
	}
	if req.OccupationTypeID != nil {
		updates["occupation_type_id"] = *req.OccupationTypeID
	}

	// Recalculate premium if relevant fields are changed
	if req.BuildingPrice != nil || req.OccupationTypeID != nil || req.Duration != nil {
		occupationTypeID := policy.OccupationTypeID
		if req.OccupationTypeID != nil {
			occupationTypeID = *req.OccupationTypeID
		}

		var occupationType models.OccupationType
		if result := db.First(&occupationType, "id = ?", occupationTypeID); result.Error != nil {
			return nil, errors.New("Tipe okupasi tidak ditemukan")
		}

		buildingPrice := policy.BuildingPrice
		if req.BuildingPrice != nil {
			buildingPrice = *req.BuildingPrice
		}

		duration := policy.Duration
		if req.Duration != nil {
			duration = *req.Duration
		}

		premium := (buildingPrice*occupationType.PremiumRate)/1000*float64(duration) + 2500
		updates["premium"] = premium
	}

	if len(updates) > 0 {
		if result := db.Model(policy).Updates(updates); result.Error != nil {
			return nil, errors.New("Gagal mengupdate polis")
		}
	}

	return s.FindByID(id)
}

// ─── Delete ───

// Delete removes a policy by ID.
func (s *PolicyService) Delete(id string) error {
	db := database.GetDB()

	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if result := db.Delete(&models.Policy{}, "id = ?", id); result.Error != nil {
		return errors.New("Gagal menghapus polis")
	}

	return nil
}

// ─── Helpers ───

// generatePolicyNumber generates: 001.{year}.XXXXX
func (s *PolicyService) generatePolicyNumber() (string, error) {
	db := database.GetDB()

	var lastPolicy models.Policy
	result := db.Order("policy_number desc").
		Select("policy_number").
		First(&lastPolicy)

	nextNumber := 1

	if result.Error == nil && lastPolicy.PolicyNumber != "" {
		parts := strings.Split(lastPolicy.PolicyNumber, ".")
		if len(parts) > 0 {
			lastNum, err := strconv.Atoi(parts[len(parts)-1])
			if err == nil {
				nextNumber = lastNum + 1
			}
		}
	}

	year := time.Now().Year()
	return fmt.Sprintf("001.%d.%05d", year, nextNumber), nil
}

// generateApplicationNumber generates: 00001{year}XXXXXX
func (s *PolicyService) generateApplicationNumber() (string, error) {
	db := database.GetDB()

	var lastPolicy models.Policy
	result := db.Order("application_number desc").
		Select("application_number").
		First(&lastPolicy)

	nextNumber := 1

	if result.Error == nil && lastPolicy.ApplicationNumber != "" {
		appNum := lastPolicy.ApplicationNumber
		if len(appNum) >= 6 {
			lastSix := appNum[len(appNum)-6:]
			lastNum, err := strconv.Atoi(lastSix)
			if err == nil {
				nextNumber = lastNum + 1
			}
		}
	}

	year := time.Now().Year()
	return fmt.Sprintf("00001%d%06d", year, nextNumber), nil
}
