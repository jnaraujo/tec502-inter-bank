package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func GetToken(c *fiber.Ctx) error {
	token := storage.Token.Get()
	if token.IsZero() {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "no token",
		})
	}
	return c.Status(http.StatusOK).JSON(token)
}
