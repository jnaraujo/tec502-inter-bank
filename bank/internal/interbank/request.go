package interbank

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/shopspring/decimal"
)

type Code string

const (
	TransferSuccess  Code = "transfer_success"
	WrongBankCode    Code = "wrong_bank_code"
	ReceiverNotFound Code = "receiver_not_found"
)

type PaymentRequestResponse struct {
	Code Code
}

func SendPaymentRequest(from, to IBK, amount decimal.Decimal) (*PaymentRequestResponse, error) {
	body := map[string]any{
		"from_user_ibk": from,
		"to_user_ibk":   to,
		"amount":        amount,
	}
	m, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:3001/interbank/transfer", "application/json", bytes.NewReader(m))
	if err != nil {
		return nil, errors.New("could not send payment request")
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read response")
	}

	jsonResp := PaymentRequestResponse{}
	err = json.Unmarshal(resBody, &jsonResp)
	if err != nil {
		return nil, errors.New("invalid response")
	}

	return &jsonResp, nil
}
