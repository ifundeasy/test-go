package application

import (
	"context"
	"encoding/json"
	"log"

	queue "test-go/internal/adapters/secondary/messaging"
	"test-go/internal/core/entities"
	"test-go/internal/core/ports"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	repo        ports.ProductRepository
	mongoDB     *mongo.Database
	redisClient *redis.Client
	queue       *queue.RabbitMQ
}

// NewProductService creates a new instance of ProductService
func NewProductService(mongoDB *mongo.Database, redisClient *redis.Client, queue *queue.RabbitMQ) *ProductService {
	return &ProductService{
		mongoDB:     mongoDB,
		redisClient: redisClient,
		queue:       queue,
	}
}

// CreateProduct handles the creation of a new product
func (s *ProductService) CreateProduct(ctx context.Context, name string, price float32) (string, error) {
	product := &entities.Product{
		Name:  name,
		Price: price,
	}

	// Save the product to MongoDB
	id, err := s.repo.Create(ctx, product)
	if err != nil {
		return "", err
	}

	// Set product ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	product.ID = objID

	// Convert the product to JSON and save to Redis
	productJSON, _ := json.Marshal(product)
	s.redisClient.Set(ctx, "product:"+objID.Hex(), productJSON, 0)

	// Publish event to the queue
	err = s.queue.PublishProductCreated(product)
	if err != nil {
		log.Printf("Failed to publish product created event: %v", err)
	}

	return objID.Hex(), nil
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*entities.Product, error) {
	// Try to get product from Redis
	val, err := s.redisClient.Get(ctx, "product:"+id).Result()
	if err == redis.Nil {
		// If not found in Redis, fetch from MongoDB
		product, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		// Save the product in Redis for future requests
		productJSON, _ := json.Marshal(product)
		s.redisClient.Set(ctx, "product:"+id, productJSON, 0)

		return product, nil
	} else if err != nil {
		return nil, err
	}

	// If found in Redis, unmarshal the JSON
	product := &entities.Product{}
	if err := json.Unmarshal([]byte(val), product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct handles updating an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id string, name string, price float32) error {
	// Retrieve and update the product
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	product.Name = name
	product.Price = price

	// Update in MongoDB
	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	// Convert the updated product to JSON and save to Redis
	productJSON, _ := json.Marshal(product)
	s.redisClient.Set(ctx, "product:"+id, productJSON, 0)

	// Publish event to the queue
	err = s.queue.PublishProductUpdated(product)
	if err != nil {
		log.Printf("Failed to publish product updated event: %v", err)
	}

	return nil
}

// DeleteProduct handles deleting a product by its ID
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// Delete from MongoDB
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete from Redis
	s.redisClient.Del(ctx, "product:"+id)

	// Publish event to the queue
	err := s.queue.PublishProductDeleted(id)
	if err != nil {
		log.Printf("Failed to publish product deleted event: %v", err)
	}

	return nil
}

// ListProducts retrieves all products
func (s *ProductService) ListProducts(ctx context.Context) ([]*entities.Product, error) {
	// Example without Redis caching for list operation
	return s.repo.FindAll(ctx)
}
