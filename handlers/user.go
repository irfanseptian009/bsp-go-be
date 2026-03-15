package handlers

import (
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
