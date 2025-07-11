package entity

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                 uint              `json:"id" gorm:"primaryKey"`
	UUID               uuid.UUID         `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	UserID             uint              `json:"user_id" gorm:"index;not null"`
	AssetName          string            `json:"asset_name" gorm:"not null"`
	ContractNumber     string            `json:"contract_number" gorm:"type:varchar(255);index;not null"`
	OnTheRoad          decimal.Decimal   `json:"on_the_road" gorm:"type:decimal(20,2);not null"`
	AdminFee           decimal.Decimal   `json:"admin_fee" gorm:"type:decimal(20,2);not null"`
	TotalLoanAmount    decimal.Decimal   `json:"total_loan_amount" gorm:"->;type:decimal(20,2) GENERATED ALWAYS AS (on_the_road + admin_fee) STORED"`
	MonthlyInstallment decimal.Decimal   `json:"monthly_installment" gorm:"type:decimal(20,2);not null"`
	InterestAmount     decimal.Decimal   `json:"interest_amount" gorm:"type:decimal(20,2);not null"`
	Tenor              uint              `json:"tenor" gorm:"type:smallint unsigned;not null"`
	StartDate          time.Time         `json:"start_date" gorm:"type:date;not null"`
	EndDate            time.Time         `json:"end_date" gorm:"type:date"`
	Status             TransactionStatus `json:"status" gorm:"type:enum('active', 'paid', 'canceled');default:'active'"`

	// Relationship
	User         User                     `json:"user" gorm:"foreignKey:UserID"`
	Installments []TransactionInstallment `json:"installments" gorm:"foreignKey:TransactionID"`
	Payments     []Payment                `json:"payments" gorm:"foreignKey:TransactionID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
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

type TransactionStatus string

const (
	TransactionActive   TransactionStatus = "active"
	TransactionPaid     TransactionStatus = "paid"
	TransactionCanceled TransactionStatus = "canceled"
)

func (t *TransactionStatus) Scan(value interface{}) error {
	*t = TransactionStatus(value.([]byte))
	return nil
}

func (t TransactionStatus) Value() (driver.Value, error) {
	return string(t), nil
}
