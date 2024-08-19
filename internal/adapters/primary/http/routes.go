package http

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *ProductHandler) {
	app.Post("/api/v1/products", handler.CreateProduct)
	app.Get("/api/v1/products/:id", handler.GetProductByID)
	app.Put("/api/v1/products/:id", handler.UpdateProduct)
	app.Delete("/api/v1/products/:id", handler.DeleteProduct)
	app.Get("/api/v1/products", handler.ListProducts)
}
