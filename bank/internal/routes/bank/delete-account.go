package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func DeleteAccountRoute(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid account id",
		})
	}

	storage.Accounts.RLock()
	_, exists := storage.Accounts.FindUserById(id)
	storage.Accounts.RUnlock()

	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "user not found",
		})
	}

	storage.Accounts.Lock()
	storage.Accounts.Delete(id)
	storage.Accounts.Unlock()

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "account deleted",
	})
}
