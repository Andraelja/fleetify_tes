package controllers

import (
	"backend/services"

	"backend/app/responses"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}
	err := services.RegisterUser(data["username"], data["password"], data["role"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.RegisterResponse{
		Success: true,
		Message: "Register successfully!",
		Data:    data,
	})
}

func LoginHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	token, err := services.LoginUser(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid username or password",
		})
	}

	return c.JSON(responses.LoginResponse{
		Success: true,
		Token:   token,
		Data: fiber.Map{
			"username": req.Username,
		},
	})
}
