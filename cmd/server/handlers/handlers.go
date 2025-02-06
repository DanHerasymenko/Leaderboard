package handlers

import (
	_ "Leaderboard/cmd/server/docs"
	"Leaderboard/cmd/server/handlers/auth"
	"Leaderboard/cmd/server/handlers/counter"
	"Leaderboard/cmd/server/handlers/health"
	"Leaderboard/cmd/server/handlers/score"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handlers struct {
	Health  *health.Handler
	Counter *counter.Handler
	Auth    *auth.Handler
	Score   *score.Handler
}

func NewHandlers(cfg *config.Config, clnts *client.Clients, srvs *services.Services) *Handlers {
	return &Handlers{
		Health:  health.NewHandler(cfg),
		Counter: counter.NewHandler(cfg, clnts, srvs),
		Auth:    auth.NewHandler(cfg, clnts, srvs),
		Score:   score.NewHandler(cfg, clnts, srvs),
	}
}

func (h *Handlers) RegisterRoutes(router fiber.Router) {

	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/health", h.Health.Health)

	cg := router.Group("/counter")
	cg.Get("/increment", h.Counter.Increment)

	ag := router.Group("/auth")
	ag.Post("/singup", h.Auth.SingUp)
	ag.Post("/singin", h.Auth.SingIn)

	sg := router.Group("/score")
	sg.Post("/submit", h.Score.SubmitScore)
	sg.Post("/list", h.Score.ListScores)

}
