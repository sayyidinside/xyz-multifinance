package transactions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterTransactionRoutes(route fiber.Router, handler handler.TransactionHandler) {
	transaction := route.Group("/data")

	transaction.Use(middleware.Authentication())

	transaction.Get(
		"/:uuid",
		middleware.Authorization(false, true, []string{}),
		handler.GetTransaction,
	)

	transaction.Get(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.GetAllTransaction,
	)

	transaction.Post(
		"/",
		middleware.Authorization(false, true, []string{}),
		handler.Create,
	)

	transaction.Post(
		"/cancel/:uuid",
		middleware.Authorization(false, true, []string{}),
		handler.Cancel,
	)
}
