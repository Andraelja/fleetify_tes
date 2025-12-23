package middleware

import (
	"backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"message": "Missing token"})
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid token"})
	}
	c.Locals("username", claims.Username)
	return c.Next()
}
