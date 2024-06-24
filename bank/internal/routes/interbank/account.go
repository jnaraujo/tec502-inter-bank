package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAccountRoute(c *fiber.Ctx) error {
	document := c.Params("document")
	user, exists := storage.Accounts.FindUserByDocument(document)
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	return c.Status(http.StatusOK).JSON(&user)
}
