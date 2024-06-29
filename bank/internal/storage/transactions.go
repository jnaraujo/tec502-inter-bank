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

func (ts *transactionsStorage) Save(tr models.Transaction) {
	ts.mu.Lock()
	ts.data[tr.Id] = tr
	ts.mu.Unlock()
}

func (ts *transactionsStorage) CreateDepositTransaction(author interbank.IBK, amount decimal.Decimal) models.Transaction {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	operation := *models.NewDepositOperation(
		author,
		amount,
	)

	transaction := *models.NewTransaction(
		author,
		[]models.Operation{operation},
	)
	ts.data[transaction.Id] = transaction

	return transaction
}

func (ts *transactionsStorage) FindUserTransactionsById(userId int) []models.Transaction {
	transactions := []models.Transaction{}

	user, _ := Accounts.FindAccountById(userId)

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

	found := t.UpdateOperation(op)
	return found
}

func (ts *transactionsStorage) UpdateTransactionStatus(t models.Transaction, s models.TransactionStatus) models.Transaction {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t.UpdatedAt = time.Now()
	t.Status = s
	ts.data[t.Id] = t

	return t
}

func (ts *transactionsStorage) FindTransactionById(id uuid.UUID) *models.Transaction {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tr := ts.data[id]
	return &tr
}
