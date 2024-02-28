package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bugrakocabay/konsume/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBPlugin struct {
	db *mongo.Database
}

// Connect establishes a connection to the MongoDB database
func (m *MongoDBPlugin) Connect(connectionString, dbName string) error {
	slog.Info("Connecting to MongoDB database")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	m.db = client.Database(dbName)
	slog.Info("Connected to the MongoDB database")
	return nil
}

// Insert stores data into the MongoDB database
func (m *MongoDBPlugin) Insert(data map[string]interface{}, dbRouteConfig config.DatabaseRouteConfig) error {
	transformedData := make(map[string]interface{})
	for key, value := range data {
		if newKey, ok := dbRouteConfig.Mapping[key]; ok {
			transformedData[newKey] = value
		} else {
			transformedData[key] = value
		}
	}
	collection := m.db.Collection(dbRouteConfig.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, transformedData)
	if err != nil {
		slog.Error("Failed to insert data into MongoDB", "error", err, "collection", dbRouteConfig.Collection)
		return fmt.Errorf("error inserting data into MongoDB collection %s: %w", dbRouteConfig.Collection, err)
	}

	slog.Info("Data inserted into MongoDB collection successfully", "collection", dbRouteConfig.Collection)
	return nil
}

// Close closes the connection to the MongoDB database
func (m *MongoDBPlugin) Close() error {
	slog.Info("Closing connection to MongoDB database")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.db.Client().Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

var Plugin MongoDBPlugin
