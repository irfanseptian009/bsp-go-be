package models

import (
	"time"

	"gorm.io/gorm"
)

// ─── Policy Model ───

type Policy struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PolicyNumber      string    `gorm:"column:policy_number;uniqueIndex;not null" json:"policyNumber"`
	ApplicationNumber string    `gorm:"column:application_number;uniqueIndex;not null" json:"applicationNumber"`
	Name              string    `gorm:"not null" json:"name"`
	BranchID          string    `gorm:"column:branch_id;not null" json:"branchId"`
	BirthDate         time.Time `gorm:"column:birth_date;not null" json:"birthDate"`
	Duration          int       `gorm:"not null" json:"duration"`
	BuildingPrice     float64   `gorm:"column:building_price;not null" json:"buildingPrice"`
	OccupationTypeID  string    `gorm:"column:occupation_type_id;not null" json:"occupationTypeId"`
	Premium           float64   `gorm:"not null" json:"premium"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	Branch         Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	OccupationType OccupationType `gorm:"foreignKey:OccupationTypeID" json:"occupationType,omitempty"`
}

func (Policy) TableName() string {
	return "policies"
}

func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = newUUID()
	}

	return nil
}
