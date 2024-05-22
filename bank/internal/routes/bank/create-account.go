package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type createAccountBody struct {
	Name string `json:"name" validate:"required,lte=255"`
}

func CreateAccountRoute(c *fiber.Ctx) error {
	var body createAccountBody
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	user := storage.CreateUser(body.Name)
	return c.Status(http.StatusCreated).JSON(&user)
}
