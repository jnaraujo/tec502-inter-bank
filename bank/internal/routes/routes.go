package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes/bank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes/interbank"
)

func RegisterInterBankRoutes(app *fiber.App) {
	router := app.Group("/interbank")

	router.Post("/add-credit", interbank.AddCreditRoute)
	router.Post("/sub-credit", interbank.SubCreditRoute)
	router.Get("/account/:document", interbank.FindAccountRoute)

	router.Post("/token", interbank.SetToken)
	router.Get("/token", interbank.GetToken)
	router.Get("/token/ok", interbank.CanReceiveToken)
}

func RegisterBankRoutes(app *fiber.App) {
	router := app.Group("/api")

	router.Get("/", bank.GetRootRoute)

	router.Post("/accounts", bank.CreateAccountRoute)
	router.Get("/accounts/:id", bank.FindAccountRoute)
	router.Delete("/accounts/:id", bank.DeleteAccountRoute)
	router.Get("/accounts/:id/transactions", bank.ListAccountTransactionsRoute)
	router.Get("/accounts/:id/all", bank.FindAllAccountsRoute)

	router.Post("/payments/deposit", bank.DepositRoute)
	router.Post("/payments/pay", bank.PayRoute)
}
