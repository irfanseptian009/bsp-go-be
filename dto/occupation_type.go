package dto

// ─── OccupationType DTOs ───

type CreateOccupationTypeRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	PremiumRate float64 `json:"premiumRate" binding:"required,gt=0"`
}

type UpdateOccupationTypeRequest struct {
	Code        *string  `json:"code"`
	Name        *string  `json:"name"`
	PremiumRate *float64 `json:"premiumRate" binding:"omitempty,gt=0"`
}
