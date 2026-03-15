package dto

// ─── User DTOs ───

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" binding:"omitempty,email"`
}
