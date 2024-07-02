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
	ParentId   *TransactionId    `json:"parent_id"` // ID da transação de origem, no caso de ser uma transação final
	Owner      interbank.IBK     `json:"owner"`
	Type       TransactionType   `json:"type"`
	Operations []Operation       `json:"operations"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Status     TransactionStatus `json:"status"`
}

func NewPackageTransaction(owner interbank.IBK, operations []Operation) *Transaction {
	return &Transaction{
		Id:         uuid.New(),
		Owner:      owner,
		Operations: operations,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     TransactionStatusPending,
		Type:       TransactionTypePackage,
	}
}

func NewFinalTransaction(parentID *TransactionId, owner interbank.IBK, operations []Operation) *Transaction {
	return &Transaction{
		Id:         uuid.New(),
		ParentId:   parentID,
		Owner:      owner,
		Operations: operations,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     TransactionStatusPending,
		Type:       TransactionTypeFinal,
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
