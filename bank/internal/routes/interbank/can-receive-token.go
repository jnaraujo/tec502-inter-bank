package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CanReceiveToken(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
