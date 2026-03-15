package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/services"
	"github.com/irfanseptian/fims-backend/utils"
)

// OccupationTypeHandler handles occupation type HTTP requests.
type OccupationTypeHandler struct {
	service *services.OccupationTypeService
}

// NewOccupationTypeHandler creates a new OccupationTypeHandler.
func NewOccupationTypeHandler(service *services.OccupationTypeService) *OccupationTypeHandler {
	return &OccupationTypeHandler{service: service}
}

// FindAll handles GET /api/occupation-types
// @Summary List occupation types
// @Tags Occupation Types
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.OccupationType
// @Failure 401 {object} map[string]interface{}
// @Router /occupation-types [get]
func (h *OccupationTypeHandler) FindAll(c *gin.Context) {
	types, err := h.service.FindAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, types)
}

// FindByID handles GET /api/occupation-types/:id
// @Summary Get occupation type by ID
// @Tags Occupation Types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Occupation Type ID"
// @Success 200 {object} models.OccupationType
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /occupation-types/{id} [get]
func (h *OccupationTypeHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	occupationType, err := h.service.FindByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, occupationType)
}

// Create handles POST /api/occupation-types
// @Summary Create occupation type
// @Tags Occupation Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.CreateOccupationTypeRequest true "Create occupation type payload"
// @Success 201 {object} models.OccupationType
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /occupation-types [post]
func (h *OccupationTypeHandler) Create(c *gin.Context) {
	var req dto.CreateOccupationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	occupationType, err := h.service.Create(req)
	if err != nil {
		if err.Error() == "Kode okupasi sudah digunakan" {
			utils.Conflict(c, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Created(c, occupationType)
}

// Update handles PATCH /api/occupation-types/:id
// @Summary Update occupation type
// @Tags Occupation Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Occupation Type ID"
// @Param payload body dto.UpdateOccupationTypeRequest true "Update occupation type payload"
// @Success 200 {object} models.OccupationType
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /occupation-types/{id} [patch]
func (h *OccupationTypeHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateOccupationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	occupationType, err := h.service.Update(id, req)
	if err != nil {
		if err.Error() == "Tipe okupasi tidak ditemukan" {
			utils.NotFound(c, err.Error())
			return
		}

		if err.Error() == "Kode okupasi sudah digunakan" {
			utils.Conflict(c, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, occupationType)
}

// Delete handles DELETE /api/occupation-types/:id
// @Summary Delete occupation type
// @Tags Occupation Types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Occupation Type ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /occupation-types/{id} [delete]
func (h *OccupationTypeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "Tipe okupasi tidak ditemukan" {
			utils.NotFound(c, err.Error())
			return
		}

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Tipe okupasi berhasil dihapus"})
}
