package transactions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterRoutes(route fiber.Router, handler *handler.TransactionManagementHandler) {
	transactions := route.Group("/transactions/")

	RegisterLimitRoutes(transactions, handler.LimitHandler)
	RegisterTransactionRoutes(transactions, handler.TransactionHandler)
}
