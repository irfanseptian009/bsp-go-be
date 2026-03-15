package dto

// ─── Policy DTOs ───

type CreatePolicyRequest struct {
	Name             string  `json:"name" binding:"required"`
	BranchID         string  `json:"branchId" binding:"required"`
	BirthDate        string  `json:"birthDate" binding:"required"`
	Duration         int     `json:"duration" binding:"required,min=1,max=10"`
	BuildingPrice    float64 `json:"buildingPrice" binding:"required,gt=0"`
	OccupationTypeID string  `json:"occupationTypeId" binding:"required"`
}

type UpdatePolicyRequest struct {
	Name             *string  `json:"name"`
	BranchID         *string  `json:"branchId"`
	BirthDate        *string  `json:"birthDate"`
	Duration         *int     `json:"duration" binding:"omitempty,min=1,max=10"`
	BuildingPrice    *float64 `json:"buildingPrice" binding:"omitempty,gt=0"`
	OccupationTypeID *string  `json:"occupationTypeId"`
}

type SearchPolicyQuery struct {
	Name             string `form:"name"`
	BranchID         string `form:"branchId"`
	OccupationTypeID string `form:"occupationTypeId"`
}
