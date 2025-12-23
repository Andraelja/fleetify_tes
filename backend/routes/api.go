package routes

import (
	"backend/app/controllers"
	"backend/middleware"
	"backend/services"

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

	supplierService := services.NewSupplierService()
	supplierController := controllers.NewSupplierController(supplierService)

	supplier := app.Group("/supplier", middleware.AuthMiddleware)
	supplier.Post("/", supplierController.Create)
	supplier.Get("/", supplierController.GetAll)
	supplier.Get("/:id", supplierController.GetById)
	supplier.Put("/:id", supplierController.Update)
	supplier.Delete("/:id", supplierController.Delete)
}
