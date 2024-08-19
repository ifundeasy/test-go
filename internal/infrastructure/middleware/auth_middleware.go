package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a middleware function that handles authentication
func AuthMiddleware(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")

	// Check if the header is empty
	if authHeader == "" {
		log.Println("Missing Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse the header to get the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate the token (for example, using a JWT library)
	if !validateToken(token) {
		log.Println("Invalid or expired token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Proceed to the next handler if the token is valid
	return c.Next()
}

// validateToken validates the token (this is just a placeholder, replace with actual validation logic)
func validateToken(token string) bool {
	// This is where you would typically validate the token (e.g., JWT)
	// For the sake of example, we'll just check if the token equals "valid-token"
	return token == "valid-token"
}
