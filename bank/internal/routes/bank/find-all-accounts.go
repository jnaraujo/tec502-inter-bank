package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAllAccountsRoute(c *fiber.Ctx) error {
	accountId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid account id",
		})
	}

	account, ok := storage.Accounts.FindAccountById(accountId)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	if account.Type == models.AccountTypeJoint {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": "Contas conjuntas não podem ver as contas de outros usuários",
		})
	}

	accounts := services.FindAllUserAccountsInterBank(account.Documents[0])
	return c.Status(http.StatusOK).JSON(&accounts)
}
