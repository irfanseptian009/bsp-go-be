package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── Standard JSON Response Helpers ───

// Success sends a 200 JSON response with data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Created sends a 201 JSON response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// ErrorResponse sends an error JSON response.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"message":    message,
	})
}

// ValidationError sends a 400 JSON response for validation errors.
func ValidationError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFound sends a 404 JSON response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource tidak ditemukan"
	}
	ErrorResponse(c, http.StatusNotFound, message)
}

// Unauthorized sends a 401 JSON response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 JSON response.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, http.StatusForbidden, message)
}

// Conflict sends a 409 JSON response.
func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message)
}

// InternalError sends a 500 JSON response.
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, message)
}
