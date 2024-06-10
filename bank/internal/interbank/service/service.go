package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

func BeginTransaction(bankId int) bool {
	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := http.Post(fmt.Sprintf("%s/interbank/lock", ip), "application/json", nil)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func AddCredit(bankId int, to interbank.IBK, amount decimal.Decimal) error {
	ip := fmt.Sprintf("http://localhost:300%d", bankId)

	reqBody, _ := json.Marshal(map[string]string{
		"to":     to.String(),
		"amount": amount.String(),
	})
	resp, err := http.Post(fmt.Sprintf("%s/interbank/add-credit", ip), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
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
	ip := fmt.Sprintf("http://localhost:300%d", bankId)

	reqBody, _ := json.Marshal(map[string]string{
		"from":   from.String(),
		"amount": amount.String(),
	})
	resp, err := http.Post(fmt.Sprintf("%s/interbank/sub-credit", ip), "application/json", bytes.NewBuffer(reqBody))
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

func CommitTransaction(bankId int) bool {
	ip := fmt.Sprintf("http://localhost:300%d", bankId)
	resp, err := http.Post(fmt.Sprintf("%s/interbank/unlock", ip), "application/json", nil)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
