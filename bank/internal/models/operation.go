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
	OperationTypeAdd     OperationType = "add"     // Adiciona dinheiro
	OperationTypeSub     OperationType = "sub"     // Subtrai dinheiro
	OperationTypeDeposit OperationType = "deposit" // Adiciona dinheiro
)

type Operation struct {
	Id        uuid.UUID
	User      interbank.IBK   `json:"user"`
	Type      OperationType   `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Status    OperationStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewOperation(user interbank.IBK, opType OperationType, amount decimal.Decimal) *Operation {
	return &Operation{
		Id:        uuid.New(),
		User:      user,
		Type:      opType,
		Amount:    amount,
		Status:    OperationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewAddOperation(user interbank.IBK, amount decimal.Decimal) *Operation {
	return NewOperation(user, OperationTypeAdd, amount)
}

func NewSubOperation(user interbank.IBK, amount decimal.Decimal) *Operation {
	return NewOperation(user, OperationTypeSub, amount)
}

func NewDepositOperation(user interbank.IBK, amount decimal.Decimal) *Operation {
	return NewOperation(user, OperationTypeDeposit, amount)
}

func NewTransferOperations(from, to interbank.IBK, amount decimal.Decimal) []Operation {
	operations := []Operation{}

	operations = append(operations, *NewSubOperation(from, amount))
	operations = append(operations, *NewAddOperation(to, amount))

	return operations
}
