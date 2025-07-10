package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserDetail struct {
		ID           uint                `json:"id"`
		UUID         uuid.UUID           `json:"uuid"`
		RoleID       uint                `json:"role_id"`
		Role         string              `json:"role"`
		Name         string              `json:"name"`
		Username     string              `json:"username"`
		Email        string              `json:"email"`
		ValidatedAt  sql.NullTime        `json:"validated_at"`
		Profile      *UserProfileDetail  `json:"profile,omitempty"`
		Document     *UserDocumentDetail `json:"document,omitempty"`
		Limits       []LimitList         `json:"limits"`
		Transactions []TransactionList   `json:"transactions"`
		CreatedAt    time.Time           `json:"created_at"`
		UpdatedAt    time.Time           `json:"updated_at"`
	}

	LogUserInfo struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	UserList struct {
		ID       uint      `json:"id"`
		UUID     uuid.UUID `json:"uuid"`
		Name     string    `json:"name"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
		Role     string    `json:"role"`
	}

	RegistrationDetailModel struct {
		ID          uint         `json:"id"`
		UUID        uuid.UUID    `json:"uuid"`
		Username    string       `json:"username"`
		Email       string       `json:"email"`
		ValidatedAt sql.NullTime `json:"validated_at"`
		Role        string       `json:"role"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}

	RegistrationList struct {
		ID          uint         `json:"id"`
		UUID        uuid.UUID    `json:"uuid"`
		Username    string       `json:"username"`
		Email       string       `json:"email"`
		ValidatedAt sql.NullTime `json:"validated_at"`
		Role        string       `json:"role"`
	}

	UserInput struct {
		Name       string `json:"name" form:"name" validate:"required"`
		Username   string `json:"username" form:"username" validate:"required"`
		Email      string `json:"email" form:"email" validate:"required"`
		Password   string `json:"password" form:"password" validate:"required"`
		RePassword string `json:"repassword" form:"repassword" validate:"required,eqfield=Password"`
		RoleID     uint   `json:"role_id" form:"role_id" validate:"required"`
	}

	UserRegisterInput struct {
		Username   string `json:"username" form:"username" validate:"required"`
		Email      string `json:"email" form:"email" validate:"required"`
		Password   string `json:"password" form:"password" validate:"required"`
		RePassword string `json:"repassword" form:"repassword" validate:"required,eqfield=Password"`
	}

	UserUpdateInput struct {
		Name     string `json:"name" form:"name" validate:"required"`
		Username string `json:"username" form:"username" validate:"required"`
		Email    string `json:"email" form:"email" validate:"required"`
		RoleID   uint   `json:"role_id" form:"role_id" validate:"required"`
	}

	ChangePasswordInput struct {
		Password   string `json:"password" form:"password" validate:"required"`
		RePassword string `json:"repassword" form:"repassword" validate:"required,eqfield=Password"`
	}
)

func UserToDetailModel(user *entity.User) *UserDetail {
	var profileName string
	var profile *UserProfileDetail
	var document *UserDocumentDetail

	if user.Profile != nil {
		profileName = user.Profile.Name
		profile = UserProfileToDetailModel(user.Profile)
	}

	if user.Document != nil {
		document = UserDocumentToDetailModel(user.Document)
	}

	return &UserDetail{
		ID:           user.ID,
		UUID:         user.UUID,
		RoleID:       user.RoleID,
		Role:         user.Role.Name,
		Name:         profileName,
		Username:     user.Username,
		Email:        user.Email,
		ValidatedAt:  user.ValidatedAt,
		Profile:      profile,
		Document:     document,
		Limits:       LimitToListModels(user.Limits),
		Transactions: TransactionToListModels(user.Transactions),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func UserToModel(user *entity.User) *UserList {
	var profileName string

	if user.Profile != nil {
		profileName = user.Profile.Name
	}

	return &UserList{
		ID:       user.ID,
		UUID:     user.UUID,
		Name:     profileName,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
	}
}

func UserToListModel(users *[]entity.User) *[]UserList {
	listUsers := []UserList{}

	for _, user := range *users {
		listUsers = append(listUsers, *UserToModel(&user))
	}

	return &listUsers
}

func RegistrationToDetailModel(user *entity.User) *RegistrationDetailModel {
	return &RegistrationDetailModel{
		ID:          user.ID,
		UUID:        user.UUID,
		Username:    user.Username,
		Email:       user.Email,
		ValidatedAt: user.ValidatedAt,
		Role:        user.Role.Name,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func RegistrationToListModel(user *entity.User) *RegistrationList {
	return &RegistrationList{
		ID:          user.ID,
		UUID:        user.UUID,
		Username:    user.Username,
		Email:       user.Email,
		ValidatedAt: user.ValidatedAt,
		Role:        user.Role.Name,
	}
}

func RegistrationToListModels(users *[]entity.User) (listModels []RegistrationList) {
	for _, registration := range *users {
		listModels = append(listModels, *RegistrationToListModel(&registration))
	}

	return
}

func UserInputToEntity(userInput *UserInput) *entity.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

	return &entity.User{
		Username:    userInput.Username,
		Email:       userInput.Email,
		Password:    string(hashedPassword),
		RoleID:      userInput.RoleID,
		ValidatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
}

func UserUpdateInputToEntity(input *UserUpdateInput) *entity.User {

	return &entity.User{
		Username: input.Username,
		Email:    input.Email,
		RoleID:   input.RoleID,
	}
}

func (input *UserRegisterInput) ToEntity() *entity.User {
	return &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}
}

func ChangePasswordToEntity(input *ChangePasswordInput) *entity.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	return &entity.User{
		Password: string(hashedPassword),
	}
}

func (input *UserInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
	input.Username = sanitizer.Sanitize(input.Username)
	input.Email = sanitizer.Sanitize(input.Email)
	input.Password = sanitizer.Sanitize(input.Password)
	input.RePassword = sanitizer.Sanitize(input.RePassword)
}

func SanitizeUserUpdateInput(input *UserUpdateInput) {
	sanitizer := bluemonday.StrictPolicy()

	input.Name = sanitizer.Sanitize(input.Name)
	input.Username = sanitizer.Sanitize(input.Username)
	input.Email = sanitizer.Sanitize(input.Email)
}

func SanitizeChangePasswordInput(input *ChangePasswordInput) {
	sanitizer := bluemonday.StrictPolicy()

	input.Password = sanitizer.Sanitize(input.Password)
	input.RePassword = sanitizer.Sanitize(input.RePassword)
}

func (input *UserRegisterInput) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()

	input.Username = sanitizer.Sanitize(input.Username)
	input.Email = sanitizer.Sanitize(input.Email)
	input.Password = sanitizer.Sanitize(input.Password)
	input.RePassword = sanitizer.Sanitize(input.RePassword)
}
