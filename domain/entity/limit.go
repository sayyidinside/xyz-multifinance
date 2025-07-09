package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Limit struct {
	ID          uint            `json:"id"`
	UUID        uuid.UUID       `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	UserID      uint            `json:"user_id"`
	Tenor       uint            `json:"tenor" gorm:"type:smallint unsigned"`
	LimitAmount decimal.Decimal `json:"limit_amount" gorm:"type:decimal(20,2)"`
	User        User            `json:"user" gorm:"foreignKey:UserID"`
	gorm.Model
}

func (Limit) TableName() string {
	return "limits"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (l *Limit) BeforeCreate(tx *gorm.DB) (err error) {
	if l.UUID == uuid.Nil {
		l.UUID = uuid.New()
	}
	return
}
