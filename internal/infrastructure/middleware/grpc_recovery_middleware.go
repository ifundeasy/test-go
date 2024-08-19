package middleware

import (
	"context"
	"runtime/debug"

	"test-go/internal/infrastructure/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryRecoveryInterceptor is a gRPC interceptor that recovers from panics and logs the error
func UnaryRecoveryInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Use recover to catch any panic that occurs during the handling of the request
		defer func() {
			if r := recover(); r != nil {
				// Log the panic and stack trace
				logger.Error("Recovered from panic: " + logPanic(r))
				logger.Error("Stack trace: " + string(debug.Stack()))

				// Convert the panic to a gRPC error
				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		// Continue handling the request
		return handler(ctx, req)
	}
}

// logPanicInterceptor formats the panic information into a string
func logPanicInterceptor(r interface{}) string {
	switch v := r.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "Unknown panic"
	}
}
