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

type usersStorage struct {
	sync.RWMutex
	data map[int]models.User
}

var Users = &usersStorage{
	RWMutex: sync.RWMutex{},
	data:    make(map[int]models.User),
}

func (us *usersStorage) CreateUser(name string) models.User {
	us.Lock()
	user := models.User{
		Id:        len(us.data) + 1,
		Name:      name,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}

	user.InterBankKey = interbank.UserKey{
		BankId: config.Env.BankId,
		UserId: interbank.NewUserId(uint32(user.Id)),
	}

	us.data[user.Id] = user
	us.Unlock()

	return user
}

func (us *usersStorage) FindUserById(id int) (models.User, bool) {
	us.RLock()
	user, ok := us.data[id]
	us.RUnlock()
	return user, ok
}

func (us *usersStorage) AddToUserBalance(userId int, amount decimal.Decimal) (models.User, bool) {
	us.Lock()
	user, ok := us.data[userId]
	if !ok {
		us.Unlock()
		return models.User{}, ok
	}

	user.Balance = user.Balance.Add(amount)
	us.data[userId] = user
	us.Unlock()

	return user, ok
}

func (us *usersStorage) TransferBalance(from, to int, amount decimal.Decimal) error {
	us.Lock()
	fromUser, ok := us.data[from]
	if !ok {
		us.Unlock()
		return errors.New("sender not found")
	}

	toUser, ok := us.data[to]
	if !ok {
		us.Unlock()
		return errors.New("receiver not found")
	}

	if fromUser.Balance.LessThan(amount) {
		us.Unlock()
		return errors.New("insufficient funds")
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	us.data[from] = fromUser
	us.data[to] = toUser
	us.Unlock()

	return nil
}
