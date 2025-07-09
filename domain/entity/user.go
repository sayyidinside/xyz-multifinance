package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UUID        uuid.UUID    `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	RoleID      uint         `json:"role_id" gorm:"index;not null"`
	Username    string       `json:"username" gorm:"size:100;uniqueIndex;not null"`
	Email       string       `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password    string       `json:"-"`
	ValidatedAt sql.NullTime `json:"validated_at" gorm:"index"`

	// Relationship
	Role         Role          `json:"role" gorm:"foreignKey:RoleID"`
	Limits       []Limit       `json:"limits" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID"`
	Profile      *UserProfile  `json:"profile" gorm:"foreignKey:UserID"`
	Document     *UserDocument `json:"documents" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
	return
}
