package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func UnlockRoute(c *fiber.Ctx) error {
	storage.Accounts.Unlock()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
