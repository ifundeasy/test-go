package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MockGRPCClientStream is a mock type for the gRPC ClientStream
type MockGRPCClientStream struct {
	mock.Mock
	grpc.ClientStream
}

// MockGRPCServerStream is a mock type for the gRPC ServerStream
type MockGRPCServerStream struct {
	mock.Mock
	grpc.ServerStream
}

// SendMsg is a mock implementation of the gRPC SendMsg method for ClientStream
func (m *MockGRPCClientStream) SendMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

// RecvMsg is a mock implementation of the gRPC RecvMsg method for ClientStream
func (m *MockGRPCClientStream) RecvMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

// SendMsg is a mock implementation of the gRPC SendMsg method for ServerStream
func (m *MockGRPCServerStream) SendMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

// RecvMsg is a mock implementation of the gRPC RecvMsg method for ServerStream
func (m *MockGRPCServerStream) RecvMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

// SetHeader is a mock implementation of the SetHeader method for ServerStream
func (m *MockGRPCServerStream) SetHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

// SendHeader is a mock implementation of the SendHeader method for ServerStream
func (m *MockGRPCServerStream) SendHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

// SetTrailer is a mock implementation of the SetTrailer method for ServerStream
func (m *MockGRPCServerStream) SetTrailer(md metadata.MD) {
	m.Called(md)
}

// Context is a mock implementation of the Context method for ServerStream
func (m *MockGRPCServerStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// SendMsg is a mock implementation of the SendMsg method for both ClientStream and ServerStream
func (m *MockGRPCClientStream) SendHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}
