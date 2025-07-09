package model

import (
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	LimitList struct {
		ID          uint            `json:"id"`
		UUID        uuid.UUID       `json:"uuid"`
		UserID      uint            `json:"user_id"`
		Tenor       uint            `json:"tenor"`
		LimitAmount decimal.Decimal `json:"limit_amount"`
	}
)

func LimitToListModel(limit *entity.Limit) *LimitList {
	return &LimitList{
		ID:          limit.ID,
		UUID:        limit.UUID,
		UserID:      limit.UserID,
		Tenor:       limit.Tenor,
		LimitAmount: limit.LimitAmount,
	}
}

func LimitToListModels(limits []entity.Limit) (listModels []LimitList) {
	for _, limit := range limits {
		listModels = append(listModels, *LimitToListModel(&limit))
	}

	return listModels
}
