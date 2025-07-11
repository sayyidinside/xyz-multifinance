package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	TransactionInstallmentDetail struct {
		ID                uint                 `json:"id"`
		UUID              uuid.UUID            `json:"uuid"`
		TransactionID     uint                 `json:"transaction_id"`
		InstallmentNumber uint                 `json:"installment_number"`
		AssetName         string               `json:"asset_name"`
		ContractNumber    string               `json:"contract_number"`
		Tenor             uint                 `json:"tenor"`
		DueDate           time.Time            `json:"due_date"`
		AmountDue         decimal.Decimal      `json:"amount_due"`
		AmountPaid        decimal.Decimal      `json:"amount_paid"`
		PaymentStatus     entity.PaymentStatus `json:"payment_status"`
		PaidAt            sql.NullTime         `json:"paid_at"`
		Payments          []PaymentList        `json:"payments"`
		CreatedAt         time.Time            `json:"created_at"`
		UpdatedAt         time.Time            `json:"updated_at"`
	}

	TransactionInstallmentList struct {
		ID                uint                 `json:"id"`
		UUID              uuid.UUID            `json:"uuid"`
		TransactionID     uint                 `json:"transaction_id"`
		InstallmentNumber uint                 `json:"installment_number"`
		AssetName         string               `json:"asset_name"`
		ContractNumber    string               `json:"contract_number"`
		DueDate           time.Time            `json:"due_date"`
		AmountDue         decimal.Decimal      `json:"amount_due"`
		AmountPaid        decimal.Decimal      `json:"amount_paid"`
		PaymentStatus     entity.PaymentStatus `json:"payment_status"`
		PaidAt            sql.NullTime         `json:"paid_at"`
	}
)

func TransactionInstallmentToDetailModel(transactionInstallment *entity.TransactionInstallment) *TransactionInstallmentDetail {
	return &TransactionInstallmentDetail{
		ID:                transactionInstallment.ID,
		UUID:              transactionInstallment.UUID,
		TransactionID:     transactionInstallment.TransactionID,
		InstallmentNumber: transactionInstallment.InstallmentNumber,
		AssetName:         transactionInstallment.Transaction.AssetName,
		ContractNumber:    transactionInstallment.Transaction.ContractNumber,
		Tenor:             transactionInstallment.Transaction.Tenor,
		DueDate:           transactionInstallment.DueDate,
		AmountDue:         transactionInstallment.AmountDue,
		AmountPaid:        transactionInstallment.AmountPaid,
		PaymentStatus:     transactionInstallment.PaymentStatus,
		PaidAt:            transactionInstallment.PaidAt,
		Payments:          PaymentToListModels(transactionInstallment.Payments),
		CreatedAt:         transactionInstallment.CreatedAt,
		UpdatedAt:         transactionInstallment.UpdatedAt,
	}
}

func TransactionInstallmentToListModel(transactionInstallment *entity.TransactionInstallment) *TransactionInstallmentList {
	return &TransactionInstallmentList{
		ID:                transactionInstallment.ID,
		UUID:              transactionInstallment.UUID,
		TransactionID:     transactionInstallment.TransactionID,
		InstallmentNumber: transactionInstallment.InstallmentNumber,
		AssetName:         transactionInstallment.Transaction.AssetName,
		ContractNumber:    transactionInstallment.Transaction.ContractNumber,
		DueDate:           transactionInstallment.DueDate,
		AmountDue:         transactionInstallment.AmountDue,
		AmountPaid:        transactionInstallment.AmountPaid,
		PaymentStatus:     transactionInstallment.PaymentStatus,
		PaidAt:            transactionInstallment.PaidAt,
	}
}

func TransactionInstallmentToListModels(transactionInstallments []entity.TransactionInstallment) (listModels []TransactionInstallmentList) {
	for _, installment := range transactionInstallments {
		listModels = append(listModels, *TransactionInstallmentToListModel(&installment))
	}

	return
}
