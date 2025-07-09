package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	ID           uint                `json:"id" gorm:"primaryKey"`
	UUID         uuid.UUID           `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	RoleID       uint                `json:"role_id"`
	Username     string              `json:"username" gorm:"index"`
	Name         string              `json:"name"`
	LegalName    string              `json:"legal_name"`
	Email        string              `json:"email" gorm:"index"`
	Nik          string              `json:"nik" gorm:"index;type:varchar(16)"`
	BirthPlace   string              `json:"birth_place"`
	BirthDate    time.Time           `json:"birth_date" gorm:"type:date"`
	Password     string              `json:"password"`
	Salary       decimal.NullDecimal `json:"salary" gorm:"type:decimal(20,2)"`
	ValidatedAt  sql.NullTime        `json:"validated_at" gorm:"index"`
	SelfieFile   string              `json:"selfie_file"`
	KtpFile      string              `json:"ktp_file"`
	Role         Role                `json:"role" gorm:"foreignKey:RoleID"`
	Limits       []Limit             `json:"limits" gorm:"foreignKey:UserID"`
	Transactions []Transaction       `json:"transactions" gorm:"foreignKey:UserID"`
	gorm.Model
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
	return
}
