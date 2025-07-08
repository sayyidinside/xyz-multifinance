package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UUID     uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	Name     string    `json:"name"`
	ModuleID uint      `json:"module_id"`
	Module   Module    `json:"module" gorm:"foreignKey:ModuleID"`
	gorm.Model
}

func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UUID == uuid.Nil {
		p.UUID = uuid.New()
	}
	return
}
