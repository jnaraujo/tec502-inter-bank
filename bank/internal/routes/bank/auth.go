package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type authBodySchema struct {
	Document string `json:"document" validate:"required,lte=255"`
}

func AuthRoute(c *fiber.Ctx) error {
	var body authBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	user, exists := storage.Accounts.FindUserByDocument(body.Document)
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "user does not exists",
		})
	}

	return c.Status(http.StatusCreated).JSON(&user)
}
