package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes/bank"
)

func NewServer() error {
	app := fiber.New()

	registerBankRoutes(app)

	err := app.Listen(fmt.Sprintf("%s:%d", "0.0.0.0", config.Env.ServerPort))
	if err != nil {
		return err
	}
	return nil
}

func registerBankRoutes(app *fiber.App) {
	router := app.Group("/api")

	router.Get("/", bank.GetRootRoute)

	router.Post("/accounts", bank.CreateAccountRoute)
	router.Get("/accounts/:id", bank.FindAccountRoute)
	router.Get("/accounts/:id/transactions", bank.ListAccountTransactionsRoute)

	router.Post("/payments/deposit", bank.DepositRoute)
}
