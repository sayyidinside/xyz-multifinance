package transactions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterLimitRoutes(route fiber.Router, handler handler.LimitHandler) {
	limit := route.Group("/limit")

	limit.Use(middleware.Authentication())

	limit.Get(
		"/:user_uuid",
		middleware.Authorization(false, true, []string{}),
		handler.GetLimit,
	)
}
