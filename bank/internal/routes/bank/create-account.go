package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type createAccountBody struct {
	Name     string `json:"name" validate:"required,lte=255"`
	Document string `json:"document" validate:"required,lte=255"`
}

func CreateAccountRoute(c *fiber.Ctx) error {
	var body createAccountBody
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	storage.Accounts.RLock()
	_, exists := storage.Accounts.FindUserByDocument(body.Document)
	storage.Accounts.RUnlock()

	if exists {
		return c.Status(http.StatusConflict).JSON(&fiber.Map{
			"error": "user already exists",
		})
	}

	storage.Accounts.Lock()
	user := storage.Accounts.CreateAccount(body.Name, body.Document)
	storage.Accounts.Unlock()

	return c.Status(http.StatusCreated).JSON(&user)
}
