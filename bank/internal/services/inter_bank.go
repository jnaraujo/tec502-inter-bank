package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/constants"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/shopspring/decimal"
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
		slog.Info(fmt.Sprintf("Bank: %d", bank.Id))
		if bank.Id == ibk.BankId {
			bankURL = bank.Addr
			break
		}
	}

	if bankURL == "" {
		slog.Error("Bank not found")
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/interbank/account/ibk/%s", bankURL, ibk))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("Account not found")
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	var acc models.Account
	json.Unmarshal(body, &acc)

	return &acc
}

func SubCreditFromAccount(from interbank.IBK, amount decimal.Decimal) error {
	client := http.Client{
		Timeout: constants.OperationTimeout,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"from":   from.String(),
		"amount": amount.String(),
	})

	resp, err := client.Post(fmt.Sprintf("http://localhost:300%d/interbank/sub-credit", from.BankId), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.New("bank is offline")
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]string
	json.Unmarshal(body, &response)

	if resp.StatusCode != http.StatusOK {
		return errors.New(response["message"])
	}

	return nil
}

func AddCreditToAccount(to interbank.IBK, amount decimal.Decimal) error {
	client := http.Client{
		Timeout: constants.OperationTimeout,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"to":     to.String(),
		"amount": amount.String(),
	})

	resp, err := client.Post(fmt.Sprintf("http://localhost:300%d/interbank/add-credit", to.BankId), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		if os.IsTimeout(err) {
			slog.Error("Timeout")
		}

		return errors.New("bank is offline")
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return errors.New("bank is offline")
	}

	var response map[string]string
	json.Unmarshal(body, &response)

	if resp.StatusCode != http.StatusOK {
		return errors.New(response["message"])
	}

	return nil
}

func ProcessTransaction(tr models.Transaction) error {
	for _, op := range tr.Operations {
		err := SubCreditFromAccount(op.From, op.Amount)
		if err != nil {
			RollbackOperations(tr)
			return err
		}

		err = AddCreditToAccount(op.To, op.Amount)
		if err != nil {
			// como falou na segunda parte, reverte a primeira parte
			AddCreditToAccount(op.From, op.Amount)
			RollbackOperations(tr)
			return err
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusSuccess)
	}

	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusSuccess)

	return nil
}

func RollbackOperations(tr models.Transaction) {
	for _, op := range tr.Operations {
		// so precisa reverter as que tiveram sucesso
		if op.Status == models.OperationStatusSuccess {
			SubCreditFromAccount(op.To, op.Amount)
			AddCreditToAccount(op.From, op.Amount)
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
	}

	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed)
}
