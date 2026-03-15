package models

import "gorm.io/gorm"

// ─── OccupationType Model ───

type OccupationType struct {
	ID          string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code        string  `gorm:"uniqueIndex;not null" json:"code"`
	Name        string  `gorm:"not null" json:"name"`
	PremiumRate float64 `gorm:"column:premium_rate;not null" json:"premiumRate"`

	InsuranceRequests []InsuranceRequest `gorm:"foreignKey:OccupationTypeID" json:"insuranceRequests,omitempty"`
	Policies          []Policy           `gorm:"foreignKey:OccupationTypeID" json:"policies,omitempty"`
}

func (OccupationType) TableName() string {
	return "occupation_types"
}

func (o *OccupationType) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = newUUID()
	}

	return nil
}
