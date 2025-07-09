package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	UUID             uuid.UUID       `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	TransactionID    uint            `json:"transaction_id"`
	Amount           decimal.Decimal `json:"amount" gorm:"type:decimal(20,2)"`
	InstalmentNumber uint            `json:"instalment_number" gorm:"type:smallint unsigned"`
	Transaction      Transaction     `json:"transaction" gorm:"foreignKey:TransactionID"`
	gorm.Model
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
