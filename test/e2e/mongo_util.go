package e2e

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongo(connectionString, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}

func createCollection(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.CreateCollection(ctx, collectionName)
	return err
}

func dropCollection(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.Collection(collectionName).Drop(ctx)
	return err
}

func queryCollection(db *mongo.Database, collectionName string) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
