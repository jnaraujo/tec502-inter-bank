package interbank

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/utils"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/validate"
	"github.com/shopspring/decimal"
)

type transferBodySchema struct {
	FromUserIBK interbank.UserKey `json:"from_user_ibk" validate:"required"`
	ToUserIBK   interbank.UserKey `json:"to_user_ibk" validate:"required"`
	Amount      decimal.Decimal   `json:"amount" validate:"required"`
}

func TransferRoute(c *fiber.Ctx) error {
	var body transferBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	if !utils.IsLocalUserIBK(body.ToUserIBK) {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"code": "wrong_bank_code",
		})
	}

	receiverUser, exists := storage.Users.FindUserById(int(body.ToUserIBK.UserId))
	if !exists {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"code": "receiver_not_found",
		})
	}

	transaction := storage.Transactions.CreateTransaction(
		body.FromUserIBK, body.ToUserIBK,
		body.Amount, models.TransactionTypeTransfer,
	)

	storage.Users.AddToUserBalance(receiverUser.Id, body.Amount)

	transaction = storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusCreated).JSON(&fiber.Map{
		"code": "transfer_success",
	})
}
