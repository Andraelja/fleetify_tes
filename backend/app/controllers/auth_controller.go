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
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}
	token, err := services.LoginUser(data["username"], data["password"])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	return c.JSON(responses.LoginResponse{
		Success: true,
		Token:   token,
		Data:    data,
	})
}
