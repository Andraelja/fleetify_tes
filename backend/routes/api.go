package routes

import (
	"backend/app/controllers"
	"backend/middleware"
	"backend/services"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App) {

	// Auth routes (no middleware)
	auth := app.Group("/auth")
	auth.Post("/register", controllers.RegisterHandler)
	auth.Post("/login", controllers.LoginHandler)

	// Home route (with auth middleware)
	home := app.Group("/home", middleware.AuthMiddleware)
	home.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello World!"})
	})

	// Supplier routes
	supplierService := services.NewSupplierService()
	supplierController := controllers.NewSupplierController(supplierService)

	supplier := app.Group("/supplier", middleware.AuthMiddleware)
	supplier.Post("/", supplierController.Create)
	supplier.Get("/", supplierController.GetAll)
	supplier.Get("/:id", supplierController.GetById)
	supplier.Put("/:id", supplierController.Update)
	supplier.Delete("/:id", supplierController.Delete)

	// Item routes
	itemService := services.NewItemService()
	itemController := controllers.NewItemController(itemService)

	item := app.Group("/item", middleware.AuthMiddleware)
	item.Post("/", itemController.Create)
	item.Get("/", itemController.GetAll)
	item.Get("/:id", itemController.GetById)
	item.Put("/:id", itemController.Update)
	item.Delete("/:id", itemController.Delete)

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = "https://webhook.site/your-unique-url"
	}
	
	webhookService := services.NewWebhookService(webhookURL)
	purchasingService := services.NewPurchasingService(webhookService)
	purchasingController := controllers.NewPurchasingController(purchasingService)

	purchasing := app.Group("/purchasing", middleware.AuthMiddleware)
	purchasing.Post("/", purchasingController.Create)
	purchasing.Get("/", purchasingController.FindAll)
	purchasing.Get("/:id", purchasingController.FindById)
	purchasing.Put("/:id", purchasingController.Update)
	purchasing.Delete("/:id", purchasingController.Delete)
}