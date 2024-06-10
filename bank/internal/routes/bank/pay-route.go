package bank

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank/service"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/utils"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type payRouteBodySchema struct {
	From   interbank.IBK   `json:"from_user_ibk" validate:"required"`
	To     interbank.IBK   `json:"to_user_ibk" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

func PayRoute(c *fiber.Ctx) error {
	var body payRouteBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	if !utils.IsLocalUserIBK(body.From) {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": "Sender must be from this bank",
		})
	}

	if body.From == body.To {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Sender and receiver must be different",
		})
	}

	if body.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Amount must be greater than 0",
		})
	}

	transaction := *models.NewTransaction(body.From, []models.Operation{
		*models.NewOperation(
			body.From,
			body.To,
			models.OperationTypeTransfer,
			body.Amount,
		),
	})
	storage.Transactions.Save(transaction)

	service.BeginTransaction(int(body.From.BankId))
	service.BeginTransaction(int(body.To.BankId))

	defer func() {
		service.CommitTransaction(int(body.From.BankId))
		service.CommitTransaction(int(body.To.BankId))
	}()

	err := service.SubCredit(int(body.From.BankId), body.From, body.Amount)
	if err != nil {
		fmt.Println("sub")
		storage.Transactions.UpdateOperationStatus(transaction, transaction.Operations[0], models.OperationStatusFailed)
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)

		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	err = service.AddCredit(int(body.To.BankId), body.To, body.Amount)
	if err != nil {
		fmt.Println("add")
		storage.Transactions.UpdateOperationStatus(transaction, transaction.Operations[0], models.OperationStatusFailed)
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)

		// Rollback
		service.AddCredit(int(body.From.BankId), body.From, body.Amount)

		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
}
