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
		UserID             uint   `json:"user_id" form:"user_id" xml:"user_id" validate:"required,numeric"`
		AssetName          string `json:"asset_name" form:"asset_name" xml:"asset_name" validate:"required"`
		ContractNumber     string `json:"customer_number" form:"customer_number" xml:"customer_number" validate:"required"`
		OnTheRoad          string `json:"on_the_road" form:"on_the_road" xml:"on_the_road" validate:"required"`
		AdminFee           string `json:"admin_fee" form:"admin_fee" xml:"admin_fee" validate:"required"`
		MonthlyInstallment string `json:"installment_amount" form:"installment_amount" xml:"installment_amount" validate:"required"`
		InterestAmount     string `json:"interest_amount" form:"interest_amount" xml:"interest_amount" validate:"required"`
		Tenor              uint   `json:"tenor" form:"tenor" xml:"tenor" validate:"required"`
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

func (input *TransactionInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.AssetName = sanitizer.Sanitize(input.AssetName)
	input.ContractNumber = sanitizer.Sanitize(input.ContractNumber)
	input.OnTheRoad = sanitizer.Sanitize(input.OnTheRoad)
	input.AdminFee = sanitizer.Sanitize(input.AdminFee)
	input.MonthlyInstallment = sanitizer.Sanitize(input.MonthlyInstallment)
	input.InterestAmount = sanitizer.Sanitize(input.InterestAmount)
}
