package storage

import (
	"sync"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/shopspring/decimal"
)

type transactionsStorage struct {
	sync.RWMutex
	data map[int]models.Transaction
}

var transactions = &transactionsStorage{
	RWMutex: sync.RWMutex{},
	data:    make(map[int]models.Transaction),
}

func CreateTransaction(from, to interbank.UserKey, amount decimal.Decimal, transactionType models.TransactionType) models.Transaction {
	transactions.Lock()

	t := models.Transaction{
		Id:        len(transactions.data) + 1,
		From:      from,
		To:        &to,
		Amount:    amount,
		Type:      transactionType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transactions.data[t.Id] = t

	transactions.Unlock()

	return t
}

func CreateDepositTransaction(from interbank.UserKey, amount decimal.Decimal, transactionType models.TransactionType) models.Transaction {
	transactions.Lock()

	t := models.Transaction{
		Id:        len(transactions.data) + 1,
		From:      from,
		Amount:    amount,
		Type:      transactionType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	transactions.data[t.Id] = t

	transactions.Unlock()

	return t
}

func FindUserTransactionsById(userId int) []models.Transaction {
	ts := []models.Transaction{}
	user, _ := FindUserById(userId)

	transactions.RLock()
	for _, t := range transactions.data {
		if t.From == user.InterBankKey || t.To == &user.InterBankKey {
			ts = append(ts, t)
		}
	}
	transactions.RUnlock()

	return ts
}
