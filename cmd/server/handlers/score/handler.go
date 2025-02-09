package score

import (
	rpu "Leaderboard/cmd/server/utils/req_parser"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	ss "Leaderboard/internal/services/score"
	"github.com/gofiber/fiber/v2"
	"regexp"
)

var seasonRegex = regexp.MustCompile(`^(Winter|Spring|Summer|Autumn)[0-9]{4}$`)

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

type SubmitScoreReqBody struct {
	NickName string `json:"nickName" validate:"required,min=3,max=32,alphanum"`
	Rating   int    `json:"rating" validate:"required,min=1"`
	Wins     int    `json:"wins" validate:"required,min=0"`
	Losses   int    `json:"losses" validate:"required,min=0"`
	Region   string `json:"region" validate:"required,oneof=EU NA AS SA AF"`
}

type SubmitScoreResp200Body struct {
	Score *ss.Score `json:"score"`
}

// @Summary SubmitScore
// @Description SubmitScore
// @Tags Score
// @Param body body SubmitScoreReqBody true "SubmitScore request body"
// @Success 200 {object} SubmitScoreResp200Body "SubmitScore success"
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/submit [post]
// @Security UserTokenAuth
func (h *Handler) SubmitScore(ctx *fiber.Ctx) error {

	reqBody := &SubmitScoreReqBody{}

	if err := rpu.ParseReqBody(ctx, reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	_, err := h.svsc.Auth.GetUserByName(ctx.Context(), reqBody.NickName)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	score, err := h.svsc.Score.AddNewScore(ctx.Context(), reqBody.Rating, reqBody.NickName, reqBody.Wins, reqBody.Losses, reqBody.Region)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(&SubmitScoreResp200Body{Score: score})

}

type ListScoresReqBody struct {
	Season         string  `json:"season" validate:"required"`
	LastReceivedId *string `json:"last_received_id"`
}

type ListScoresResp200Body struct {
	Scores []ss.Score `json:"scores"`
}

// @Summary ListScores
// @Description ListScores
// @Tags Score
// @Param body body ListScoresReqBody true "ListScores request body"
// @Success 200 {array} Score "ListScores success"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/list [post]
// @Security UserTokenAuth
func (h *Handler) ListScores(ctx *fiber.Ctx) error {

	reqBody := &ListScoresReqBody{}

	if err := rpu.ParseReqBody(ctx, reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	if !seasonRegex.MatchString(reqBody.Season) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid season format (expected e.g. 'Winter2025')")
	}

	scoreList, err := h.svsc.Score.ListScores(ctx.Context(), ss.ListScoresParams{
		Season:              reqBody.Season,
		Limit:               2,
		LastReceivedScoreID: reqBody.LastReceivedId,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(ListScoresResp200Body{Scores: scoreList})
}
