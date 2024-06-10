package storage

import (
	"errors"
	"sync"
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

func (as *accountsStorage) CreateAccount(name, document string) models.Account {
	user := models.Account{
		Id:        len(as.data) + 1,
		Name:      name,
		Document:  document,
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

func (as *accountsStorage) FindUserById(id int) (models.Account, bool) {
	user, ok := as.data[id]
	return user, ok
}

func (as *accountsStorage) Delete(id int) {
	delete(as.data, id)
}

func (as *accountsStorage) FindUserByDocument(document string) (models.Account, bool) {
	for _, user := range as.data {
		if user.Document == document {
			return user, true
		}
	}
	return models.Account{}, false
}

func (as *accountsStorage) AddToUserBalance(userId int, amount decimal.Decimal) (models.Account, bool) {
	user, ok := as.data[userId]
	if !ok {
		return models.Account{}, ok
	}

	user.Balance = user.Balance.Add(amount)
	as.data[userId] = user

	return user, ok
}

func (as *accountsStorage) SubFromUserBalance(userId int, amount decimal.Decimal) error {
	user, ok := as.data[userId]
	if !ok {
		return errors.New("user not found")
	}

	if user.Balance.LessThan(amount) {
		return errors.New("insufficient funds")
	}

	user.Balance = user.Balance.Sub(amount)
	as.data[userId] = user

	return nil
}

func (as *accountsStorage) TransferBalance(from, to int, amount decimal.Decimal) error {
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

func (as *accountsStorage) Lock() {
	as.mu.Lock()
}

func (as *accountsStorage) Unlock() {
	as.mu.Unlock()
}

func (as *accountsStorage) RLock() {
	as.mu.RLock()
}

func (as *accountsStorage) RUnlock() {
	as.mu.RUnlock()
}

func (as *accountsStorage) IsLocked() bool {
	isLocked := as.mu.TryLock()
	if isLocked {
		as.mu.Unlock()
	}
	return !isLocked
}
