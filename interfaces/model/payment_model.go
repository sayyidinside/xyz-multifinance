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
		ID                uint            `json:"id"`
		UUID              uuid.UUID       `json:"uuid"`
		TransactionID     uint            `json:"transaction_id"`
		InstallmentID     uint            `json:"installment_id"`
		Amount            decimal.Decimal `json:"amount"`
		PaymentMethod     string          `json:"payment_method"`
		AssetName         string          `json:"asset_name"`
		ContractNumber    string          `json:"contract_number"`
		InstallmentNumber uint            `json:"installment_number"`
		CreatedAt         time.Time       `json:"created_at"`
		UpdatedAt         time.Time       `json:"updated_at"`
	}

	PaymentList struct {
		ID               uint            `json:"id"`
		UUID             uuid.UUID       `json:"uuid"`
		TransactionID    uint            `json:"transaction_id"`
		InstallmentID    uint            `json:"installment_id"`
		Amount           decimal.Decimal `json:"amount"`
		InstalmentNumber uint            `json:"instalment_number"`
		PaymentMethod    string          `json:"payment_method"`
		CreatedAt        time.Time       `json:"created_at"`
	}

	PaymentInput struct {
		InstallmentID uint   `json:"installment_id" form:"installment_id" xml:"installment_id" validate:"required,numeric"`
		Amount        string `json:"amount" form:"amount" xml:"amount" validate:"required"`
		PaymentMethod string `json:"payment_method" form:"payment_method" xml:"payment_method" validate:"required"`
	}
)

func PaymentToDetailModel(payment *entity.Payment) *PaymentDetail {
	return &PaymentDetail{
		ID:                payment.ID,
		UUID:              payment.UUID,
		TransactionID:     payment.TransactionID,
		AssetName:         payment.Transaction.AssetName,
		ContractNumber:    payment.Transaction.ContractNumber,
		InstallmentID:     payment.InstallmentID,
		InstallmentNumber: payment.Installment.InstallmentNumber,
		PaymentMethod:     payment.PaymentMethod,
		Amount:            payment.Amount,
		CreatedAt:         payment.CreatedAt,
		UpdatedAt:         payment.UpdatedAt,
	}
}

func PaymentToListModel(payment *entity.Payment) *PaymentList {
	return &PaymentList{
		ID:               payment.ID,
		UUID:             payment.UUID,
		TransactionID:    payment.TransactionID,
		InstallmentID:    payment.InstallmentID,
		InstalmentNumber: payment.Installment.InstallmentNumber,
		PaymentMethod:    payment.PaymentMethod,
		Amount:           payment.Amount,
		CreatedAt:        payment.CreatedAt,
	}
}

func PaymentToListModels(payments []entity.Payment) (listModels []PaymentList) {
	for _, payment := range payments {
		listModels = append(listModels, *PaymentToListModel(&payment))
	}

	return listModels
}

func (input *PaymentInput) ToEntity() (*entity.Payment, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return nil, err
	}

	return &entity.Payment{
		InstallmentID: input.InstallmentID,
		Amount:        amount,
		PaymentMethod: input.PaymentMethod,
	}, nil
}

func (input *PaymentInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Amount = sanitizer.Sanitize(input.Amount)
	input.PaymentMethod = sanitizer.Sanitize(input.PaymentMethod)
}
