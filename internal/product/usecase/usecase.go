package usecase

import (
	"errors"
	"test-go/internal/pkg/queue"
	"test-go/internal/product" // Import the package that contains the interfaces

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// productUsecase implements the ProductUsecase interface
type productUsecase struct {
	productRepo product.ProductRepository
	redisClient *redis.Client
	rabbitMQ    *queue.RabbitMQ
}

// NewProductUsecase creates a new instance of ProductUsecase
func NewProductUsecase(r product.ProductRepository, redis *redis.Client, mq *queue.RabbitMQ) product.ProductUsecase {
	return &productUsecase{
		productRepo: r,
		redisClient: redis,
		rabbitMQ:    mq,
	}
}

// Fetch retrieves all products from the repository
func (u *productUsecase) Fetch() ([]product.Product, error) {
	// Try to fetch from cache first
	products, err := u.getCachedProducts()
	if err == nil && len(products) > 0 {
		return products, nil
	}

	// Fetch from database
	products, err = u.productRepo.Fetch()
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = u.cacheProducts(products)

	return products, nil
}

// GetByID retrieves a product by its ID
func (u *productUsecase) GetByID(id string) (*product.Product, error) {
	// Try to fetch from cache first
	prod, err := u.getCachedProductByID(id)
	if err == nil && prod != nil {
		return prod, nil
	}

	// Fetch from database
	prod, err = u.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = u.cacheProductByID(prod)

	return prod, nil
}

// Create adds a new product to the repository
func (u *productUsecase) Create(prod *product.Product) error {
	// Generate a new ID for the product
	prod.ID = primitive.NewObjectID().Hex()

	// Insert the product into the database
	if err := u.productRepo.Create(prod); err != nil {
		return err
	}

	// Publish to queue
	_ = u.rabbitMQ.PublishProductCreated(prod)

	// Invalidate cache
	_ = u.invalidateCache()

	return nil
}

// Update modifies an existing product in the repository
func (u *productUsecase) Update(id string, prod *product.Product) error {
	// Update the product in the database
	if err := u.productRepo.Update(id, prod); err != nil {
		return err
	}

	// Publish to queue
	_ = u.rabbitMQ.PublishProductUpdated(prod)

	// Invalidate cache
	_ = u.invalidateCache()

	return nil
}

// Delete removes a product by its ID from the repository
func (u *productUsecase) Delete(id string) error {
	// Delete the product from the database
	if err := u.productRepo.Delete(id); err != nil {
		return err
	}

	// Publish to queue
	_ = u.rabbitMQ.PublishProductDeleted(id)

	// Invalidate cache
	_ = u.invalidateCache()

	return nil
}

// Private methods for cache management

// getCachedProducts tries to retrieve the list of products from Redis cache
func (u *productUsecase) getCachedProducts() ([]product.Product, error) {
	// Implement Redis cache retrieval logic here
	return nil, errors.New("not implemented")
}

// cacheProducts caches the list of products in Redis
func (u *productUsecase) cacheProducts(prods []product.Product) error {
	// Implement Redis cache storing logic here
	return nil
}

// getCachedProductByID tries to retrieve a single product by ID from Redis cache
func (u *productUsecase) getCachedProductByID(id string) (*product.Product, error) {
	// Implement Redis cache retrieval logic here
	return nil, errors.New("not implemented")
}

// cacheProductByID caches a single product by ID in Redis
func (u *productUsecase) cacheProductByID(prod *product.Product) error {
	// Implement Redis cache storing logic here
	return nil
}

// invalidateCache invalidates the Redis cache for products
func (u *productUsecase) invalidateCache() error {
	// Implement Redis cache invalidation logic here
	return nil
}
