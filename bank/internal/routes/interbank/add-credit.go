package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type addCreditBodySchema struct {
	To     interbank.IBK   `json:"to" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

func AddCreditRoute(c *fiber.Ctx) error {
	var body addCreditBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	_, ok := storage.Accounts.AddToAccountBalance(int(body.To.UserId), body.Amount)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Conta não encontrada",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "OCrédito adicionado com sucesso",
	})
}
