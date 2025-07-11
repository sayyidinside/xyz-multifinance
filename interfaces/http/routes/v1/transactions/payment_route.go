package transactions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterPaymentRoutes(route fiber.Router, handler handler.PaymentHandler) {
	payment := route.Group("/payment")

	payment.Use(middleware.Authentication())

	payment.Get(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.GetAllPayment,
	)

	payment.Get(
		"/:uuid",
		middleware.Authorization(false, true, []string{}),
		handler.GetPayment,
	)

	payment.Post(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.Create,
	)
}
