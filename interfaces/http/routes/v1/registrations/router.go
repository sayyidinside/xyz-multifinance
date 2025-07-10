package registrations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterRoutes(route fiber.Router, handler handler.RegistrationHandler) {
	registrationRoutes := route.Group("/registration")

	registrationRoutes.Post("/", handler.Register)

	registrationRoutes.Post(
		"/activate/:uuid",
		middleware.Authentication(),
		middleware.Authorization(true, true, []string{}),
		handler.Activate,
	)
	registrationRoutes.Get(
		"/",
		middleware.Authentication(),
		middleware.Authorization(true, true, []string{}),
		handler.GetAllRegistration,
	)

	registrationRoutes.Get(
		"/:uuid",
		middleware.Authentication(),
		middleware.Authorization(true, true, []string{}),
		handler.GetRegistration,
	)
}
