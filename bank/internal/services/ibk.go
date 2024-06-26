package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func FindAllUserAccountsInterBank(document string) []models.Account {
	accounts := []models.Account{}

	banks := storage.Ring.List()
	for _, bank := range banks {
		resp, err := http.Get(fmt.Sprintf("http://%s/interbank/account/%s", bank.Addr, document))
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

	return accounts
}
