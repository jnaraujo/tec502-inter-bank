package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/shopspring/decimal"
)

const (
	transactionTimeout = 5 * time.Second
	operationTimeout   = 1 * time.Second
)

func LockAccountsFromTransaction(tr models.Transaction) {
	banks := map[int]bool{}
	for _, op := range tr.Operations {
		banks[int(op.From.BankId)] = true
		banks[int(op.To.BankId)] = true
	}

	for bankId := range banks {
		err := BeginTransaction(bankId)
		if err != nil {
			// Provavelmente um deadlock
			if os.IsTimeout(err) {
				fmt.Println("Provavelmente um deadlock! Tentando novamente...")
				// termina todas as transações
				// e espera um tempo para tentar novamente
				UnlockAccountsFromTransaction(tr)
				// espera um tempo aleatório entre 0 e 300ms
				time.Sleep(rand.N(300 * time.Millisecond))
				LockAccountsFromTransaction(tr)
				return
			}

			// Se não for um deadlock, é um erro
			RollbackOperations(tr)
			return
		}
	}
}

func UnlockAccountsFromTransaction(tr models.Transaction) {
	banks := map[int]bool{}
	for _, op := range tr.Operations {
		banks[int(op.From.BankId)] = true
		banks[int(op.To.BankId)] = true
	}

	for bankId := range banks {
		EndTransaction(bankId)
	}
}

func RollbackOperations(tr models.Transaction) {
	for _, op := range tr.Operations {
		// so precisa reverter as que tiveram sucesso
		if op.Status == models.OperationStatusSuccess {
			SubCredit(int(op.To.BankId), op.To, op.Amount)
			AddCredit(int(op.From.BankId), op.From, op.Amount)
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
	}

	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed)
}

func BeginTransaction(bankId int) error {
	client := http.Client{
		Timeout: transactionTimeout,
	}

	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := client.Post(fmt.Sprintf("%s/interbank/lock", ip), "application/json", nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("error starting transaction")
	}

	return nil
}

func EndTransaction(bankId int) error {
	client := http.Client{
		Timeout: transactionTimeout,
	}

	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := client.Post(fmt.Sprintf("%s/interbank/unlock", ip), "application/json", nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("error ending transaction")
	}

	return nil
}

func AddCredit(bankId int, to interbank.IBK, amount decimal.Decimal) error {
	client := http.Client{
		Timeout: operationTimeout,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"to":     to.String(),
		"amount": amount.String(),
	})

	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := client.Post(fmt.Sprintf("%s/interbank/add-credit", ip), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Println("timeout!")
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

func SubCredit(bankId int, from interbank.IBK, amount decimal.Decimal) error {
	client := http.Client{
		Timeout: operationTimeout,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"from":   from.String(),
		"amount": amount.String(),
	})
	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := client.Post(fmt.Sprintf("%s/interbank/sub-credit", ip), "application/json", bytes.NewBuffer(reqBody))
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
