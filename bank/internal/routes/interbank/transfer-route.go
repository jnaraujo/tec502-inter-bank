package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type transferBodySchema struct {
	FromUserIBK interbank.IBK   `json:"from_user_ibk" validate:"required"`
	ToUserIBK   interbank.IBK   `json:"to_user_ibk" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
}

func TransferRoute(c *fiber.Ctx) error {
	var body transferBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	return c.SendStatus(http.StatusNotImplemented)
}
