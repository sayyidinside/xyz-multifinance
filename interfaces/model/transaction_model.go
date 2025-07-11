package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	TransactionDetail struct {
		ID                 uint                         `json:"id"`
		UUID               uuid.UUID                    `json:"uuid"`
		UserID             uint                         `json:"user_id"`
		CustomerName       string                       `json:"customer_name"`
		AssetName          string                       `json:"asset_name"`
		ContractNumber     string                       `json:"contract_number"`
		OnTheRoad          decimal.Decimal              `json:"on_the_road"`
		AdminFee           decimal.Decimal              `json:"admin_fee"`
		TotalLoanAmount    decimal.Decimal              `json:"total_loan_amount"`
		MonthlyInstallment decimal.Decimal              `json:"monthly_installment"`
		InterestAmount     decimal.Decimal              `json:"interest_amount"`
		Tenor              uint                         `json:"tenor"`
		StartDate          time.Time                    `json:"start_date"`
		EndDate            time.Time                    `json:"end_date"`
		Status             entity.TransactionStatus     `json:"status"`
		Installments       []TransactionInstallmentList `json:"installments"`
		Payments           []PaymentList                `json:"payments"`
		CreatedAt          time.Time                    `json:"created_at"`
		UpdatedAt          time.Time                    `json:"updated_at"`
	}

	TransactionList struct {
		ID              uint                     `json:"id"`
		UUID            uuid.UUID                `json:"uuid"`
		UserID          uint                     `json:"user_id"`
		CustomerName    string                   `json:"customer_name"`
		AssetName       string                   `json:"asset_name"`
		ContractNumber  string                   `json:"contract_number"`
		OnTheRoad       decimal.Decimal          `json:"on_the_road"`
		AdminFee        decimal.Decimal          `json:"admin_fee"`
		TotalLoanAmount decimal.Decimal          `json:"total_loan_amount"`
		Tenor           uint                     `json:"tenor"`
		StartDate       time.Time                `json:"start_date"`
		EndDate         time.Time                `json:"end_date"`
		Status          entity.TransactionStatus `json:"status"`
	}

	TransactionInput struct {
		AssetName      string `json:"asset_name" form:"asset_name" xml:"asset_name" validate:"required"`
		ContractNumber string `json:"contract_number" form:"contract_number" xml:"contract_number" validate:"required"`
		OnTheRoad      string `json:"on_the_road" form:"on_the_road" xml:"on_the_road" validate:"required,numeric"`
		AdminFee       string `json:"admin_fee" form:"admin_fee" xml:"admin_fee" validate:"required,numeric"`
		InterestAmount string `json:"interest_amount" form:"interest_amount" xml:"interest_amount" validate:"required,numeric"`
		Tenor          uint   `json:"tenor" form:"tenor" xml:"tenor" validate:"required"`
	}
)

func TransactionToDetailModel(transaction *entity.Transaction) *TransactionDetail {
	return &TransactionDetail{
		ID:                 transaction.ID,
		UUID:               transaction.UUID,
		UserID:             transaction.UserID,
		CustomerName:       transaction.User.Profile.Name,
		AssetName:          transaction.AssetName,
		ContractNumber:     transaction.ContractNumber,
		OnTheRoad:          transaction.OnTheRoad,
		AdminFee:           transaction.AdminFee,
		TotalLoanAmount:    transaction.TotalLoanAmount,
		MonthlyInstallment: transaction.MonthlyInstallment,
		InterestAmount:     transaction.InterestAmount,
		Tenor:              transaction.Tenor,
		Status:             transaction.Status,
		StartDate:          transaction.StartDate,
		EndDate:            transaction.EndDate,
		Installments:       TransactionInstallmentToListModels(transaction.Installments),
		Payments:           PaymentToListModels(transaction.Payments),
		CreatedAt:          transaction.CreatedAt,
		UpdatedAt:          transaction.UpdatedAt,
	}
}

func TransactionToListModel(transaction *entity.Transaction) *TransactionList {
	return &TransactionList{
		ID:              transaction.ID,
		UUID:            transaction.UUID,
		UserID:          transaction.UserID,
		CustomerName:    transaction.User.Profile.Name,
		AssetName:       transaction.AssetName,
		ContractNumber:  transaction.ContractNumber,
		OnTheRoad:       transaction.OnTheRoad,
		AdminFee:        transaction.AdminFee,
		TotalLoanAmount: transaction.TotalLoanAmount,
		Tenor:           transaction.Tenor,
		StartDate:       transaction.StartDate,
		EndDate:         transaction.EndDate,
		Status:          transaction.Status,
	}
}

func TransactionToListModels(transactions []entity.Transaction) (listModels []TransactionList) {
	for _, transaction := range transactions {
		listModels = append(listModels, *TransactionToListModel(&transaction))
	}

	return listModels
}

func (input *TransactionInput) ToEntity() (*entity.Transaction, error) {
	otr_decimal, err := decimal.NewFromString(input.OnTheRoad)
	if err != nil {
		return nil, err
	}
	admin_decimal, err := decimal.NewFromString(input.AdminFee)
	if err != nil {
		return nil, err
	}
	interest_decimal, err := decimal.NewFromString(input.InterestAmount)
	if err != nil {
		return nil, err
	}
	tenor_decimal := decimal.NewFromUint64(uint64(input.Tenor))
	monthly_decimal := ((otr_decimal.Add(admin_decimal)).Div(tenor_decimal)).Add(interest_decimal)

	return &entity.Transaction{
		AssetName:          input.AssetName,
		ContractNumber:     input.ContractNumber,
		OnTheRoad:          otr_decimal,
		AdminFee:           admin_decimal,
		MonthlyInstallment: monthly_decimal,
		InterestAmount:     interest_decimal,
		StartDate:          time.Now(),
		EndDate:            time.Now().AddDate(0, int(input.Tenor), 1),
		Tenor:              input.Tenor,
	}, nil
}

func (input *TransactionInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.AssetName = sanitizer.Sanitize(input.AssetName)
	input.ContractNumber = sanitizer.Sanitize(input.ContractNumber)
	input.OnTheRoad = sanitizer.Sanitize(input.OnTheRoad)
	input.AdminFee = sanitizer.Sanitize(input.AdminFee)
	input.InterestAmount = sanitizer.Sanitize(input.InterestAmount)
}
