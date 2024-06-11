package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank/service"
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

		operations = append(operations, *models.NewOperation(op.From, op.To, models.OperationTypeTransfer, op.Amount))
	}

	transaction := *models.NewTransaction(body.Author, operations)
	storage.Transactions.Save(transaction)

	service.LockAccountsFromTransaction(transaction)
	defer func() {
		service.UnlockAccountsFromTransaction(transaction)
	}()

	for _, op := range transaction.Operations {
		err := service.SubCredit(int(op.From.BankId), op.From, op.Amount)
		if err != nil {
			service.RollbackOperations(transaction)
			return c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"message": err.Error(),
			})
		}

		err = service.AddCredit(int(op.To.BankId), op.To, op.Amount)
		if err != nil {
			// como falou na segunda parte, reverte a primeira parte
			service.AddCredit(int(op.From.BankId), op.From, op.Amount)
			service.RollbackOperations(transaction)
			return c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"message": err.Error(),
			})
		}

		storage.Transactions.UpdateOperationStatus(transaction, op, models.OperationStatusSuccess)
	}

	storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
}
