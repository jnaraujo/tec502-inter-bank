package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func UnlockRoute(c *fiber.Ctx) error {
	if !storage.Accounts.IsLocked() {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "accounts are not locked",
		})
	}

	storage.Accounts.Unlock()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
