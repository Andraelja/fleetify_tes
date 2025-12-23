package routes

import (
	"backend/app/controllers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App) {

	auth := app.Group("/auth")
	auth.Post("/register", controllers.RegisterHandler)
	auth.Post("/login", controllers.LoginHandler)

	home := app.Group("/home", middleware.AuthMiddleware)
	home.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello World!"})
	})
}
