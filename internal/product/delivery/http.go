package delivery

import (
	"net/http"
	"test-go/internal/product"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a generic error response
type ErrorResponse map[string]interface{}

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	usecase  product.ProductUsecase
	validate *validator.Validate
}

// NewProductHandler registers the routes for product handling
func NewProductHandler(app *fiber.App, u product.ProductUsecase) {
	handler := &ProductHandler{
		usecase:  u,
		validate: validator.New(),
	}

	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Post("/products", handler.CreateProduct)
	app.Put("/products/:id", handler.UpdateProduct)
	app.Delete("/products/:id", handler.DeleteProduct)
}

// GetProducts handles the HTTP request to retrieve all products
// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Produce  json
// @Success 200 {array} product.Product
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.usecase.Fetch()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{"error": err.Error()})
	}
	return c.JSON(products)
}

// GetProduct handles the HTTP request to retrieve a product by ID
// @Summary Get a product by ID
// @Description Get details of a product by its ID
// @Tags products
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} product.Product
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.usecase.GetByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{"error": "Product not found"})
	}

	return c.JSON(product)
}

// CreateProduct handles the HTTP request to create a new product
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body product.Product true "Product Data"
// @Success 201 {object} product.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product product.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{"error": "Failed to parse request body"})
	}

	if err := h.validate.Struct(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{"error": err.Error()})
	}

	if err := h.usecase.Create(&product); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(product)
}

// UpdateProduct handles the HTTP request to update an existing product by ID
// @Summary Update an existing product
// @Description Update the details of an existing product by its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body product.Product true "Product Data"
// @Success 200 {object} product.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product product.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{"error": "Failed to parse request body"})
	}

	product.ID = id

	if err := h.validate.Struct(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{"error": err.Error()})
	}

	if err := h.usecase.Update(id, &product); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{"error": err.Error()})
	}

	return c.JSON(product)
}

// DeleteProduct handles the HTTP request to delete a product by ID
// @Summary Delete a product by ID
// @Description Delete an existing product by its ID
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.usecase.Delete(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{"error": err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}
