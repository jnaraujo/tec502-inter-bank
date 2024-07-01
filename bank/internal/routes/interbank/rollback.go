package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
)

type rollbackBodySchema struct {
	TxId uuid.UUID `json:"tx_id" validate:"required"`
	Step string    `json:"step" validate:"required,oneof=credit debit"`
}

func RollbackRoute(c *fiber.Ctx) error {
	var body rollbackBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	transaction := storage.Transactions.FindTransactionById(body.TxId)
	if transaction == nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Transação não encontrada",
		})
	}

	if len(transaction.Operations) != 1 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "A transação deve ter apenas uma operação",
		})
	}

	operation := transaction.Operations[0]

	fromAcc := storage.Accounts.FindAccountByIBK(operation.From)
	toAcc := storage.Accounts.FindAccountByIBK(operation.To)

	if body.Step == "debit" && fromAcc == nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Conta de origem não encontrada",
		})
	}

	if body.Step == "credit" && toAcc == nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Conta de destino não encontrada",
		})
	}

	if operation.Status == models.OperationStatusSuccess {
		if body.Step == "debit" {
			storage.Accounts.AddToAccountBalance(fromAcc.Id, operation.Amount)
		} else if body.Step == "credit" {
			storage.Accounts.SubFromAccountBalance(toAcc.Id, operation.Amount)
		}
	}

	if operation.Status == models.OperationStatusPending {
		if body.Step == "debit" {
			storage.Accounts.SubFromBlockedAccountBalance(fromAcc.Id, operation.Amount)
		} else if body.Step == "credit" {
			storage.Accounts.SubFromPendingAccountBalance(toAcc.Id, operation.Amount)
		}
	}

	storage.Transactions.UpdateOperationStatus(*transaction, operation, models.OperationStatusFailed)
	transaction = storage.Transactions.UpdateTransactionStatus(*transaction, models.TransactionStatusFailed)

	return c.Status(http.StatusOK).JSON(transaction)
}
