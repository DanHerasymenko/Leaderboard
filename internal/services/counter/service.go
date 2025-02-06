package counter

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/constants"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	cfg   *config.Config
	clnts *client.Clients

	counterColl *mongo.Collection
}

func NewService(cfg *config.Config, clnts *client.Clients) *Service {
	return &Service{
		cfg:   cfg,
		clnts: clnts,

		counterColl: clnts.Mongo.Collection(constants.CounterMongoCollection),
	}
}

func (s *Service) Increment(ctx context.Context) (int, error) {
	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{"counter": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	res := s.counterColl.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", res.Err())
	}

	cntr := &Document{}
	err := res.Decode(&cntr)
	if err != nil {
		return 0, fmt.Errorf("failed to decode counter: %w", err)
	}

	return cntr.Counter, nil
}
