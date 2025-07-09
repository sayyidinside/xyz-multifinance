package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TransactionActive   TransactionStatus = "active"
	TransactionPaid     TransactionStatus = "paid"
	TransactionCanceled TransactionStatus = "canceled"
)

type Transaction struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	UUID             uuid.UUID         `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	UserID           uint              `json:"user_id"`
	AssetName        string            `json:"asset_name"`
	ContractNumber   string            `json:"contract_number" gorm:"type:varchar(255);uniqueIndex"`
	OnTheRoad        decimal.Decimal   `json:"on_the_road" gorm:"type:decimal(20,2)"`
	AdminFee         decimal.Decimal   `json:"admin_fee" gorm:"type:decimal(20,2)"`
	InstalmentAmount decimal.Decimal   `json:"instalment_amount" gorm:"type:decimal(20,2)"`
	InterestAmount   decimal.Decimal   `json:"interest_amount" gorm:"type:decimal(20,2)"`
	Tenor            uint              `json:"tenor" gorm:"type:smallint unsigned"`
	Status           TransactionStatus `gorm:"type:enum('active', 'paid', 'canceled');default:'active'"`
	User             User              `json:"user" gorm:"foreignKey:UserID"`
	Payments         []Payment         `json:"payments" gorm:"foreignKey:TransactionID"`
	gorm.Model
}

func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the struct.
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == uuid.Nil {
		t.UUID = uuid.New()
	}
	return
}
