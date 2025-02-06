package counter

import (
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg   *config.Config
	clnts *client.Clients
	svsc  *services.Services
}

type CounterDoc struct {
}

type IncrementResp200Body struct {
	Counter int `json:"counter"`
}

func NewHandler(cfg *config.Config, clnts *client.Clients, srvs *services.Services) *Handler {
	return &Handler{
		cfg:   cfg,
		clnts: clnts,
		svsc:  srvs,
	}
}

// @Summary Increment counter
// @Description Increment counter by 1
// @Tags Counter
// @Success 200 {object} IncrementResp200Body
// @Router /counter/increment [get]
func (h *Handler) Increment(ctx *fiber.Ctx) error {

	counter, err := h.svsc.Counter.Increment(ctx.Context())
	if err != nil {

	}

	return ctx.JSON(&IncrementResp200Body{Counter: counter})
}
