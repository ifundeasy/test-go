package db

import (
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var mongoDatabase *mongo.Database
var mongoOnce sync.Once

func GetMongoInstance(mongoURI string, dBName string) *mongo.Database {
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Check the connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		mongoDatabase = client.Database(dBName)
		log.Println("Connected to MongoDB successfully")
	})

	return mongoDatabase
}

func CloseMongoInstance(ctx context.Context) error {
	if mongoDatabase == nil {
		return nil
	}
	return mongoDatabase.Client().Disconnect(ctx)
}
