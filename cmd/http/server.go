package main

import (
	"log"

	"test-go/internal/adapters/primary/http"
	_ "test-go/internal/adapters/primary/http/swagger"
	queue "test-go/internal/adapters/secondary/messaging"
	"test-go/internal/application"
	"test-go/internal/infrastructure/config"
	"test-go/internal/infrastructure/db"
	"test-go/internal/infrastructure/logging"
	"test-go/internal/infrastructure/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Product API
// @version 1.0
// @description API for managing products
// @host localhost:3002
// @BasePath /

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// Load configuration
	conf := config.LoadConfig()

	mongoDB := db.GetMongoInstance(conf.MongoURI, conf.MongoDbName)
	redisClient := db.GetRedisInstance(conf.RedisURI, conf.RedisDbName)
	rabbitMQ := queue.GetRmqInstance(conf.RabbitMqURI)

	// Initialize the logger
	logger := logging.NewLogger("HTTP: ")

	// Create a new Fiber app
	app := fiber.New()

	// Apply middleware
	app.Use(middleware.RecoveryMiddleware(logger)) // Handle panics and log them
	app.Use(middleware.LoggingMiddleware(logger))  // Custom logging middleware for detailed logs
	// app.Use(middleware.AuthMiddleware)             // Authentication middleware

	// Create service and handler, passing MongoDB and Redis clients to the service
	productService := application.NewProductService(mongoDB, redisClient, rabbitMQ)
	productHandler := http.NewProductHandler(productService)

	// Set up routes
	http.SetupRoutes(app, productHandler)

	// Register Swagger route
	app.Get("/api/docs/*", swagger.HandlerDefault) // Swagger endpoint

	// Start the HTTP server
	log.Printf("HTTP server is running on port %s", conf.HttpPort)
	if err := app.Listen(":" + conf.HttpPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
