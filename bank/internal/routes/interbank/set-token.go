package interbank

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type setTokenBodySchema struct {
	To interbank.BankId `json:"to" validate:"required"`
	Ts time.Time        `json:"ts" validate:"required"`
}

func SetToken(c *fiber.Ctx) error {
	var body setTokenBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	storage.Token.Set(token.Token{
		Owner: body.To,
		Ts:    body.Ts,
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "token set",
	})
}
