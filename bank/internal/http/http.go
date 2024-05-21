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
	router := app.Group("/bank")

	router.Get("/", bank.GetRootRoute)
}
