package storage

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/shopspring/decimal"
)

type accountsStorage struct {
	mu   sync.RWMutex
	data map[int]models.Account
}

var Accounts = &accountsStorage{
	mu:   sync.RWMutex{},
	data: make(map[int]models.Account),
}

var accCounter atomic.Int64

func (as *accountsStorage) CreateAccount(name, document string, accType models.AccountType) models.Account {
	as.mu.Lock()
	defer as.mu.Unlock()

	id := accCounter.Add(1)
	user := models.Account{
		Id:        int(id),
		Name:      name,
		Documents: []string{document},
		Type:      accType,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}

	user.InterBankKey = interbank.IBK{
		BankId: config.Env.BankId,
		UserId: interbank.NewUserId(uint32(user.Id)),
	}

	as.data[user.Id] = user

	return user
}

func (as *accountsStorage) CreateJointAccount(name string, document []string) models.Account {
	as.mu.Lock()
	defer as.mu.Unlock()

	id := accCounter.Add(1)
	user := models.Account{
		Id:        int(id),
		Name:      name,
		Documents: document,
		Type:      models.AccountTypeJoint,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}

	user.InterBankKey = interbank.IBK{
		BankId: config.Env.BankId,
		UserId: interbank.NewUserId(uint32(user.Id)),
	}

	as.data[user.Id] = user

	return user
}

func (as *accountsStorage) FindAccountById(id int) (models.Account, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	user, ok := as.data[id]
	return user, ok
}

func (as *accountsStorage) Delete(id int) {
	as.mu.Lock()
	defer as.mu.Unlock()

	delete(as.data, id)
}

func (as *accountsStorage) FindUserByDocument(document string) (models.Account, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, user := range as.data {
		if slices.Contains(user.Documents, document) {
			return user, true
		}
	}
	return models.Account{}, false
}

// Busca todas as contas não conjuntas associadas a um documento
func (as *accountsStorage) FindAccountByDocument(document string) *models.Account {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, acc := range as.data {
		if acc.Type != models.AccountTypeJoint && acc.Documents[0] == document {
			return &acc
		}
	}
	return nil
}

// Busca todas as contas associadas a um documento (incluindo contas conjuntas)
func (as *accountsStorage) FindAllAccountsByDocument(document string) []models.Account {
	as.mu.RLock()
	defer as.mu.RUnlock()

	accounts := []models.Account{}
	for _, acc := range as.data {
		if slices.Contains(acc.Documents, document) {
			accounts = append(accounts, acc)
		}
	}
	return accounts
}

func (as *accountsStorage) FindIndividualAccountByDocument(document string) (models.Account, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, acc := range as.data {
		if acc.Type == models.AccountTypeIndividual && acc.Documents[0] == document {
			return acc, true
		}
	}
	return models.Account{}, false
}

func (as *accountsStorage) FindAccountByIBK(ibk interbank.IBK) *models.Account {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, user := range as.data {
		if user.InterBankKey == ibk {
			return &user
		}
	}
	return nil
}

func (as *accountsStorage) AddToAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	acc.Balance = acc.Balance.Add(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) SubFromAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	if acc.Balance.LessThan(amount) {
		return false
	}

	acc.Balance = acc.Balance.Sub(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) AddToPendingAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	acc.PendingBalance = acc.PendingBalance.Add(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) SubFromPendingAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	acc.PendingBalance = acc.PendingBalance.Sub(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) AddToBlockedAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	acc.BlockedBalance = acc.BlockedBalance.Add(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) SubFromBlockedAccountBalance(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, ok := as.data[accId]
	if !ok {
		return false
	}

	acc.BlockedBalance = acc.BlockedBalance.Sub(amount)
	as.data[accId] = acc

	return true
}

func (as *accountsStorage) CanSubFromAccount(accId int, amount decimal.Decimal) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	acc, exists := as.data[accId]
	finalBalance := acc.Balance.Sub(acc.BlockedBalance).Add(acc.PendingBalance)

	return exists && finalBalance.GreaterThanOrEqual(amount)
}

func (as *accountsStorage) TransferBalance(from, to int, amount decimal.Decimal) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	fromUser, ok := as.data[from]
	if !ok {
		return errors.New("sender not found")
	}

	toUser, ok := as.data[to]
	if !ok {
		return errors.New("receiver not found")
	}

	if fromUser.Balance.LessThan(amount) {
		return errors.New("insufficient funds")
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	as.data[from] = fromUser
	as.data[to] = toUser

	return nil
}
