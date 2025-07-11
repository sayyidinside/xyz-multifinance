package transactions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterInstallmentRoutes(route fiber.Router, handler handler.InstallmentHandler) {
	installment := route.Group("/installment")

	installment.Use(middleware.Authentication())

	installment.Get(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.GetAllInstallment,
	)
	installment.Get(
		"/:uuid",
		middleware.Authorization(false, true, []string{}),
		handler.GetInstallment,
	)
}
