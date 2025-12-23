package bootstrap

import (
     "github.com/gofiber/fiber/v2"
     "backend/routes"
     "log"
)

func StartServer() {
     app := fiber.New()
     routes.ApiRoutes(app)
     log.Println("Server running on port 8080")
     app.Listen(":8080")
}
