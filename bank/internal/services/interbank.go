package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/constants"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

var interbankClient = &http.Client{
	Timeout: constants.MaxInterBankRequestTimeout,
}

func FindAllUserAccountsInterBank(document string) []models.Account {
	accounts := []models.Account{}

	banks := storage.Ring.List()
	for _, bank := range banks {
		resp, err := interbankClient.Get(fmt.Sprintf("http://%s/interbank/account/%s", bank.Addr, document))
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
		if bank.Id == ibk.BankId {
			bankURL = bank.Addr
			break
		}
	}
	if bankURL == "" {
		slog.Error("Bank not found")
		return nil
	}

	resp, err := interbankClient.Get(fmt.Sprintf("http://%s/interbank/account/ibk/%s", bankURL, ibk))
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

type txProcess struct {
	Tx   *models.Transaction
	Step Step
}

func ProcessTransaction(tr models.Transaction) error {
	externalTransactions := []txProcess{}

	isSuccess := true
	for _, op := range tr.Operations {
		txDebit := Prepare(tr.Id, op, StepDebit)
		if txDebit == nil {
			isSuccess = false
			break
		}
		externalTransactions = append(externalTransactions, txProcess{Tx: txDebit, Step: StepDebit})

		txCredit := Prepare(tr.Id, op, StepCredit)
		if txCredit == nil {
			isSuccess = false
			break
		}
		externalTransactions = append(externalTransactions, txProcess{Tx: txCredit, Step: StepCredit})
	}

	if !isSuccess {
		slog.Warn("Erro em alguma das transações")
		// se ocorreu algum erro, as transações ja feitas devem sofrer rollback
		for _, tx := range externalTransactions {
			Rollback(tx.Tx.Id, tx.Tx.Operations[0], tx.Step)
		}

		for _, op := range tr.Operations {
			storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
		}
		storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed)

		return errors.New("transaction failed")
	}

	// se tudo estiver correto, as transações são confirmadas
	for _, tx := range externalTransactions {
		ok := Commit(tx.Tx.Id, tx.Tx.Operations[0], tx.Step)
		if !ok {
			isSuccess = false
			break
		}
	}

	// se ocorreu algum erro, as transações ja feitas devem sofrer rollback
	// mesmo as que tiveram sucesso
	if !isSuccess {
		slog.Warn("Erro em alguma das transações na confirmação")
		for _, tx := range externalTransactions {
			Rollback(tx.Tx.Id, tx.Tx.Operations[0], tx.Step)
		}

		for _, op := range tr.Operations {
			storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusFailed)
		}
		storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusFailed)

		return errors.New("transaction failed")
	}

	for _, op := range tr.Operations {
		storage.Transactions.UpdateOperationStatus(tr, op, models.OperationStatusSuccess)
	}
	storage.Transactions.UpdateTransactionStatus(tr, models.TransactionStatusSuccess)

	return nil
}

var txClient = &http.Client{
	Timeout: constants.OperationTimeout,
}

type Step string

const (
	StepCredit Step = "credit"
	StepDebit  Step = "debit"
)

func Prepare(parentId models.TransactionId, op models.Operation, step Step) *models.Transaction {
	reqBody, _ := json.Marshal(map[string]any{
		"parent_id": parentId,
		"operation": map[string]string{
			"from":   op.From.String(),
			"to":     op.To.String(),
			"amount": op.Amount.String(),
		},
		"step": step,
	})

	bank := storage.Ring.Find(op.From.BankId) // StepDebit
	if step == StepCredit {
		bank = storage.Ring.Find(op.To.BankId)
	}
	if bank == nil {
		slog.Error("Bank not found")
		return nil
	}

	resp, err := txClient.Post(fmt.Sprintf("http://%s/interbank/prepare", bank.Addr), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error preparing transaction", "status", resp.StatusCode, "body", string(body))
		return nil
	}

	var response models.Transaction
	json.Unmarshal(body, &response)

	return &response
}

func Rollback(txId uuid.UUID, op models.Operation, step Step) bool {
	reqBody, _ := json.Marshal(map[string]string{
		"tx_id": txId.String(),
		"step":  string(step),
	})

	bank := storage.Ring.Find(op.From.BankId) // StepDebit
	if step == StepCredit {
		bank = storage.Ring.Find(op.To.BankId)
	}
	if bank == nil {
		slog.Error("Bank not found")
		return false
	}

	resp, err := txClient.Post(fmt.Sprintf("http://%s/interbank/rollback", bank.Addr), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]string
	json.Unmarshal(body, &response)

	return resp.StatusCode == http.StatusOK
}

func Commit(txId uuid.UUID, op models.Operation, step Step) bool {
	reqBody, _ := json.Marshal(map[string]string{
		"tx_id": txId.String(),
		"step":  string(step),
	})

	bank := storage.Ring.Find(op.From.BankId) // StepDebit
	if step == StepCredit {
		bank = storage.Ring.Find(op.To.BankId)
	}
	if bank == nil {
		slog.Error("Bank not found")
		return false
	}

	resp, err := txClient.Post(fmt.Sprintf("http://%s/interbank/commit", bank.Addr), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]string
	json.Unmarshal(body, &response)

	return resp.StatusCode == http.StatusOK
}
