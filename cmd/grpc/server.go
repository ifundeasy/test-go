package main

import (
	"net"

	grpcHandler "test-go/internal/adapters/primary/grpc"
	"test-go/internal/adapters/primary/grpc/proto"
	queue "test-go/internal/adapters/secondary/messaging"
	"test-go/internal/application"
	"test-go/internal/infrastructure/config"
	"test-go/internal/infrastructure/db"
	"test-go/internal/infrastructure/logging"
	"test-go/internal/infrastructure/middleware"

	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	conf := config.LoadConfig()

	mongoDB := db.GetMongoInstance(conf.MongoURI, conf.MongoDbName)
	redisClient := db.GetRedisInstance(conf.RedisURI, conf.RedisDbName)
	rabbitMQ := queue.GetRmqInstance(conf.RabbitMqURI)

	// Initialize the logger
	logger := logging.NewLogger("gRPC: ")

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryLoggingInterceptor(logger),  // Logging interceptor
			middleware.UnaryRecoveryInterceptor(logger), // Recovery interceptor
		),
	)

	// Initialize ProductService and ProductHandler
	productService := application.NewProductService(mongoDB, redisClient, rabbitMQ)
	productHandler := grpcHandler.NewProductHandler(productService)

	// Register the ProductService server
	proto.RegisterProductServiceServer(grpcServer, productHandler)

	// Listen on the specified gRPC port
	listener, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		logger.Fatal("Failed to start gRPC server: " + err.Error())
	}

	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server: " + err.Error())
	}
}
