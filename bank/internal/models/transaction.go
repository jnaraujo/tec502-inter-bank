package models

import (
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"  // Depósito de dinheiro
	TransactionTypeWithdraw TransactionType = "withdraw" // Saque de dinheiro
	TransactionTypeTransfer TransactionType = "transfer" // Transferência de dinheiro
	TransactionTypePayment  TransactionType = "payment"  // Pagamento de boleto, fatura, etc
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	Id        int                `json:"id"`
	From      interbank.UserKey  `json:"from"`
	To        *interbank.UserKey `json:"to,omitempty"`
	Amount    decimal.Decimal    `json:"amount"`
	Type      TransactionType    `json:"type"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Status    TransactionStatus  `json:"status"`
}
