package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterProfileRoutes(route fiber.Router, handler handler.ProfileHandler) {
	user := route.Group("/profile")

	user.Use(middleware.Authentication())

	user.Put(
		"/",
		handler.UpdateProfile,
	)
}
