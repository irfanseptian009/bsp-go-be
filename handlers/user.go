package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/services"
	"github.com/irfanseptian/fims-backend/utils"
)

// UserHandler handles user HTTP requests.
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetProfile handles GET /api/users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userId")

	user, err := h.service.FindByID(userID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, user)
}

// UpdateProfile handles PATCH /api/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userId")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	user, err := h.service.Update(userID, req)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, user)
}

// UploadProfilePhoto handles POST /api/users/me/photo
func (h *UserHandler) UploadProfilePhoto(c *gin.Context) {
	userID := c.GetString("userId")

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		fileHeader, err = c.FormFile("file")
	}
	if err != nil {
		fileHeader, err = c.FormFile("avatar")
	}
	if err != nil {
		utils.ValidationError(c, "File foto wajib diisi dengan multipart/form-data (field: photo)")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utils.InternalError(c, "Gagal membaca file foto")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, 5*1024*1024+1))
	if err != nil {
		utils.InternalError(c, "Gagal memproses file foto")
		return
	}

	user, err := h.service.UpdateProfilePhoto(userID, fileHeader.Filename, content)
	if err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	utils.Success(c, user)
}
