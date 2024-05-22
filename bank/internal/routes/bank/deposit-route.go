package bank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type depositBodySchema struct {
	UserId int             `json:"user_id" validate:"required"`
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

func DepositRoute(c *fiber.Ctx) error {
	var body depositBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	user, exists := storage.FindUserById(body.UserId)
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	if body.Amount.LessThan(decimal.NewFromInt(0)) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Amount must be greater than 0",
		})
	}

	transaction := storage.CreateDepositTransaction(user.InterBankKey, body.Amount, models.TransactionTypeDeposit)

	user, ok := storage.AddToUserBalance(user.Id, body.Amount)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error adding to user balance",
		})
	}

	return c.Status(http.StatusOK).JSON(&transaction)
}
