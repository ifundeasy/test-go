package main

import (
	"log"
	"net"

	"test-go/internal/config"
	"test-go/internal/pkg/db"
	"test-go/internal/pkg/queue"
	"test-go/internal/product/delivery"
	"test-go/internal/product/pb"
	"test-go/internal/product/repository"
	"test-go/internal/product/usecase"
	_ "test-go/internal/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

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

	// Initialize HTTP server
	app := fiber.New()

	// Route untuk dokumentasi Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	delivery.NewProductHandler(app, productUsecase)

	// Initialize gRPC server and handler
	grpcServer := grpc.NewServer()
	productHandler := delivery.NewGRPCProductHandler(productUsecase)

	// Register the ProductServiceServer with the gRPC server
	pb.RegisterProductServiceServer(grpcServer, productHandler)

	// Start the HTTP server in a separate goroutine to prevent blocking
	go func() {
		log.Printf("HTTP server is running on port %s", conf.HttpPort)
		if err := app.Listen(":" + conf.HttpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start listening on the specified gRPC port
	listener, err := net.Listen("tcp", ":"+conf.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", conf.GrpcPort, err)
	}

	log.Printf("gRPC server is running on port %s", conf.GrpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
