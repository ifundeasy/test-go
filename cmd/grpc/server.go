package main

import (
	"log"
	"net"
	"os"

	"test-go/internal/config"
	"test-go/internal/pkg/db"
	"test-go/internal/pkg/queue"
	"test-go/internal/product/delivery"
	"test-go/internal/product/pb"
	"test-go/internal/product/repository"
	"test-go/internal/product/usecase"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// getEnv fetches an environment variable and provides a default value if not set
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()
	conf := config.LoadConfig()

	// Initialize MongoDB, Redis, and RabbitMQ as singletons
	mongoDatabase := db.GetMongoInstance(conf.MongoURI, conf.MongoDbName)
	redisClient := db.GetRedisInstance(conf.RedisURI, conf.RedisDbName)
	rabbitMQ := queue.GetRmqInstance(conf.RabbitMqURI)

	// Initialize repositories and use cases
	productRepo := repository.NewMongoProductRepository(mongoDatabase)
	productUsecase := usecase.NewProductUsecase(productRepo, redisClient, rabbitMQ)

	// Initialize gRPC server and handler
	grpcServer := grpc.NewServer()
	productHandler := delivery.NewGRPCProductHandler(productUsecase)

	// Register the ProductServiceServer with the gRPC server
	pb.RegisterProductServiceServer(grpcServer, productHandler)

	// Start listening on the specified port
	listener, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", conf.GrpcPort, err)
	}

	log.Printf("gRPC server is running on port %s", conf.GrpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
