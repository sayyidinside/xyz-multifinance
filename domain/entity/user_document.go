package entity

import "time"

type UserDocument struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id" gorm:"not null;uniqueIndex"`
	SelfieFile string `json:"selfie_file" gorm:"not null"`
	KtpFile    string `json:"ktp_file" gorm:"not null"`

	// Relationship
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserDocument) TableName() string {
	return "user_documents"
}
