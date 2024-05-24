package bank

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

type payRouteBodySchema struct {
	FromUserIBK interbank.UserKey `json:"from_user_ibk" validate:"required"`
	ToUserIBK   interbank.UserKey `json:"to_user_ibk" validate:"required"`
	Amount      decimal.Decimal   `json:"amount" validate:"required"`
}

func PayRoute(c *fiber.Ctx) error {
	var body payRouteBodySchema
	if errs := validate.ParseAndValidate(c.Body(), &body); len(errs) > 0 {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": errs})
	}

	if !utils.IsLocalUserIBK(body.FromUserIBK) {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": "Sender must be from this bank",
		})
	}

	if body.FromUserIBK == body.ToUserIBK {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Sender and receiver must be different",
		})
	}

	if body.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Amount must be greater than 0",
		})
	}

	// transação interna
	if utils.IsLocalUserIBK(body.ToUserIBK) {
		return handleInternalTransaction(
			c,
			body.FromUserIBK, body.ToUserIBK,
			body.Amount,
		)
	}

	fromUserId := int(body.FromUserIBK.UserId)
	err := storage.Users.SubFromUserBalance(fromUserId, body.Amount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	transaction := storage.Transactions.CreateTransaction(
		body.FromUserIBK, body.ToUserIBK,
		body.Amount, models.TransactionTypeTransfer,
	)

	resp, err := interbank.SendPaymentRequest(body.FromUserIBK, body.ToUserIBK, body.Amount)
	if err != nil {
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)
		storage.Users.AddToUserBalance(fromUserId, body.Amount)

		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if resp.Code != interbank.TransferSuccess {
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)
		storage.Users.AddToUserBalance(fromUserId, body.Amount)

		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": resp.Code,
		})
	}

	transaction = storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusNotImplemented).JSON(&transaction)
}

func handleInternalTransaction(c *fiber.Ctx, fromUserIBK, toUserIBK interbank.UserKey, amount decimal.Decimal) error {
	transaction := storage.Transactions.CreateTransaction(
		fromUserIBK, toUserIBK,
		amount, models.TransactionTypeTransfer,
	)

	fromUserId := int(fromUserIBK.UserId)
	toUserId := int(toUserIBK.UserId)

	err := storage.Users.TransferBalance(fromUserId, toUserId, amount)
	if err != nil {
		storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusFailed)
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	transaction = storage.Transactions.UpdateTransactionStatus(transaction, models.TransactionStatusSuccess)

	return c.Status(http.StatusInternalServerError).JSON(&transaction)
}
