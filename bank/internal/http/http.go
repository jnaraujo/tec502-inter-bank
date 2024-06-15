package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/routes"
)

func NewServer(port int) error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.RegisterBankRoutes(app)
	routes.RegisterInterBankRoutes(app)

	err := app.Listen(fmt.Sprintf("%s:%d", "0.0.0.0", port))
	if err != nil {
		return err
	}
	return nil
}
