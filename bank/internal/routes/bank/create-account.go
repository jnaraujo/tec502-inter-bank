package bank

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type createAccountBody struct {
	Name      string             `json:"name" validate:"required,lte=255"`
	Type      models.AccountType `json:"type" validate:"required,oneof=individual legal joint"`
	Documents []string           `json:"documents" validate:"required,min=1"`
}

func CreateAccountRoute(c *fiber.Ctx) error {
	var body createAccountBody
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	seen := make(map[string]bool)
	if body.Type == models.AccountTypeJoint {
		if len(body.Documents) < 2 {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "Contas conjuntas precisam de pelo menos dois documentos",
			})
		}

		for _, doc := range body.Documents {
			if seen[doc] {
				return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
					"error": "Documentos duplicados",
				})
			}
			seen[doc] = true
		}

		for _, doc := range body.Documents {
			_, exists := storage.Accounts.FindAccountByDocument(doc)
			if !exists {
				return c.Status(http.StatusConflict).JSON(&fiber.Map{
					"error": fmt.Sprintf("O documento %s não está associado a nenhuma conta", doc),
				})
			}
		}

		account := storage.Accounts.CreateJointAccount(body.Name, body.Documents)
		return c.Status(http.StatusBadRequest).JSON(&account)
	}

	if len(body.Documents) > 1 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "Contas individuais e jurídicas não podem ter mais de um documento",
		})
	}

	_, exists := storage.Accounts.FindAccountByDocument(body.Documents[0])
	if exists {
		return c.Status(http.StatusConflict).JSON(&fiber.Map{
			"error": "Conta já existe",
		})
	}

	user := storage.Accounts.CreateAccount(body.Name, body.Documents[0], body.Type)
	return c.Status(http.StatusCreated).JSON(&user)
}
