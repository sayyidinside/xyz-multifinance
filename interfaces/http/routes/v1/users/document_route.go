package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterDocumentRoutes(route fiber.Router, handler handler.DocumentHandler) {
	document := route.Group("/document")

	document.Use(middleware.Authentication())

	document.Put(
		"/",
		handler.UpdateDocument,
	)

	document.Get(
		"/",
		handler.GetDocument,
	)
}
