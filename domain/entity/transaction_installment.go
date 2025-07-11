package entity

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusOverdue PaymentStatus = "overdue"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type TransactionInstallment struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	UUID              uuid.UUID       `json:"uuid" gorm:"uniqueIndex;type:char(36);not null"`
	TransactionID     uint            `json:"transaction_id" gorm:"not null"`
	InstallmentNumber uint            `json:"installment_number" gorm:"not null;type:smallint unsigned"`
	DueDate           time.Time       `json:"due_date" gorm:"type:date;not null"`
	AmountDue         decimal.Decimal `json:"amount_due" gorm:"type:decimal(20,2);not null"`
	AmountPaid        decimal.Decimal `json:"amount_paid" gorm:"type:decimal(20,2);default:0"`
	PaymentStatus     PaymentStatus   `json:"payment_status" gorm:"type:enum('pending', 'partial', 'paid', 'overdue', 'failed');default:'pending'"`
	PaidAt            sql.NullTime    `json:"paid_at"`

	// Relationships
	Transaction Transaction `json:"transaction" gorm:"foreignKey:TransactionID"`
	Payments    []Payment   `json:"payments" gorm:"foreignKey:InstallmentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (TransactionInstallment) TableName() string {
	return "transaction_installments"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (t *TransactionInstallment) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == uuid.Nil {
		t.UUID = uuid.New()
	}
	return
}

func (p *PaymentStatus) Scan(value interface{}) error {
	*p = PaymentStatus(value.([]byte))
	return nil
}

func (p PaymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}
