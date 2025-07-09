package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDocument struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UUID       uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	SelfieFile string    `json:"selfie_file"`
	KtpFile    string    `json:"ktp_file"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	gorm.Model
}

func (UserDocument) TableName() string {
	return "user_documents"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (up *UserDocument) BeforeCreate(tx *gorm.DB) (err error) {
	if up.UUID == uuid.Nil {
		up.UUID = uuid.New()
	}
	return
}
