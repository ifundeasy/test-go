package repository

import (
	"context"
	"errors"
	"test-go/internal/product"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoProductRepository implements the ProductRepository interface using MongoDB
type MongoProductRepository struct {
	collection *mongo.Collection
}

// NewMongoProductRepository creates a new instance of MongoProductRepository
func NewMongoProductRepository(mongoDatabase *mongo.Database) product.ProductRepository {
	collection := mongoDatabase.Collection("products")
	return &MongoProductRepository{collection}
}

// Fetch retrieves all products from the MongoDB collection
func (r *MongoProductRepository) Fetch() ([]product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []product.Product
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var prod product.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID retrieves a product by its ID from the MongoDB collection
func (r *MongoProductRepository) GetByID(id string) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var prod product.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&prod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &prod, nil
}

// Create inserts a new product into the MongoDB collection
func (r *MongoProductRepository) Create(prod *product.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prod.ID = primitive.NewObjectID().Hex()
	_, err := r.collection.InsertOne(ctx, prod)
	return err
}

// Update modifies an existing product in the MongoDB collection
func (r *MongoProductRepository) Update(id string, prod *product.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": prod,
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a product by its ID from the MongoDB collection
func (r *MongoProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
