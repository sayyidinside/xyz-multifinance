package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	UUID          uuid.UUID       `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	TransactionID uint            `json:"transaction_id" gorm:"index;not null"`
	InstallmentID uint            `json:"installment_id" gorm:"index;not null"`
	Amount        decimal.Decimal `json:"amount" gorm:"type:decimal(20,2);not null"`
	PaymentMethod string          `json:"payment_method" gorm:"size:50;not null"`

	// Relationship
	Transaction Transaction            `json:"transaction" gorm:"foreignKey:TransactionID"`
	Installment TransactionInstallment `json:"installment" gorm:"foreignKey:InstallmentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UUID == uuid.Nil {
		p.UUID = uuid.New()
	}
	return
}
