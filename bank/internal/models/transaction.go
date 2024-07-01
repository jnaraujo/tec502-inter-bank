package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type TransactionType string

const (
	TransactionTypePackage TransactionType = "package" // representa um pacote de operações
	TransactionTypeFinal   TransactionType = "final"   // representa uma operação final
)

type TransactionId = uuid.UUID

type Transaction struct {
	Id         TransactionId     `json:"id"`
	Owner      interbank.IBK     `json:"owner"`
	Type       TransactionType   `json:"type"`
	Operations []Operation       `json:"operations"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Status     TransactionStatus `json:"status"`
}

func NewTransaction(owner interbank.IBK, operations []Operation, trType TransactionType) *Transaction {
	return &Transaction{
		Id:         uuid.New(),
		Owner:      owner,
		Operations: operations,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     TransactionStatusPending,
		Type:       trType,
	}
}

func (tr *Transaction) UpdateOperation(operation Operation) bool {
	for idx := range tr.Operations {
		if tr.Operations[idx].Id == operation.Id {
			tr.Operations[idx] = operation
			return true
		}
	}
	return false
}
