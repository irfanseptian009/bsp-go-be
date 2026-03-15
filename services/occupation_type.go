package services

import (
	"errors"
	"strings"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// OccupationTypeService handles occupation type business logic.
type OccupationTypeService struct{}

// NewOccupationTypeService creates a new OccupationTypeService.
func NewOccupationTypeService() *OccupationTypeService {
	return &OccupationTypeService{}
}

// FindAll returns all occupation types ordered by name.
func (s *OccupationTypeService) FindAll() ([]models.OccupationType, error) {
	db := database.GetDB()

	var types []models.OccupationType
	if result := db.Order("name asc").Find(&types); result.Error != nil {
		return nil, result.Error
	}

	return types, nil
}

// FindByID returns an occupation type by ID.
func (s *OccupationTypeService) FindByID(id string) (*models.OccupationType, error) {
	db := database.GetDB()

	var occupationType models.OccupationType
	if result := db.First(&occupationType, "id = ?", id); result.Error != nil {
		return nil, errors.New("Tipe okupasi tidak ditemukan")
	}

	return &occupationType, nil
}

// Create creates a new occupation type.
func (s *OccupationTypeService) Create(req dto.CreateOccupationTypeRequest) (*models.OccupationType, error) {
	db := database.GetDB()

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("Kode dan nama wajib diisi")
	}

	var existing models.OccupationType
	if result := db.Where("code = ?", code).First(&existing); result.Error == nil {
		return nil, errors.New("Kode okupasi sudah digunakan")
	}

	occupationType := models.OccupationType{
		Code:        code,
		Name:        name,
		PremiumRate: req.PremiumRate,
	}

	if result := db.Create(&occupationType); result.Error != nil {
		if strings.Contains(strings.ToLower(result.Error.Error()), "duplicate key") {
			return nil, errors.New("Kode okupasi sudah digunakan")
		}
		return nil, errors.New("Gagal membuat tipe okupasi")
	}

	return &occupationType, nil
}

// Update updates an existing occupation type.
func (s *OccupationTypeService) Update(id string, req dto.UpdateOccupationTypeRequest) (*models.OccupationType, error) {
	db := database.GetDB()

	occupationType, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, errors.New("Kode tidak boleh kosong")
		}
		updates["code"] = code
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("Nama tidak boleh kosong")
		}
		updates["name"] = name
	}
	if req.PremiumRate != nil {
		updates["premium_rate"] = *req.PremiumRate
	}

	if len(updates) > 0 {
		if result := db.Model(occupationType).Updates(updates); result.Error != nil {
			if strings.Contains(strings.ToLower(result.Error.Error()), "duplicate key") {
				return nil, errors.New("Kode okupasi sudah digunakan")
			}
			return nil, errors.New("Gagal mengupdate tipe okupasi")
		}
	}

	return s.FindByID(id)
}

// Delete removes an occupation type by ID.
func (s *OccupationTypeService) Delete(id string) error {
	db := database.GetDB()

	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if result := db.Delete(&models.OccupationType{}, "id = ?", id); result.Error != nil {
		return errors.New("Gagal menghapus tipe okupasi")
	}

	return nil
}
