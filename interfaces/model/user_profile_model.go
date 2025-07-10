package model

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
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

	UserProfileUpdate struct {
		Name       string `json:"name"`
		LegalName  string `json:"legal_name"`
		Nik        string `json:"nik"`
		BirthPlace string `json:"birth_place"`
		BirthDate  string `json:"birth_date"`
		Salary     string `json:"salary"`
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

func (input *UserProfileUpdate) ToEntity() (*entity.UserProfile, error) {
	birth_date, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, err
	}

	salary, err := decimal.NewFromString(input.Salary)
	if err != nil {
		return nil, err
	}

	return &entity.UserProfile{
		Name:       input.Name,
		LegalName:  input.LegalName,
		Nik:        input.Nik,
		BirthPlace: input.BirthPlace,
		BirthDate:  birth_date,
		Salary:     decimal.NullDecimal{Decimal: salary, Valid: true},
	}, nil
}

func (input *UserProfileUpdate) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
	input.LegalName = sanitizer.Sanitize(input.LegalName)
	input.Nik = sanitizer.Sanitize(input.Nik)
	input.BirthPlace = sanitizer.Sanitize(input.BirthPlace)
	input.BirthDate = sanitizer.Sanitize(input.BirthDate)
	input.Salary = sanitizer.Sanitize(input.Salary)
}
