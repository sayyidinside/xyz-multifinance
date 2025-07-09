package model

import (
	"time"

	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/shopspring/decimal"
)

type (
	UserProfileDetail struct {
		ID         uint                `json:"id"`
		UserID     uint                `json:"user_id"`
		Name       string              `json:"name"`
		LegalName  string              `json:"legal_name"`
		Nik        string              `json:"nik"`
		BirthPlace string              `json:"birth_place"`
		BirthDate  time.Time           `json:"birth_date"`
		Salary     decimal.NullDecimal `json:"salary"`
	}
)

func UserProfileToDetailModel(userProfile *entity.UserProfile) *UserProfileDetail {
	return &UserProfileDetail{
		ID:         userProfile.ID,
		UserID:     userProfile.UserID,
		Name:       userProfile.Name,
		LegalName:  userProfile.LegalName,
		Nik:        userProfile.Nik,
		BirthPlace: userProfile.BirthPlace,
		BirthDate:  userProfile.BirthDate,
		Salary:     userProfile.Salary,
	}
}
