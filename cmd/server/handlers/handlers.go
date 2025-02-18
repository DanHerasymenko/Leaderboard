package handlers

import (
	_ "Leaderboard/cmd/server/docs"
	"Leaderboard/cmd/server/handlers/auth"
	"Leaderboard/cmd/server/handlers/health"
	"Leaderboard/cmd/server/handlers/middlewares"
	"Leaderboard/cmd/server/handlers/score"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handlers struct {
	Health *health.Handler
	Auth   *auth.Handler
	Score  *score.Handler

	mdlwrs *middlewares.Middleware
}

func NewHandlers(cfg *config.Config, clnts *client.Clients, srvs *services.Services, mdlwrs *middlewares.Middleware) *Handlers {
	return &Handlers{
		Health: health.NewHandler(cfg),
		Auth:   auth.NewHandler(cfg, clnts, srvs),
		Score:  score.NewHandler(cfg, clnts, srvs),

		mdlwrs: mdlwrs,
	}
}

func (h *Handlers) RegisterRoutes(router fiber.Router) {

	//Swagger
	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/health", h.Health.Health)

	//API
	api := router.Group("/api")
	api.Use(h.mdlwrs.Log.Handle)

	ag := api.Group("/auth")
	ag.Post("/singup", h.Auth.SingUp)
	ag.Post("/singin", h.Auth.SingIn)

	sg := api.Group("/score")
	sg.Post("/submit", h.mdlwrs.Auth.Handle, h.Score.SubmitScore)
	sg.Post("/list", h.Score.ListScores)
	sg.Get("/top", h.Score.GetTopScores)
	sg.Delete("/delete", h.mdlwrs.Auth.Handle, h.mdlwrs.Auth.HandleIsAdmin, h.Score.DeleteAllScores)

	//Websocket
	sg.Get("/listen_list", websocket.New(h.Score.ListenScores))

	//Web pages
	web := router.Group("/leaderboard")
	web.Use(h.mdlwrs.Log.Handle)

	list := web.Group("/list")
	//list.Use(h.mdlwrs.Auth.Handle)
	list.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./web/list.html")
	})
}
