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
		ID          uint         `json:"id"`
		UUID        uuid.UUID    `json:"uuid"`
		RoleID      uint         `json:"role_id"`
		Role        string       `json:"role"`
		Name        string       `json:"name"`
		Username    string       `json:"username"`
		Email       string       `json:"email"`
		ValidatedAt sql.NullTime `json:"validated_at"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
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

	UserInput struct {
		Name       string `json:"name" form:"name" validate:"required"`
		Username   string `json:"username" form:"username" validate:"required"`
		Email      string `json:"email" form:"email" validate:"required"`
		Password   string `json:"password" form:"password" validate:"required"`
		RePassword string `json:"repassword" form:"repassword" validate:"required,eqfield=Password"`
		RoleID     uint   `json:"role_id" form:"role_id" validate:"required"`
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
	return &UserDetail{
		ID:          user.ID,
		UUID:        user.UUID,
		RoleID:      user.RoleID,
		Role:        user.Role.Name,
		Username:    user.Username,
		Email:       user.Email,
		ValidatedAt: user.ValidatedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func UserToModel(user *entity.User) *UserList {
	return &UserList{
		ID:       user.ID,
		UUID:     user.UUID,
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

func ChangePasswordToEntity(input *ChangePasswordInput) *entity.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	return &entity.User{
		Password: string(hashedPassword),
	}
}

func SanitizeUserInput(input *UserInput) {
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
