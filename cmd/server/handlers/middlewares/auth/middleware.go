package auth

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/logger"
	"Leaderboard/internal/services"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"

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

const (
	authHeaderName = "X-User-Token"
	ctxKey         = "user_id"
)

var ErrUserIdExpected = errors.New("user id expected")

func (m *Middleware) Handle(ctx *fiber.Ctx) error {

	token := ctx.Get(authHeaderName)

	if token == "" {
		return fiber.ErrUnauthorized
	}

	userID, err := m.svcs.Auth.VerifyAuthToken(token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals(ctxKey, *userID)

	fu.SetLoggerAttr(ctx, slog.String(ctxKey, *userID))

	return ctx.Next()
}

func MustGetUserID(ctx *fiber.Ctx) string {
	val := ctx.Locals(ctxKey)

	if uid, ok := val.(string); ok {
		return uid
	}

	logger.Panic(ctx.Context(), ErrUserIdExpected)
	return ""

}
