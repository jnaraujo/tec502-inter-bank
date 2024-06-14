package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func ListAccountTransactionsRoute(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid account id",
		})
	}

	_, exists := storage.Accounts.FindUserById(id)
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	transactions := storage.Transactions.FindUserTransactionsById(id)
	return c.Status(http.StatusOK).JSON(&transactions)
}
