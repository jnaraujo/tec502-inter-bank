package interbank

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/transaction_processor"
)

func ReceiveToken(c *fiber.Ctx) error {
	fmt.Println("Token received")
	go func() {
		transaction_processor.Process()
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "token received",
	})
}
