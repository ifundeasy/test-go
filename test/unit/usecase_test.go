package unit

import (
	"context"
	"test-go/internal/product"
	"test-go/internal/product/usecase"
	"test-go/test/mock"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFetch(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock.MockMongoCollection)
	mockRedis := new(mock.MockRedisClient)
	mockMQ := new(mock.MockRabbitMQ)

	u := usecase.NewProductUsecase(mockRepo, mockRedis, mockMQ)

	expectedProducts := []product.Product{
		{ID: "1", Name: "Product 1", Price: 10.0},
		{ID: "2", Name: "Product 2", Price: 20.0},
	}

	// Simulate a cache miss and fetching from the database
	mockRedis.On("Get", ctx, "products").Return(redis.NewStringResult("", redis.Nil))
	cursor := &mongo.Cursor{} // You can replace this with an actual mock cursor implementation if necessary
	mockRepo.On("Find", ctx, bson.D{}, (*options.FindOptions)(nil)).Return(cursor, nil)
	mockRepo.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*product.Product)
		*arg = expectedProducts[0] // Simulate decoding first product
	}).Return(nil)

	products, err := u.Fetch()

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)

	mockRedis.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock.MockMongoCollection)
	mockRedis := new(mock.MockRedisClient)
	mockMQ := new(mock.MockRabbitMQ)

	u := usecase.NewProductUsecase(mockRepo, mockRedis, mockMQ)

	expectedProduct := &product.Product{ID: "1", Name: "Product 1", Price: 10.0}

	// Simulate a cache miss and fetching from the database
	mockRedis.On("Get", ctx, "product:1").Return(redis.NewStringResult("", redis.Nil))
	singleResult := &mongo.SingleResult{} // You can replace this with an actual mock single result implementation if necessary
	mockRepo.On("FindOne", ctx, bson.M{"_id": "1"}, (*options.FindOneOptions)(nil)).Return(singleResult)
	mockRepo.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*product.Product)
		*arg = *expectedProduct // Simulate decoding the product
	}).Return(nil)

	product, err := u.GetByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)

	mockRedis.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock.MockMongoCollection)
	mockRedis := new(mock.MockRedisClient)
	mockMQ := new(mock.MockRabbitMQ)

	u := usecase.NewProductUsecase(mockRepo, mockRedis, mockMQ)

	newProduct := &product.Product{ID: "1", Name: "Product 1", Price: 10.0}

	mockRepo.On("InsertOne", ctx, newProduct).Return(&mongo.InsertOneResult{}, nil)
	mockMQ.On("PublishProductCreated", newProduct).Return(nil)
	mockRedis.On("Del", ctx, "products").Return(redis.NewIntCmd(ctx))

	err := u.Create(newProduct)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockMQ.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock.MockMongoCollection)
	mockRedis := new(mock.MockRedisClient)
	mockMQ := new(mock.MockRabbitMQ)

	u := usecase.NewProductUsecase(mockRepo, mockRedis, mockMQ)

	updatedProduct := &product.Product{ID: "1", Name: "Updated Product", Price: 15.0}

	mockRepo.On("UpdateOne", ctx, bson.M{"_id": "1"}, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	mockMQ.On("PublishProductUpdated", updatedProduct).Return(nil)
	mockRedis.On("Del", ctx, "product:1", "products").Return(redis.NewIntCmd(ctx))

	err := u.Update("1", updatedProduct)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockMQ.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock.MockMongoCollection)
	mockRedis := new(mock.MockRedisClient)
	mockMQ := new(mock.MockRabbitMQ)

	u := usecase.NewProductUsecase(mockRepo, mockRedis, mockMQ)

	mockRepo.On("DeleteOne", ctx, bson.M{"_id": "1"}).Return(&mongo.DeleteResult{}, nil)
	mockMQ.On("PublishProductDeleted", "1").Return(nil)
	mockRedis.On("Del", ctx, "product:1", "products").Return(redis.NewIntCmd(ctx))

	err := u.Delete("1")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockMQ.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}
