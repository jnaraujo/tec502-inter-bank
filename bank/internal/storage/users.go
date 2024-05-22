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

var users = &usersStorage{
	RWMutex: sync.RWMutex{},
	data:    make(map[int]models.User),
}

func CreateUser(name string) models.User {
	users.Lock()
	user := models.User{
		Id:        len(users.data) + 1,
		Name:      name,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}

	user.InterBankKey = interbank.UserKey{
		BankId: config.Env.BankId,
		UserId: interbank.NewUserId(uint32(user.Id)),
	}

	users.data[user.Id] = user
	users.Unlock()

	return user
}

func FindUserById(id int) (models.User, bool) {
	users.RLock()
	user, ok := users.data[id]
	users.RUnlock()
	return user, ok
}

func AddToUserBalance(userId int, amount decimal.Decimal) (models.User, bool) {
	users.Lock()
	user, ok := users.data[userId]
	if !ok {
		users.Unlock()
		return models.User{}, ok
	}

	user.Balance = user.Balance.Add(amount)
	users.data[userId] = user
	users.Unlock()

	return user, ok
}

func TransferBalance(from, to int, amount decimal.Decimal) error {
	users.Lock()
	fromUser, ok := users.data[from]
	if !ok {
		users.Unlock()
		return errors.New("sender not found")
	}

	toUser, ok := users.data[to]
	if !ok {
		users.Unlock()
		return errors.New("receiver not found")
	}

	if fromUser.Balance.LessThan(amount) {
		users.Unlock()
		return errors.New("insufficient funds")
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	users.data[from] = fromUser
	users.data[to] = toUser
	users.Unlock()

	return nil
}
