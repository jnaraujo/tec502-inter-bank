package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CreateAccountRoute(c *fiber.Ctx) error {
	return c.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "Account created",
	})
}
