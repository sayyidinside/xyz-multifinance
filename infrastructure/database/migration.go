package database

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Module{})
	db.AutoMigrate(&entity.Permission{})
	db.AutoMigrate(&entity.Role{})
	db.AutoMigrate(&entity.RolePermission{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.UserProfile{})
	db.AutoMigrate(&entity.UserDocument{})
	db.AutoMigrate(&entity.RefreshToken{})
	db.AutoMigrate(&entity.Limit{})
	db.AutoMigrate(&entity.Transaction{})
	db.AutoMigrate(&entity.Payment{})
}
