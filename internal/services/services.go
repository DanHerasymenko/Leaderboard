package services

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services/auth"
	"Leaderboard/internal/services/counter"
	"Leaderboard/internal/services/score"
)

type Services struct {
	Counter *counter.Service
	Auth    *auth.Service
	Score   *score.Service
}

func NewServices(cfg *config.Config, clnts *client.Clients) *Services {
	return &Services{
		Counter: counter.NewService(cfg, clnts),
		Auth:    auth.NewService(cfg, clnts),
		Score:   score.NewService(cfg, clnts),
	}
}
