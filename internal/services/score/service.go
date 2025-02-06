package score

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/constants"
	"Leaderboard/internal/utils"
	"context"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Service struct {
	cfg        *config.Config
	clnts      *client.Clients
	scoresColl *mongo.Collection
}

func NewService(cfg *config.Config, clnts *client.Clients) *Service {
	return &Service{
		cfg:        cfg,
		clnts:      clnts,
		scoresColl: clnts.Mongo.Collection(constants.ScoresMongoCollection),
	}
}

func (s *Service) AddNewScore(ctx context.Context, rating int, nickname string, wins int, losses int, region string) (*Score, error) {

	timeNow, currentSeason := utils.GetCurrentSeasonAndTime()

	filter := bson.M{"scoreDetails.nickname": nickname}

	var existingScore Score

	err := s.scoresColl.FindOne(ctx, filter).Decode(&existingScore)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to check existing score, no such nickname: %w", err)
	}

	// If player exists, keep the same ScoreID
	scoreID := ulid.MustNew(ulid.Timestamp(timeNow), ulid.DefaultEntropy()).String()
	if err == nil {
		scoreID = existingScore.ScoreID
	}

	// Prepare the updated fields
	update := bson.M{
		"$set": bson.M{
			"season":   currentSeason,
			"scoredAt": timeNow.Format(time.RFC3339),
			"scoreDetails": bson.M{
				"worldRank":    0, // TODO: Implement ranking position
				"rating":       rating,
				"nickname":     nickname,
				"wins":         wins,
				"losses":       losses,
				"winLoseRatio": utils.CalculateWinLoseRatio(wins, losses),
				"region":       region,
			},
		},
		"$setOnInsert": bson.M{
			"scoreID": scoreID, // Only set when inserting a new record
		},
	}

	// Perform upsert (update if exists, insert if not)
	opts := options.Update().SetUpsert(true)
	_, err = s.scoresColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert score: %w", err)
	}

	// Return the updated or inserted score
	return &Score{
		ScoreID:  scoreID,
		Season:   currentSeason,
		ScoredAt: timeNow.Format(time.RFC3339),
		ScoreDetails: Details{
			WorldRank:    0, // TODO: Implement ranking position
			Rating:       rating,
			Nickname:     nickname,
			Wins:         wins,
			Losses:       losses,
			WinLoseRatio: utils.CalculateWinLoseRatio(wins, losses),
			Region:       region,
		},
	}, nil
}

type ListScoresParams struct {
	Season              string
	Limit               int64
	LastReceivedScoreID *string
}

func (s *Service) ListScores(ctx context.Context, params ListScoresParams) ([]Score, error) {

	filter := bson.M{"season": params.Season}

	if params.LastReceivedScoreID != nil {
		filter["scoreID"] = bson.M{"$lt": *params.LastReceivedScoreID}
	}

	opts := options.Find().SetLimit(params.Limit)

	cur, err := s.scoresColl.Find(ctx, filter, opts)

	if err != nil {
		return nil, fmt.Errorf("failed to find scores: %w", err)
	}
	defer cur.Close(ctx)

	res := []Score{}

	if err = cur.All(ctx, &res); err != nil {
		return nil, fmt.Errorf("failed to decode scores: %w", err)
	}

	return res, nil

}
