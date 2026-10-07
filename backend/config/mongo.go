package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	// Command Monitor
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			log.Printf("[MongoDB] ▶ Command Started\n  ↪︎ Type: %s\n  ↪︎ Command: %v\n", evt.CommandName, evt.Command)
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			log.Printf("[MongoDB] ❌ Command Failed\n  ↪︎ Type: %s\n  ↪︎ Failure: %v\n", evt.CommandName, evt.Failure)
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			log.Printf("[MongoDB] ✅ Command Succeeded\n  ↪︎ Type: %s\n  ↪︎ Duration: %v\n", evt.CommandName, evt.Duration)
		},
	}

	// MongoDB Connection Options
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetMonitor(monitor)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("[MongoDB] ❌ Connection Error: %v", err)
	}

	DB = client.Database("fixedasset-golang")
	log.Println("✅ Connected to MongoDB: fixedasset-golang")
}

func GetCollection(name string) *mongo.Collection {
	log.Printf("[MongoDB] 📦 Accessing collection: %s", name)
	return DB.Collection(name)
}
