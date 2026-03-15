package models

import "gorm.io/gorm"

// ─── Branch Model ───

type Branch struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code string `gorm:"uniqueIndex;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`

	Policies []Policy `gorm:"foreignKey:BranchID" json:"policies,omitempty"`
}

func (Branch) TableName() string {
	return "branches"
}

func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = newUUID()
	}

	return nil
}
