package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Limit struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	UserID        uint            `json:"user_id" gorm:"not null;index:idx_limit,unique"`
	Tenor         uint            `json:"tenor" gorm:"type:smallint unsigned;not null;index:idx_limit,unique"`
	CurrentLimit  decimal.Decimal `json:"current_limit" gorm:"type:decimal(20,2);not null"`
	OriginalLimit decimal.Decimal `json:"original_limit" gorm:"type:decimal(20,2);not null"`

	// Relationship
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Limit) TableName() string {
	return "limits"
}
