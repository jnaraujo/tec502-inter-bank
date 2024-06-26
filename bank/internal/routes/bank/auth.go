package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type authBodySchema struct {
	IBK interbank.IBK `json:"user_ibk" validate:"required"`
}

func AuthRoute(c *fiber.Ctx) error {
	var body authBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	acc := storage.Accounts.FindAccountByIBK(body.IBK)
	if acc == nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "Conta não encontrada",
		})
	}

	return c.Status(http.StatusCreated).JSON(acc)
}
