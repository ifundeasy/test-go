package ports

import (
	"context"
	"errors"
	"test-go/internal/core/entities"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) (string, error)
	FindByID(ctx context.Context, id string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entities.Product, error)
}

// ErrProductNotFound is returned when a product is not found in the repository
var ErrProductNotFound = errors.New("product not found")
