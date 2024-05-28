package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type createAccountBody struct {
	Name  string `json:"name" validate:"required,lte=255"`
	Email string `json:"email" validate:"required,lte=255,email"`
}

func CreateAccountRoute(c *fiber.Ctx) error {
	var body createAccountBody
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	_, exists := storage.Users.FindUserByEmail(body.Email)
	if exists {
		return c.Status(http.StatusConflict).JSON(&fiber.Map{
			"error": "user already exists",
		})
	}

	user := storage.Users.CreateUser(body.Name, body.Email)
	return c.Status(http.StatusCreated).JSON(&user)
}
