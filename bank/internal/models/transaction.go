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

type Transaction struct {
	Id         uuid.UUID         `json:"id"`
	Author     interbank.IBK     `json:"author"`
	Operations []Operation       `json:"operations"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Status     TransactionStatus `json:"status"`
}

func NewTransaction(author interbank.IBK, operations []Operation) *Transaction {
	return &Transaction{
		Id:         uuid.New(),
		Author:     author,
		Operations: operations,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     TransactionStatusPending,
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
