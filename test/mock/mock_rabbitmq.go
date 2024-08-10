package mock

import (
	"github.com/stretchr/testify/mock"
)

// MockRabbitMQ is a mock type for the RabbitMQ
type MockRabbitMQ struct {
	mock.Mock
}

// PublishProductCreated is a mock implementation of the method that publishes a product created message
func (m *MockRabbitMQ) PublishProductCreated(product interface{}) error {
	args := m.Called(product)
	return args.Error(0)
}

// PublishProductUpdated is a mock implementation of the method that publishes a product updated message
func (m *MockRabbitMQ) PublishProductUpdated(product interface{}) error {
	args := m.Called(product)
	return args.Error(0)
}

// PublishProductDeleted is a mock implementation of the method that publishes a product deleted message
func (m *MockRabbitMQ) PublishProductDeleted(productID string) error {
	args := m.Called(productID)
	return args.Error(0)
}

// Close is a mock implementation of the method that closes the RabbitMQ connection
func (m *MockRabbitMQ) Close() error {
	args := m.Called()
	return args.Error(0)
}
