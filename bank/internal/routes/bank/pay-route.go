package bank

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type operationSchema struct {
	From   interbank.IBK   `json:"from" validate:"required"`
	To     interbank.IBK   `json:"to" validate:"required"`
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

	author := storage.Accounts.FindAccountByIBK(body.Author)
	if author == nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Conta de autor não encontrada",
		})
	}

	// Um usuário pessoa física ou jurídica só pode fazer transferências com suas próprias contas
	// Uma conta conjunta só pode fazer transferências com a própria conta
	var userAccounts []models.Account
	if author.Type != models.AccountTypeJoint {
		userAccounts = services.FindAllUserAccountsInterBank(author.Documents[0])
	} else {
		userAccounts = []models.Account{*author}
	}

	var operations []models.Operation
	for _, op := range body.Operations {
		fromAcc := services.FindAccountInterBank(op.From)
		if fromAcc == nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": fmt.Sprintf("Conta de origem %s não encontrada", op.From),
			})
		}

		if !slices.ContainsFunc(userAccounts, func(acc models.Account) bool {
			return slices.Contains(fromAcc.Documents, acc.Documents[0])
		}) {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Usuário não pode fazer transferências com contas de terceiros",
			})
		}

		if op.From == op.To {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Conta de origem e destino não podem ser iguais",
			})
		}

		if op.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "Valor da operação deve ser maior que zero",
			})
		}

		op := *models.NewOperation(op.From, op.To, models.OperationTypeTransfer, op.Amount)
		operations = append(operations, op)
	}

	transaction := *models.NewTransaction(body.Author, operations)
	storage.Transactions.Save(transaction)
	storage.TransactionQueue.Add(transaction.Id)

	return c.Status(http.StatusCreated).JSON(&transaction)
}
