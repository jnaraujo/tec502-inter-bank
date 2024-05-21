package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetRootRoute(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Bank App",
	})
}
