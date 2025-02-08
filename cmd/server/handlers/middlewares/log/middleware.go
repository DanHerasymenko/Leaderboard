package log

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/logger"
	"Leaderboard/internal/services"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"time"
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

	logger.Info(ctx.Context(), "request start", slog.String("request_id", reqID), slog.String("request_path", reqPath))

	startedAt := time.Now()
	err := ctx.Next()
	duration := time.Since(startedAt)

	fe := &fiber.Error{}
	if errors.As(err, &fe) {
		logger.Error(ctx.Context(), err, slog.String("request_id", reqID), slog.String("request_path", reqPath), slog.String("error", fe.Message))
		return err
	}

	logger.Info(ctx.Context(),
		"request end",
		slog.String("request_id", reqID),
		slog.String("request_path", reqPath),
		slog.Int64("duration", duration.Milliseconds()),
		slog.Int("status", ctx.Response().StatusCode()))

	return err

}
