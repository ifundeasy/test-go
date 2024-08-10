package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort    string
	GrpcPort    string
	MongoURI    string
	MongoDbName string
	RedisURI    string
	RedisDbName string
	RabbitMqURI string
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		HttpPort:    getEnv("HTTP_PORT"),
		GrpcPort:    getEnv("GRPC_PORT"),
		MongoURI:    getEnv("MONGO_URI"),
		MongoDbName: getEnv("MONGO_DBNAME"),
		RedisURI:    getEnv("REDIS_URI"),
		RedisDbName: getEnv("REDIS_DBNAME"),
		RabbitMqURI: getEnv("RABBITMQ_URI"),
	}
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s not set", key)
	return ""
}

// Optional helper functions to parse other types from environment variables

// GetEnvAsInt reads an environment variable as integer or returns a default value if not set
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetEnvAsBool reads an environment variable as boolean or returns a default value if not set
func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
