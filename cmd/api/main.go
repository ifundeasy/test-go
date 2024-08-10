// @title Test-Go API
// @version 1.0
// @description This is a sample server for a Go hexagonal architecture.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

package main

import (
	"log"

	"test-go/internal/config"
	"test-go/internal/pkg/db"
	"test-go/internal/pkg/queue"
	"test-go/internal/product/delivery"
	"test-go/internal/product/repository"
	"test-go/internal/product/usecase"
	_ "test-go/internal/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
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

	// Initialize HTTP delivery
	app := fiber.New()

	// Route untuk dokumentasi Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	delivery.NewProductHandler(app, productUsecase)

	// Start the server
	log.Printf("Server is running on port %s", conf.HttpPort)
	if err := app.Listen(":" + conf.HttpPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
