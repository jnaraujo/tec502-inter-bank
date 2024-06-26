package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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

	account, ok := storage.Accounts.FindUserById(accountId)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	accounts := services.FindAllUserAccountsInterBank(account.Document)
	return c.Status(http.StatusOK).JSON(&accounts)
}
