package registrations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterRoutes(route fiber.Router, handler handler.RegistrationHandler) {
	registrationRoutes := route.Group("/registration")

	registrationRoutes.Post("/", handler.Register)
	registrationRoutes.Post("/activate/:uuid", handler.Activate)
	registrationRoutes.Get("/", handler.GetAllRegistration)
	registrationRoutes.Get("/:uuid", handler.GetRegistration)
}
