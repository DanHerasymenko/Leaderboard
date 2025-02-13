package services

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services/auth"
	"Leaderboard/internal/services/score"
)

type Services struct {
	Auth  *auth.Service
	Score *score.Service
}

func NewServices(cfg *config.Config, clnts *client.Clients) *Services {
	return &Services{
		Auth:  auth.NewService(cfg, clnts),
		Score: score.NewService(cfg, clnts),
	}
}
