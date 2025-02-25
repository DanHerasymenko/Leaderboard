package score

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/constants"
	"Leaderboard/internal/logger"
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
	eventList  *EventList
}

func NewService(cfg *config.Config, clnts *client.Clients) *Service {
	return &Service{
		cfg:        cfg,
		clnts:      clnts,
		scoresColl: clnts.Mongo.Collection(constants.ScoresMongoCollection),
		eventList:  StartEventList(),
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
				"worldRank":    0,
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

	opts := options.Update().SetUpsert(true)
	_, err = s.scoresColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert score: %w", err)
	}

	score := &Score{
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
	}

	err = s.updateWorldRanks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update rankings: %w", err)
	}

	s.eventList.PostScore(score)

	return score, nil
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

	opts := options.Find().
		SetSort(bson.D{{"ranking", -1}}).
		SetLimit(params.Limit)

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

func (s *Service) DeleteAllScores(ctx context.Context) error {

	_, err := s.scoresColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("failed to delete all scores: %w", err))
		return err
	}

	return nil
}

func (s *Service) updateWorldRanks(ctx context.Context) error {
	// MongoDB Aggregation to Rank Players by Rating
	pipeline := mongo.Pipeline{
		{{Key: "$setWindowFields", Value: bson.M{
			"partitionBy": nil,                               // Global ranking, no partitioning
			"sortBy":      bson.M{"scoreDetails.rating": -1}, // Descending order
			"output": bson.M{
				"scoreDetails.worldRank": bson.M{
					"$rank": bson.M{},
				},
			},
		}}},
	}

	// Apply the pipeline to update rankings
	_, err := s.scoresColl.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("failed to update rankings: %w", err)
	}

	return nil
}

func (s *Service) SubsribeList(l chan<- *Score) {
	s.eventList.Subscribe(l)
}

func (s *Service) UnsubscribeList(l chan<- *Score) {
	s.eventList.Unsubscribe(l)
}
