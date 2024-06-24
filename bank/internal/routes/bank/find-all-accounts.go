package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAllAccountsRoute(c *fiber.Ctx) error {
	accountId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid account id",
		})
	}

	account, ok := storage.Accounts.FindUserById(accountId)
	if !ok {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "User does not exists",
		})
	}

	accounts := []models.Account{}

	banks := storage.Ring.List()
	for _, bank := range banks {
		resp, err := http.Get(fmt.Sprintf("http://%s/interbank/account/%s", bank.Addr, account.Document))
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			continue
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var acc models.Account
		json.Unmarshal(body, &acc)

		accounts = append(accounts, acc)
	}

	return c.Status(http.StatusOK).JSON(&accounts)
}
