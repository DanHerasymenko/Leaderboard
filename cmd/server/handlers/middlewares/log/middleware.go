package log

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/logger"
	"Leaderboard/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"log/slog"
)

type Middleware struct {
	cfg   *config.Config
	clnts *client.Clients
	svcs  *services.Services
}

func NewMiddleware(cfg *config.Config, clnts *client.Clients, svcs *services.Services) *Middleware {
	return &Middleware{
		cfg:   cfg,
		clnts: clnts,
		svcs:  svcs,
	}
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {

	reqID := ulid.Make().String()
	reqPath := ctx.Path()

	logger.Info(ctx.Context(), "request started", slog.String("request_id", reqID), slog.String("request_path", reqPath))

	return ctx.Next()
}
