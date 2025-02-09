package middlewares

import (
	"Leaderboard/cmd/server/handlers/middlewares/auth"
	"Leaderboard/cmd/server/handlers/middlewares/log"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
)

type Middleware struct {
	Log  *log.Middleware
	Auth *auth.Middleware
}

func NewMiddleware(cfg *config.Config, clnts *client.Clients, svcs *services.Services) *Middleware {
	return &Middleware{
		Log:  log.NewMiddleware(cfg, clnts, svcs),
		Auth: auth.NewMiddleware(cfg, clnts, svcs),
	}
}
