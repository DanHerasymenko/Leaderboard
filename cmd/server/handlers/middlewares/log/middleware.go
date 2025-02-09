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

	fu "Leaderboard/cmd/server/utils/fiber"
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

	fu.SetLoggerAttr(ctx, slog.String("request_id", reqID), slog.String("request_path", reqPath))

	logger.Info(ctx.Context(), "request start")

	startedAt := time.Now()
	err := ctx.Next()
	duration := time.Since(startedAt)

	fu.SetLoggerAttr(ctx, slog.Int64("duration", duration.Milliseconds()))

	fe := &fiber.Error{}
	if errors.As(err, &fe) {
		fu.SetLoggerAttr(ctx, slog.String("resp_status", fe.Message))
	} else if err != nil {
		fu.SetLoggerAttr(ctx, slog.String("resp_message", err.Error()))
	}

	fu.SetLoggerAttr(ctx, slog.Int("status", ctx.Response().StatusCode()))

	logger.Info(ctx.Context(), "request end")

	return err

}
