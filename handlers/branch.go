package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/services"
	"github.com/irfanseptian/fims-backend/utils"
)

// BranchHandler handles branch HTTP requests.
type BranchHandler struct {
	service *services.BranchService
}

// NewBranchHandler creates a new BranchHandler.
func NewBranchHandler(service *services.BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

// FindAll handles GET /api/branches
func (h *BranchHandler) FindAll(c *gin.Context) {
	branches, err := h.service.FindAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, branches)
}

// FindByID handles GET /api/branches/:id
func (h *BranchHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	branch, err := h.service.FindByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, branch)
}

// Create handles POST /api/branches
func (h *BranchHandler) Create(c *gin.Context) {
	var req dto.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	branch, err := h.service.Create(req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, branch)
}

// Update handles PATCH /api/branches/:id
func (h *BranchHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	branch, err := h.service.Update(id, req)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, branch)
}

// Delete handles DELETE /api/branches/:id
func (h *BranchHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Cabang berhasil dihapus"})
}
