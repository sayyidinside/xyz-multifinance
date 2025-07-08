package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	UUID        uuid.UUID     `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	Name        string        `json:"name"`
	IsAdmin     bool          `json:"is_admin" gorm:"default:false"`
	Permissions *[]Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	gorm.Model
}

func (Role) TableName() string {
	return "roles"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.UUID == uuid.Nil {
		r.UUID = uuid.New()
	}
	return
}
