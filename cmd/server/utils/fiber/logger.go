package fiber

import (
	"Leaderboard/internal/logger"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func getLoggerAttr(ctx *fiber.Ctx) []slog.Attr {
	av := ctx.Locals(logger.CtxValueKey{})
	if av == nil {
		return []slog.Attr{}
	}

	res, ok := av.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}
	return res
}

func mergeAttrs(ctx *fiber.Ctx, attrs []slog.Attr) []slog.Attr {
	existing := getLoggerAttr(ctx)
	return append(existing, attrs...)
}

func SetLoggerAttr(ctx *fiber.Ctx, attrs ...slog.Attr) {
	attr := mergeAttrs(ctx, attrs)
	ctx.Locals(logger.CtxValueKey{}, attr)
}
