package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
)

func RegisterUserRoutes(route fiber.Router, handler handler.UserHandler) {
	user := route.Group("/data")

	user.Use(middleware.Authentication())

	user.Get(
		"/get-uuid",
		handler.GetUUID,
	)

	user.Get(
		"/:uuid",
		middleware.Authorization(false, false, []string{
			"View User",
			"Create User",
			"Update User",
			"Delete User",
		}),
		handler.GetUser,
	)

	user.Get(
		"/",
		middleware.Authorization(true, true, []string{}),
		handler.GetAllUser,
	)

	user.Post(
		"/",
		middleware.Authorization(false, false, []string{
			"Create User",
		}),
		handler.CreateUser,
	)

	user.Put(
		"/:uuid/reset-password",
		middleware.Authorization(true, false, []string{}),
		handler.ResetPassword,
	)

	user.Put(
		"/:id",
		middleware.Authorization(false, false, []string{
			"Update User",
		}),
		handler.UpdateUser,
	)

	user.Post(
		"/suspend/:uuid",
		middleware.Authorization(true, true, []string{
			"Delete User",
		}),
		handler.SuspendUser,
	)
}
