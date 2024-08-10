package mock

import (
	"context"
	"test-go/internal/product"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoCollection is a mock type for the MongoDB Collection
type MockMongoCollection struct {
	mock.Mock
}

func (m *MockMongoCollection) Fetch() ([]product.Product, error) {
	args := m.Called()
	return args.Get(0).([]product.Product), args.Error(1)
}

func (m *MockMongoCollection) GetByID(id string) (*product.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockMongoCollection) Create(prod *product.Product) error {
	args := m.Called(prod)
	return args.Error(0)
}

func (m *MockMongoCollection) Update(id string, prod *product.Product) error {
	args := m.Called(id, prod)
	return args.Error(0)
}

func (m *MockMongoCollection) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Additional MongoDB specific methods if needed for testing
func (m *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}
