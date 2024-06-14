package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes/bank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes/interbank"
)

func NewServer(port int) error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	registerBankRoutes(app)
	registerInterBankRoutes(app)

	err := app.Listen(fmt.Sprintf("%s:%d", "0.0.0.0", port))
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
	router.Delete("/accounts/:id", bank.DeleteAccountRoute)
	router.Get("/accounts/:id/transactions", bank.ListAccountTransactionsRoute)

	router.Post("/payments/deposit", bank.DepositRoute)
	router.Post("/payments/pay", bank.PayRoute)
}

func registerInterBankRoutes(app *fiber.App) {
	router := app.Group("/interbank")

	router.Post("/add-credit", interbank.AddCreditRoute)
	router.Post("/sub-credit", interbank.SubCreditRoute)
}
