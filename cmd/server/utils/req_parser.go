package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
)

// доволі масивна штука - там багато оптимізацій, кешування і т.д.
// правильно не створювати на кожен чих новий валідатор

var v = validator.New()

var seasonRegex = regexp.MustCompile(`^(Winter|Spring|Summer|Autumn)[0-9]{4}$`)

func ParseReqBody(ctx *fiber.Ctx, reqBody interface{}) error {
	if err := ctx.BodyParser(reqBody); err != nil {
		return fiber.ErrBadRequest
	}

	if err := v.Struct(reqBody); err != nil {
		return fmt.Errorf("failed to validate request body: %w", err)
	}
	return nil
}
