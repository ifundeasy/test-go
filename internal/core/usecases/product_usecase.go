package usecases

import (
	"context"

	"test-go/internal/core/entities"
	"test-go/internal/core/ports"
)

// ProductUseCase defines the use case for managing products
type ProductUseCase struct {
	repo ports.ProductRepository
}

// NewProductUseCase creates a new instance of ProductUseCase
func NewProductUseCase(repo ports.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
	}
}

// CreateProduct handles the creation of a new product
func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *entities.Product) (string, error) {
	// Additional business logic before saving the product can be added here
	return uc.repo.Create(ctx, product)
}

// GetProductByID handles retrieving a product by its ID
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entities.Product, error) {
	return uc.repo.FindByID(ctx, id)
}

// UpdateProduct handles updating an existing product
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entities.Product) error {
	// Additional business logic before updating the product can be added here
	return uc.repo.Update(ctx, product)
}

// DeleteProduct handles deleting a product by its ID
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	// Additional business logic before deleting the product can be added here
	return uc.repo.Delete(ctx, id)
}

// GetAllProducts handles retrieving all products
func (uc *ProductUseCase) GetAllProducts(ctx context.Context) ([]*entities.Product, error) {
	return uc.repo.FindAll(ctx)
}
