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

type depositBodySchema struct {
	IBK    interbank.IBK   `json:"acc_ibk" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

func DepositRoute(c *fiber.Ctx) error {
	var body depositBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"errors": errs})
	}

	acc := storage.Accounts.FindAccountByIBK(body.IBK)
	if acc == nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "Conta não encontrada",
		})
	}

	if body.Amount.LessThan(decimal.NewFromInt(0)) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Valor deve ser maior que zero",
		})
	}

	transaction := storage.Transactions.CreateDepositTransaction(acc.InterBankKey, body.Amount)

	ok := storage.Accounts.AddToAccountBalance(acc.Id, body.Amount)
	if !ok {
		storage.Transactions.UpdateOperationStatus(transaction, transaction.Operations[0], models.OperationStatusFailed)
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Erro ao adicionar valor na conta",
		})
	}

	storage.Transactions.UpdateOperationStatus(transaction, transaction.Operations[0], models.OperationStatusSuccess)
	transaction = *storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusOK).JSON(&transaction)
}
