package auth

import (
	rpu "Leaderboard/cmd/server/utils"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	as "Leaderboard/internal/services/auth"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg   *config.Config
	clnts *client.Clients
	svsc  *services.Services
}

func NewHandler(cfg *config.Config, clnts *client.Clients, srvs *services.Services) *Handler {
	return &Handler{
		cfg:   cfg,
		clnts: clnts,
		svsc:  srvs,
	}
}

type SingUpReqBody struct {
	NickName string `json:"nickName" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// @Summary SingUp
// @Description SingUp
// @Tags Auth
// @Param body body SingUpReqBody true "SingUp request body"
// @Success 200 {string} string "SingUp success"
// @Failure 409 {string} string "user already exists"
// @Router /api/auth/singup [post]
func (h *Handler) SingUp(ctx *fiber.Ctx) error {

	reqBody := &SingUpReqBody{}

	if err := rpu.ParseReqBody(ctx, reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	// Generate password hash
	passwordHash, err := h.svsc.Auth.GeneratePasswordHash(reqBody.Password)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	// Check if user already exists
	user, err := h.svsc.Auth.CreateUser(ctx.Context(), reqBody.NickName, passwordHash)

	if errors.Is(err, as.ErrUserAlreadyExists) {
		return fiber.ErrConflict
	}

	return ctx.JSON(user)
}

type SingInReqBody struct {
	NickName string `json:"nickName" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// @Summary SingIn
// @Description SingIn
// @Tags Auth
// @Param body body SingInReqBody true "SingIn request body"
// @Success 200 {string} string "SingIn success"
// @Failure 401 {string} string "invalid credentials"
// @Router /api/auth/singin [post]
func (h *Handler) SingIn(ctx *fiber.Ctx) error {

	reqBody := &SingUpReqBody{}

	if err := rpu.ParseReqBody(ctx, reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	user, err := h.svsc.Auth.GetUserByName(ctx.Context(), reqBody.NickName)
	if errors.Is(err, as.ErrUserNotFound) {
		return fiber.ErrUnauthorized
	} else if err != nil {
		return fmt.Errorf("failed to get user by name: %w", err)
	}

	ok, err := h.svsc.Auth.ComaprePasswordHashAndPassword(user.PasswordHash, reqBody.Password)
	if err != nil {
		return fmt.Errorf("failed to compare password hash and password: %w", err)
	}

	if !ok {
		return fiber.ErrUnauthorized
	}

	return ctx.JSON(user)
}
