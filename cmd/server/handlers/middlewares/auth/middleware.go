package auth

import (
	fu "Leaderboard/cmd/server/utils/fiber"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/constants"
	"Leaderboard/internal/services"
	"errors"
	"github.com/gofiber/fiber/v2"
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

const (
	authHeaderName = "X-User-Token"
	ctxUserKey     = "user_id"
	ctxRoleKey     = "role"
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

	role, err := m.svcs.Auth.GetUserRole(ctx.Context(), userID)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals(ctxUserKey, *userID)
	ctx.Locals(ctxRoleKey, role)

	fu.SetLoggerAttr(ctx, slog.String(ctxUserKey, *userID), slog.String(ctxRoleKey, role))

	return ctx.Next()
}

func (m *Middleware) HandleIsAdmin(ctx *fiber.Ctx) error {
	userRole, ok := ctx.Locals("role").(string)
	if !ok || userRole != constants.AdminRole {
		return fiber.ErrForbidden
	}
	return ctx.Next()
}

func mustBeUserID(v any) string {
	if uid, ok := v.(string); ok {
		return uid
	}
	panic(ErrUserIdExpected)
}

func MustGetUserID(ctx *fiber.Ctx) (error, string) {
	val := ctx.Locals(ctxUserKey)

	if val == nil {
		return fiber.ErrUnauthorized, ""
	}

	return nil, mustBeUserID(val)

}

//func MustGetUserWs(conn *websocket.Conn) *as.User {
//	val := conn.Locals(ctxKey)
//
//	return mustBeUser(val)
//}
