package storage

import (
	"context"
	"time"

	"github.com/ameer005/meowl/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection(uri string) (*mongo.Client, error) {
	logger.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		logger.Logger.Error("Mongodb failed to connect", "error", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Logger.Error("MongoDB ping failed", "err", err)
		return nil, err
	}

	logger.Logger.Info("Connected to MongoDB: ", "URL", uri)
	return client, nil
}
