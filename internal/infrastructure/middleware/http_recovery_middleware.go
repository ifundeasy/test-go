package middleware

import (
	"test-go/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

// RecoveryMiddleware is a middleware function that recovers from panics and logs the error
func RecoveryMiddleware(logger *logging.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use recover to catch any panic that occurs during request handling
		defer func() {
			if r := recover(); r != nil {
				// Log the panic message
				logger.Error("Recovered from panic: " + logPanic(r))

				// Respond with a 500 Internal Server Error
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			}
		}()

		// Proceed to the next middleware/handler
		return c.Next()
	}
}

// logPanic formats the panic information into a string
func logPanic(r interface{}) string {
	switch v := r.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "Unknown panic"
	}
}
