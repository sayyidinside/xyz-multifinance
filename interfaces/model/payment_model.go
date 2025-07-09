package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	PaymentDetail struct {
		ID               uint            `json:"id"`
		UUID             uuid.UUID       `json:"uuid"`
		TransactionID    uint            `json:"transaction_id"`
		AssetName        string          `json:"asset_name"`
		ContractNumber   string          `json:"contract_number"`
		Amount           decimal.Decimal `json:"amount"`
		InstalmentNumber uint            `json:"instalment_number"`
		CreatedAt        time.Time       `json:"created_at"`
		UpdatedAt        time.Time       `json:"updated_at"`
	}

	PaymentList struct {
		ID               uint            `json:"id"`
		UUID             uuid.UUID       `json:"uuid"`
		TransactionID    uint            `json:"transaction_id"`
		Amount           decimal.Decimal `json:"amount"`
		InstalmentNumber uint            `json:"instalment_number"`
		CreatedAt        time.Time       `json:"created_at"`
	}

	PaymentInput struct {
		TransactionID uint   `json:"transaction_id" form:"transaction_id" xml:"transaction_id" validate:"required,numeric"`
		Amount        string `json:"amount" form:"amount" xml:"amount" validate:"required"`
	}
)

func PaymentToDetailModel(payment *entity.Payment) *PaymentDetail {
	return &PaymentDetail{
		ID:               payment.ID,
		UUID:             payment.UUID,
		TransactionID:    payment.TransactionID,
		AssetName:        payment.Transaction.AssetName,
		ContractNumber:   payment.Transaction.ContractNumber,
		Amount:           payment.Amount,
		InstalmentNumber: payment.InstalmentNumber,
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
	}
}

func PaymentToListModel(payment *entity.Payment) *PaymentList {
	return &PaymentList{
		ID:               payment.ID,
		UUID:             payment.UUID,
		TransactionID:    payment.TransactionID,
		Amount:           payment.Amount,
		InstalmentNumber: payment.InstalmentNumber,
		CreatedAt:        payment.CreatedAt,
	}
}

func PaymentToListModels(payments []entity.Payment) (listModels []PaymentList) {
	for _, payment := range payments {
		listModels = append(listModels, *PaymentToListModel(&payment))
	}

	return listModels
}

func (input *PaymentInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Amount = sanitizer.Sanitize(input.Amount)
}
