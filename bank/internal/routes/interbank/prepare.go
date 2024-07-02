package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type prepareBodySchema struct {
	ParentId  models.TransactionId `json:"parent_id"`
	Operation struct {
		From   interbank.IBK   `json:"from" validate:"required"`
		To     interbank.IBK   `json:"to" validate:"required"`
		Amount decimal.Decimal `json:"amount" validate:"required"`
	} `json:"operation" validate:"required"`
	Step string `json:"step" validate:"required,oneof=credit debit"`
}

func PrepareRoute(c *fiber.Ctx) error {
	var body prepareBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"errors": errs})
	}

	fromAcc := storage.Accounts.FindAccountByIBK(body.Operation.From)
	toAcc := storage.Accounts.FindAccountByIBK(body.Operation.To)

	var owner interbank.IBK

	if body.Step == "debit" {
		if fromAcc == nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Conta não encontrada",
			})
		}

		if !storage.Accounts.CanSubFromAccount(fromAcc.Id, body.Operation.Amount) {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Saldo insuficiente",
			})
		}
		storage.Accounts.AddToBlockedAccountBalance(fromAcc.Id, body.Operation.Amount)
		owner = body.Operation.From
	} else { // credit
		if toAcc == nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Conta não encontrada",
			})
		}
		storage.Accounts.AddToPendingAccountBalance(toAcc.Id, body.Operation.Amount)
		owner = body.Operation.To
	}

	transaction := models.NewFinalTransaction(&body.ParentId, owner, []models.Operation{
		*models.NewOperation(body.Operation.From, body.Operation.To, models.OperationTypeTransfer, body.Operation.Amount),
	})

	storage.Transactions.Save(*transaction)

	return c.Status(http.StatusOK).JSON(transaction)
}
