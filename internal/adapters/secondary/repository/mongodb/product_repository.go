package mongodb

import (
	"context"
	"log"
	"time"

	"test-go/internal/core/entities"
	"test-go/internal/core/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepository implements the ports.ProductRepository interface
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *mongo.Database) ports.ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

// Create inserts a new product into the MongoDB collection
func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) (string, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return "", err
	}

	log.Printf("Product created with ID: %s", result.InsertedID.(primitive.ObjectID).Hex())
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FindByID retrieves a product by its ID from the MongoDB collection
func (r *ProductRepository) FindByID(ctx context.Context, id string) (*entities.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product entities.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Update modifies an existing product in the MongoDB collection
func (r *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": product.ID}
	update := bson.M{
		"$set": product,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	log.Printf("Product with ID: %s updated successfully", product.ID.Hex())
	return nil
}

// Delete removes a product by its ID from the MongoDB collection
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	log.Printf("Product with ID: %s deleted successfully", id)
	return nil
}

// FindAll retrieves all products from the MongoDB collection
func (r *ProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*entities.Product
	for cursor.Next(ctx) {
		var product entities.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
