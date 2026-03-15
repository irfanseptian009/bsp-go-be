package dto

// ─── Branch DTOs ───

type CreateBranchRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type UpdateBranchRequest struct {
	Code *string `json:"code"`
	Name *string `json:"name"`
}
