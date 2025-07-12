package tests

import (
	"database/sql"
	"time"

	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func GenerateAdminTestToken() string {

	user := entity.User{
		ID:       1,
		Username: "admin",
		Email:    "admin@email.id",
		Role: entity.Role{
			IsAdmin: true,
			Permissions: []entity.Permission{
				{Name: "View User"},
				{Name: "Update User"},
			},
		},
		ValidatedAt: sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour * -2)},
	}
	signedToken, _ := helpers.GenerateToken(&user, config.AppConfig.JwtAccessTime, config.AppConfig.JwtAccessPrivateSecret, false)

	return signedToken
}

func GenerateUserTestToken() string {

	user := entity.User{
		ID:       2,
		Username: "user",
		Email:    "user@email.id",
		Role: entity.Role{
			IsAdmin: false,
			Permissions: []entity.Permission{
				{Name: "View User"},
				{Name: "Update User"},
			},
		},
		ValidatedAt: sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour * -2)},
	}
	signedToken, _ := helpers.GenerateToken(&user, config.AppConfig.JwtAccessTime, config.AppConfig.JwtAccessPrivateSecret, false)

	return signedToken
}
