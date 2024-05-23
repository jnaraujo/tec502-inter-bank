package storage

import (
	"slices"
	"sync"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/shopspring/decimal"
)

type transactionsStorage struct {
	mu   sync.RWMutex
	data map[int]models.Transaction
}

var Transactions = &transactionsStorage{
	mu:   sync.RWMutex{},
	data: make(map[int]models.Transaction),
}

func (ts *transactionsStorage) CreateTransaction(from, to interbank.UserKey, amount decimal.Decimal, transactionType models.TransactionType) models.Transaction {
	ts.mu.Lock()

	transaction := models.Transaction{
		Id:        len(ts.data) + 1,
		From:      from,
		To:        &to,
		Amount:    amount,
		Type:      transactionType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ts.data[transaction.Id] = transaction

	ts.mu.Unlock()

	return transaction
}

func (ts *transactionsStorage) CreateDepositTransaction(from interbank.UserKey, amount decimal.Decimal, transactionType models.TransactionType) models.Transaction {
	ts.mu.Lock()

	transaction := models.Transaction{
		Id:        len(ts.data) + 1,
		From:      from,
		Amount:    amount,
		Type:      transactionType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ts.data[transaction.Id] = transaction

	ts.mu.Unlock()

	return transaction
}

func (ts *transactionsStorage) FindUserTransactionsById(userId int) []models.Transaction {
	transactions := []models.Transaction{}
	user, _ := Users.FindUserById(userId)

	ts.mu.RLock()
	for _, t := range ts.data {
		if t.From == user.InterBankKey || t.To == &user.InterBankKey {
			transactions = append(transactions, t)
		}
	}
	ts.mu.RUnlock()

	slices.SortStableFunc(transactions, func(a models.Transaction, b models.Transaction) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	return transactions
}
