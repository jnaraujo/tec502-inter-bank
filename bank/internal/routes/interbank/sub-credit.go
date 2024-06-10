package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type subCreditBodySchema struct {
	From   interbank.IBK   `json:"from" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

func SubCreditRoute(c *fiber.Ctx) error {
	var body subCreditBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	if !storage.Accounts.IsLocked() {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": "You must start a transaction.",
		})
	}

	err := storage.Accounts.SubFromUserBalance(int(body.From.UserId), body.Amount)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Success",
	})
}
