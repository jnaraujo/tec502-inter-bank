package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

type OperationStatus string

const (
	OperationStatusPending OperationStatus = "pending"
	OperationStatusSuccess OperationStatus = "success"
	OperationStatusFailed  OperationStatus = "failed"
)

type OperationType string

const (
	OperationTypeDeposit  OperationType = "deposit"  // Depósito de dinheiro
	OperationTypeTransfer OperationType = "transfer" // Transferência de dinheiro
)

type Operation struct {
	Id        uuid.UUID       `json:"id"`
	From      interbank.IBK   `json:"from"`
	To        interbank.IBK   `json:"to"`
	Type      OperationType   `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Status    OperationStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewOperation(from, to interbank.IBK, opType OperationType, amount decimal.Decimal) *Operation {
	return &Operation{
		Id:        uuid.New(),
		From:      from,
		To:        to,
		Type:      opType,
		Amount:    amount,
		Status:    OperationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewDepositOperation(user interbank.IBK, amount decimal.Decimal) *Operation {
	return NewOperation(user, user, OperationTypeDeposit, amount)
}
