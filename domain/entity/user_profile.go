package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID         uint                `json:"id" gorm:"primaryKey"`
	UUID       uuid.UUID           `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	UserID     uint                `json:"user_id" gorm:"not null"`
	Name       string              `json:"name" gorm:"not null"`
	LegalName  string              `json:"legal_name" gorm:"not null"`
	Nik        string              `json:"nik" gorm:"unique;not null;uniqueIndex;type:varchar(16)"`
	BirthPlace string              `json:"birth_place" gorm:"not null"`
	BirthDate  time.Time           `json:"birth_date" gorm:"type:date;not null"`
	Salary     decimal.NullDecimal `json:"salary" gorm:"type:decimal(20,2);not null"`
	User       User                `json:"user" gorm:"foreignKey:UserID"`
	gorm.Model
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (up *UserProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if up.UUID == uuid.Nil {
		up.UUID = uuid.New()
	}
	return
}
