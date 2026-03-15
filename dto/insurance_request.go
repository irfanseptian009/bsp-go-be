package dto

// ─── InsuranceRequest DTOs ───

type CreateInsuranceRequestRequest struct {
	OccupationTypeID  string  `json:"occupationTypeId" binding:"required"`
	BuildingPrice     float64 `json:"buildingPrice" binding:"required,gt=0"`
	Duration          int     `json:"duration" binding:"required,min=1,max=10"`
	ConstructionClass string  `json:"constructionClass" binding:"required,oneof=KELAS_1 KELAS_2 KELAS_3"`
	Address           string  `json:"address" binding:"required"`
	Province          string  `json:"province" binding:"required"`
	City              string  `json:"city" binding:"required"`
	District          string  `json:"district" binding:"required"`
	Area              string  `json:"area" binding:"required"`
	Earthquake        *bool   `json:"earthquake"`
}
