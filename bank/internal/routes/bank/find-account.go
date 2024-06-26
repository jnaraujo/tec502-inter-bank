package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAccountRoute(c *fiber.Ctx) error {
	accountId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid account id",
		})
	}

	user, exists := storage.Accounts.FindAccountById(accountId)
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	return c.Status(http.StatusOK).JSON(&user)
}
