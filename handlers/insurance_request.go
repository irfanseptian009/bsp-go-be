package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/services"
	"github.com/irfanseptian/fims-backend/utils"
)

// InsuranceRequestHandler handles insurance request HTTP requests.
type InsuranceRequestHandler struct {
	service *services.InsuranceRequestService
}

// NewInsuranceRequestHandler creates a new InsuranceRequestHandler.
func NewInsuranceRequestHandler(service *services.InsuranceRequestService) *InsuranceRequestHandler {
	return &InsuranceRequestHandler{service: service}
}

// Create handles POST /api/insurance-requests (Customer only)
func (h *InsuranceRequestHandler) Create(c *gin.Context) {
	userID := c.GetString("userId")

	var req dto.CreateInsuranceRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	result, err := h.service.Create(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Created(c, result)
}

// FindMyRequests handles GET /api/insurance-requests/my-requests (Customer only)
func (h *InsuranceRequestHandler) FindMyRequests(c *gin.Context) {
	userID := c.GetString("userId")

	requests, err := h.service.FindByUser(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, requests)
}

// FindAll handles GET /api/insurance-requests (Admin only)
func (h *InsuranceRequestHandler) FindAll(c *gin.Context) {
	requests, err := h.service.FindAll()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, requests)
}

// FindByID handles GET /api/insurance-requests/:id
func (h *InsuranceRequestHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	request, err := h.service.FindByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, request)
}

// FindByInvoiceNumber handles GET /api/insurance-requests/invoice/:invoiceNumber
func (h *InsuranceRequestHandler) FindByInvoiceNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoiceNumber")

	request, err := h.service.FindByInvoiceNumber(invoiceNumber)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, request)
}

// Approve handles PATCH /api/insurance-requests/:id/approve (Admin only)
func (h *InsuranceRequestHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.Approve(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, result)
}

// Reject handles PATCH /api/insurance-requests/:id/reject (Admin only)
func (h *InsuranceRequestHandler) Reject(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.Reject(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, result)
}
