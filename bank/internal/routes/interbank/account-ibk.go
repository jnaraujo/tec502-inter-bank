package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAccountIBKRoute(c *fiber.Ctx) error {
	ibk, err := interbank.NewIBKFromStr(c.Params("ibk"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "IBK inválido",
		})
	}

	acc := storage.Accounts.FindAccountByIBK(*ibk)
	if acc == nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Conta não encontrada",
		})
	}

	return c.Status(http.StatusOK).JSON(&acc)
}
