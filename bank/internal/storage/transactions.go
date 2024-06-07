package storage

import (
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/shopspring/decimal"
)

type transactionsStorage struct {
	mu   sync.RWMutex
	data map[uuid.UUID]models.Transaction
}

var Transactions = &transactionsStorage{
	mu:   sync.RWMutex{},
	data: make(map[uuid.UUID]models.Transaction),
}

func (ts *transactionsStorage) CreateTransaction(author interbank.IBK, operations []models.Operation) models.Transaction {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	transaction := models.Transaction{
		Id:         uuid.New(),
		Author:     author,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     models.TransactionStatusPending,
		Operations: operations,
	}
	ts.data[transaction.Id] = transaction

	return transaction
}

func (ts *transactionsStorage) CreateDepositTransaction(author interbank.IBK, amount decimal.Decimal) models.Transaction {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	operation := models.Operation{
		From:      author,
		To:        author,
		Type:      models.OperationTypeDeposit,
		Amount:    amount,
		Status:    models.OperationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transaction := models.Transaction{
		Id:         uuid.New(),
		Author:     author,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     models.TransactionStatusPending,
		Operations: []models.Operation{operation},
	}
	ts.data[transaction.Id] = transaction

	return transaction
}

func (ts *transactionsStorage) FindUserTransactionsById(userId int) []models.Transaction {
	transactions := []models.Transaction{}
	user, _ := Users.FindUserById(userId)

	ts.mu.RLock()
	for _, t := range ts.data {
		if t.Author == user.InterBankKey {
			transactions = append(transactions, t)
		}
	}
	ts.mu.RUnlock()

	slices.SortStableFunc(transactions, func(a models.Transaction, b models.Transaction) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	return transactions
}

func (ts *transactionsStorage) UpdateOperationStatus(t models.Transaction, op models.Operation, status models.OperationStatus) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t.UpdatedAt = time.Now()
	op.UpdatedAt = time.Now()

	op.Status = status

	found := false
	for idx, o := range t.Operations {
		if o.Id == op.Id {
			t.Operations[idx] = op
			found = true
		}
	}

	if !found {
		return false
	}

	ts.data[t.Id] = t

	return true
}

func (ts *transactionsStorage) UpdateTransactionStatus(t models.Transaction, s models.TransactionStatus) models.Transaction {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t.UpdatedAt = time.Now()
	t.Status = s
	ts.data[t.Id] = t

	return t
}
