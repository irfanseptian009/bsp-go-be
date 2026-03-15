package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Enums ───

type RequestStatus string

const (
	StatusPending  RequestStatus = "PENDING"
	StatusApproved RequestStatus = "APPROVED"
	StatusRejected RequestStatus = "REJECTED"
)

type ConstructionClass string

const (
	ConstructionKelas1 ConstructionClass = "KELAS_1"
	ConstructionKelas2 ConstructionClass = "KELAS_2"
	ConstructionKelas3 ConstructionClass = "KELAS_3"
)

// ─── InsuranceRequest Model ───

type InsuranceRequest struct {
	ID                string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID            string            `gorm:"column:user_id;not null" json:"userId"`
	InvoiceNumber     string            `gorm:"column:invoice_number;uniqueIndex;not null" json:"invoiceNumber"`
	OccupationTypeID  string            `gorm:"column:occupation_type_id;not null" json:"occupationTypeId"`
	BuildingPrice     float64           `gorm:"column:building_price;not null" json:"buildingPrice"`
	Duration          int               `gorm:"not null" json:"duration"`
	ConstructionClass ConstructionClass `gorm:"column:construction_class;type:varchar(20);not null" json:"constructionClass"`
	Address           string            `gorm:"not null" json:"address"`
	Province          string            `gorm:"not null" json:"province"`
	City              string            `gorm:"not null" json:"city"`
	District          string            `gorm:"not null" json:"district"`
	Area              string            `gorm:"not null" json:"area"`
	Earthquake        bool              `gorm:"default:false" json:"earthquake"`
	BasicPremium      float64           `gorm:"column:basic_premium;not null" json:"basicPremium"`
	AdminFee          float64           `gorm:"column:admin_fee;default:10000" json:"adminFee"`
	TotalAmount       float64           `gorm:"column:total_amount;not null" json:"totalAmount"`
	Status            RequestStatus     `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	PolicyNumber      *string           `gorm:"column:policy_number" json:"policyNumber"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OccupationType OccupationType `gorm:"foreignKey:OccupationTypeID" json:"occupationType,omitempty"`
}

func (InsuranceRequest) TableName() string {
	return "insurance_requests"
}

// BeforeCreate ensures UUID is always set even when DB default is missing.
func (r *InsuranceRequest) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}

	return nil
}
