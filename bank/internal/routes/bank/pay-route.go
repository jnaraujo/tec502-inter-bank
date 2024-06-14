package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type operationSchema struct {
	From   interbank.IBK   `json:"from_user_ibk" validate:"required"`
	To     interbank.IBK   `json:"to_user_ibk" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

type payRouteBodySchema struct {
	Author     interbank.IBK     `json:"author" validate:"required"`
	Operations []operationSchema `json:"operations" validate:"required,min=1"`
}

func PayRoute(c *fiber.Ctx) error {
	var body payRouteBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	var operations []models.Operation
	for _, op := range body.Operations {
		// TODO: verificar se as operações são do mesmo usuário

		if op.From == op.To {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Sender and receiver must be different",
			})
		}

		if op.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Amount must be greater than 0",
			})
		}

		ops := models.NewTransferOperations(op.From, op.To, op.Amount)
		operations = append(operations, ops...)
	}

	transaction := *models.NewTransaction(body.Author, operations)
	storage.Transactions.Save(transaction)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
}
