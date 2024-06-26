package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
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

		var acc []models.Account
		json.Unmarshal(body, &acc)

		accounts = append(accounts, acc...)
	}

	return accounts
}

func FindAccountInterBank(ibk interbank.IBK) *models.Account {
	bankURL := ""
	for _, bank := range storage.Ring.List() {
		fmt.Println(bank.Id, ibk.BankId)
		if bank.Id == ibk.BankId {
			bankURL = bank.Addr
			break
		}
	}

	if bankURL == "" {
		fmt.Println("Bank not found")
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/interbank/account/ibk/%s", bankURL, ibk))
	if err != nil {
		fmt.Println("Error getting account")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Account not found")
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body")
		return nil
	}

	var acc models.Account
	json.Unmarshal(body, &acc)

	return &acc
}
