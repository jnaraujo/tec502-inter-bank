package ibt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func Process(tr models.Transaction) error {
	for _, op := range tr.Operations {
		err := SubCredit(int(op.From.BankId), op.From, op.Amount)
		if err != nil {
			RollbackOperations(tr)
			return err
		}

		err = AddCredit(int(op.To.BankId), op.To, op.Amount)
		if err != nil {
			// como falou na segunda parte, reverte a primeira parte
			AddCredit(int(op.From.BankId), op.From, op.Amount)
			RollbackOperations(tr)
			return err
		}

		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusSuccess)
	}

	return nil
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
