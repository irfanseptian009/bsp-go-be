package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/services"
	"github.com/irfanseptian/fims-backend/utils"
)

// PolicyHandler handles policy HTTP requests.
type PolicyHandler struct {
	service *services.PolicyService
}

// NewPolicyHandler creates a new PolicyHandler.
func NewPolicyHandler(service *services.PolicyService) *PolicyHandler {
	return &PolicyHandler{service: service}
}

// FindAll handles GET /api/policies
func (h *PolicyHandler) FindAll(c *gin.Context) {
	var query dto.SearchPolicyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	policies, err := h.service.FindAll(query)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, policies)
}

// FindByID handles GET /api/policies/:id
func (h *PolicyHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	policy, err := h.service.FindByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, policy)
}

// Create handles POST /api/policies (Admin only)
func (h *PolicyHandler) Create(c *gin.Context) {
	var req dto.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	policy, err := h.service.Create(req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, policy)
}

// Update handles PATCH /api/policies/:id (Admin only)
func (h *PolicyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	policy, err := h.service.Update(id, req)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, policy)
}

// Delete handles DELETE /api/policies/:id (Admin only)
func (h *PolicyHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Polis berhasil dihapus"})
}
