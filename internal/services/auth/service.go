package auth

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/constants"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cfg   *config.Config
	clnts *client.Clients

	usersColl *mongo.Collection
}

func NewService(cfg *config.Config, clnts *client.Clients) *Service {
	return &Service{
		cfg:       cfg,
		clnts:     clnts,
		usersColl: clnts.Mongo.Collection(constants.UsersMongoCollection),
	}
}

func (s *Service) GeneratePasswordHash(password string) (string, error) {

	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}

	return string(res), nil
}

func (s *Service) ComaprePasswordHashAndPassword(hash, password string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to compare password hash and password: %w", err)
	}

	return true, nil
}
