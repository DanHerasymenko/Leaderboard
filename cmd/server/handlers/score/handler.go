package score

import (
	rpu "Leaderboard/cmd/server/utils/req_parser"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/services"
	as "Leaderboard/internal/services/auth"
	ss "Leaderboard/internal/services/score"
	"Leaderboard/internal/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"regexp"

	am "Leaderboard/cmd/server/handlers/middlewares/auth"
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
	//NickName string `json:"nickName" validate:"required,min=3,max=32,alphanum"`
	Rating int    `json:"rating" validate:"required,min=1"`
	Wins   int    `json:"wins" validate:"required,min=0"`
	Losses int    `json:"losses" validate:"required,min=0"`
	Region string `json:"region" validate:"required,oneof=EU NA AS SA AF"`
}

type SubmitScoreResp200Body struct {
	Score *ss.Score `json:"score"`
}

// @Summary SubmitScore
// @Description Create or updates if exists the player’s score
// @Tags Score
// @Param body body SubmitScoreReqBody true "SubmitScore request body"
// @Success 200 {object} SubmitScoreResp200Body "SubmitScore success"
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/submit [post]
// @Security UserTokenAuth
func (h *Handler) SubmitScore(ctx *fiber.Ctx) error {

	err, userID := am.MustGetUserID(ctx)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	reqBody := &SubmitScoreReqBody{}
	if err := rpu.ParseReqBody(ctx, reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	u, err := h.svsc.Auth.GetUserByParam(ctx.Context(), &as.UserSearchParameters{ID: &userID})
	if err != nil {
		return fiber.ErrUnauthorized
	}

	score, err := h.svsc.Score.AddNewScore(ctx.Context(), reqBody.Rating, u.Nickname, reqBody.Wins, reqBody.Losses, reqBody.Region)
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
// @Description Get list of scores based on season
// @Tags Score
// @Param body body ListScoresReqBody true "ListScores request body"
// @Success 200 {array} Score "ListScores success"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/list [post]
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
		Limit:               30,
		LastReceivedScoreID: reqBody.LastReceivedId,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(ListScoresResp200Body{Scores: scoreList})
}

// @Summary TopScores
// @Description Retrieves the top players based on their ranking
// @Tags Score
// @Param limit query int false "limit"
// @Default limit 3
// @Param season query string false "season"
// @Default season current season
// @Success 200 {array} Score "TopScores success"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/top [get]
func (h *Handler) GetTopScores(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit", 3)
	int64Limit := int64(limit)

	_, currentSeason := utils.GetCurrentSeasonAndTime()

	reqSeason := ctx.Query("season", currentSeason)

	if !seasonRegex.MatchString(reqSeason) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid season format (expected e.g. 'Winter2025')")
	}

	scoreList, err := h.svsc.Score.ListScores(ctx.Context(), ss.ListScoresParams{
		Season: reqSeason,
		Limit:  int64Limit,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(ListScoresResp200Body{Scores: scoreList})

}

// @Summary DeleteAllScores
// @Description Deletes the entire leaderboard (for admin use).
// @Tags Score
// @Success 200 {string} string "DeleteAllScores success"
// @Failure 500 {string} string "internal server error"
// @Router /api/score/delete [delete]
// @Security UserTokenAuth
func (h *Handler) DeleteAllScores(ctx *fiber.Ctx) error {

	err := h.svsc.Score.DeleteAllScores(ctx.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON("DeleteAllScores success")
}

// @Summary ListenScores
// @Description Listen to the score list updates
// @Tags Score
// @Success 200 {string} string "ListenScores success"
// @Router /api/score/listen_list [get]
func (h *Handler) ListenScores(conn *websocket.Conn) {

	ch := make(chan *ss.Score)
	h.svsc.Score.SubsribeList(ch)
	defer h.svsc.Score.UnsubscribeList(ch)

	for score := range ch {
		if err := conn.WriteJSON(score); err != nil {
			return
		}
	}

}
