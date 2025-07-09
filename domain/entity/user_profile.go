package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserProfile struct {
	ID         uint                `json:"id" gorm:"primaryKey"`
	UserID     uint                `json:"user_id" gorm:"not null;uniqueIndex"`
	Name       string              `json:"name" gorm:"not null"`
	LegalName  string              `json:"legal_name" gorm:"not null"`
	Nik        string              `json:"nik" gorm:"unique;not null;uniqueIndex;type:varchar(16)"`
	BirthPlace string              `json:"birth_place" gorm:"not null"`
	BirthDate  time.Time           `json:"birth_date" gorm:"type:date;not null"`
	Salary     decimal.NullDecimal `json:"salary" gorm:"type:decimal(20,2);not null"`

	// Relationship
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
