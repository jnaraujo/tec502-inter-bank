package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type commitBodySchema struct {
	TxId uuid.UUID `json:"tx_id" validate:"required"`
	Step string    `json:"step" validate:"required,oneof=credit debit"`
}

func CommitRoute(c *fiber.Ctx) error {
	var body commitBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"errors": errs})
	}

	transaction := storage.Transactions.FindTransactionById(body.TxId)
	if transaction == nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Conta não encontrada",
		})
	}

	if len(transaction.Operations) != 1 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Transação deve ter apenas uma operação",
		})
	}

	operation := transaction.Operations[0]
	if operation.Status == models.OperationStatusSuccess {
		transaction = storage.Transactions.UpdateTransactionStatus(*transaction, models.TransactionStatusFailed)
		return c.Status(http.StatusOK).JSON(transaction)
	}

	if body.Step == "debit" {
		fromAcc := storage.Accounts.FindAccountByIBK(operation.From)
		if fromAcc == nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Conta de origem não encontrada",
			})
		}

		storage.Accounts.SubFromBlockedAccountBalance(fromAcc.Id, operation.Amount)
		storage.Accounts.SubFromAccountBalance(fromAcc.Id, operation.Amount)
	} else if body.Step == "credit" {
		toAcc := storage.Accounts.FindAccountByIBK(operation.To)
		if toAcc == nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Conta de destino não encontrada",
			})
		}
		storage.Accounts.SubFromPendingAccountBalance(toAcc.Id, operation.Amount)
		storage.Accounts.AddToAccountBalance(toAcc.Id, operation.Amount)
	}

	storage.Transactions.UpdateOperationStatus(*transaction, operation, models.OperationStatusSuccess)
	transaction = storage.Transactions.UpdateTransactionStatus(*transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusOK).JSON(transaction)
}
