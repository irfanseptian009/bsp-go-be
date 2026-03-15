package services

import (
	"errors"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// UserService handles user business logic.
type UserService struct{}

// NewUserService creates a new UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// FindByID returns a user by their ID (without password).
func (s *UserService) FindByID(id string) (*models.UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if result := db.First(&user, "id = ?", id); result.Error != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	response := user.ToResponse()
	return &response, nil
}

// Update updates the user's profile (name and/or email).
func (s *UserService) Update(id string, req dto.UpdateUserRequest) (*models.UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if result := db.First(&user, "id = ?", id); result.Error != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// Build updates map
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if len(updates) > 0 {
		if result := db.Model(&user).Updates(updates); result.Error != nil {
			return nil, errors.New("Gagal mengupdate profile")
		}
	}

	// Reload user
	db.First(&user, "id = ?", id)
	response := user.ToResponse()
	return &response, nil
}
