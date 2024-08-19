package middleware

import (
	"context"
	"time"

	"test-go/internal/infrastructure/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs details about the gRPC request and response for unary RPCs
func UnaryLoggingInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler to finish processing
		h, err := handler(ctx, req)

		// Log the details
		duration := time.Since(start)
		st, _ := status.FromError(err)

		logger.InfoJSON(map[string]interface{}{
			"method":   info.FullMethod,
			"duration": duration.String(),
			"status":   st.Code().String(),
			"error":    st.Message(),
		})

		return h, err
	}
}
