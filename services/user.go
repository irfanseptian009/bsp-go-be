package services

import (
	"errors"
	"fmt"

	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/models"
)

// UserService handles user business logic.
type UserService struct {
	storageService *StorageService
}

// NewUserService creates a new UserService.
func NewUserService(storageService *StorageService) *UserService {
	return &UserService{storageService: storageService}
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

// UpdateProfilePhoto uploads profile photo to storage and updates user photo URL.
func (s *UserService) UpdateProfilePhoto(id string, fileName string, content []byte) (*models.UserResponse, error) {
	if len(content) == 0 {
		return nil, errors.New("File foto tidak boleh kosong")
	}

	if len(content) > 5*1024*1024 {
		return nil, errors.New("Ukuran foto maksimal 5MB")
	}

	if s.storageService == nil || !s.storageService.IsReady() {
		return nil, errors.New("Supabase Storage belum dikonfigurasi")
	}

	db := database.GetDB()

	var user models.User
	if result := db.First(&user, "id = ?", id); result.Error != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	photoURL, err := s.storageService.UploadProfilePhoto(id, fileName, content)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengunggah foto profil: %w", err)
	}

	updates := map[string]interface{}{
		"profile_photo_url": photoURL,
	}

	if result := db.Model(&user).Updates(updates); result.Error != nil {
		return nil, errors.New("Gagal menyimpan foto profil")
	}

	db.First(&user, "id = ?", id)
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
