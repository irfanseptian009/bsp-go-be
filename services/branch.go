package services

import (
	"errors"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// BranchService handles branch business logic.
type BranchService struct{}

// NewBranchService creates a new BranchService.
func NewBranchService() *BranchService {
	return &BranchService{}
}

// FindAll returns all branches ordered by name.
func (s *BranchService) FindAll() ([]models.Branch, error) {
	db := database.GetDB()

	var branches []models.Branch
	if result := db.Order("name asc").Find(&branches); result.Error != nil {
		return nil, result.Error
	}

	return branches, nil
}

// FindByID returns a branch by ID.
func (s *BranchService) FindByID(id string) (*models.Branch, error) {
	db := database.GetDB()

	var branch models.Branch
	if result := db.First(&branch, "id = ?", id); result.Error != nil {
		return nil, errors.New("Cabang tidak ditemukan")
	}

	return &branch, nil
}

// Create creates a new branch.
func (s *BranchService) Create(req dto.CreateBranchRequest) (*models.Branch, error) {
	db := database.GetDB()

	branch := models.Branch{
		Code: req.Code,
		Name: req.Name,
	}

	if result := db.Create(&branch); result.Error != nil {
		return nil, errors.New("Gagal membuat cabang")
	}

	return &branch, nil
}

// Update updates an existing branch.
func (s *BranchService) Update(id string, req dto.UpdateBranchRequest) (*models.Branch, error) {
	db := database.GetDB()

	branch, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if len(updates) > 0 {
		if result := db.Model(branch).Updates(updates); result.Error != nil {
			return nil, errors.New("Gagal mengupdate cabang")
		}
	}

	return s.FindByID(id)
}

// Delete removes a branch by ID.
func (s *BranchService) Delete(id string) error {
	db := database.GetDB()

	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if result := db.Delete(&models.Branch{}, "id = ?", id); result.Error != nil {
		return errors.New("Gagal menghapus cabang")
	}

	return nil
}
