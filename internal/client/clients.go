package client

import (
	"Leaderboard/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Clients struct {
	Mongo *mongo.Database
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {

	mongoOpts := options.Client().ApplyURI(cfg.MongoURI)

	mongoClient, err := mongo.Connect(ctx, mongoOpts)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &Clients{
		Mongo: mongoClient.Database(cfg.MongoDbName),
	}, nil

}
