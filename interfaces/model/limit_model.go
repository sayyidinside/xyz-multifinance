package model

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	LimitList struct {
		ID            uint            `json:"id"`
		UserID        uint            `json:"user_id"`
		Tenor         uint            `json:"tenor"`
		CurrentLimit  decimal.Decimal `json:"current_limit"`
		OriginalLimit decimal.Decimal `json:"original_limit"`
	}
)

func LimitToListModel(limit *entity.Limit) *LimitList {
	return &LimitList{
		ID:            limit.ID,
		UserID:        limit.UserID,
		Tenor:         limit.Tenor,
		CurrentLimit:  limit.CurrentLimit,
		OriginalLimit: limit.OriginalLimit,
	}
}

func LimitToListModels(limits []entity.Limit) (listModels []LimitList) {
	for _, limit := range limits {
		listModels = append(listModels, *LimitToListModel(&limit))
	}

	return listModels
}
