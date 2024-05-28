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
	mu   sync.RWMutex
	data map[int]models.User
}

var Users = &usersStorage{
	mu:   sync.RWMutex{},
	data: make(map[int]models.User),
}

func (us *usersStorage) CreateUser(name, email string) models.User {
	us.mu.Lock()
	user := models.User{
		Id:        len(us.data) + 1,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}

	user.InterBankKey = interbank.UserKey{
		BankId: config.Env.BankId,
		UserId: interbank.NewUserId(uint32(user.Id)),
	}

	us.data[user.Id] = user
	us.mu.Unlock()

	return user
}

func (us *usersStorage) FindUserById(id int) (models.User, bool) {
	us.mu.RLock()
	user, ok := us.data[id]
	us.mu.RUnlock()
	return user, ok
}

func (us *usersStorage) FindUserByEmail(email string) (models.User, bool) {
	us.mu.RLock()
	for _, user := range us.data {
		if user.Email == email {
			return user, true
		}
	}
	us.mu.RUnlock()
	return models.User{}, false
}

func (us *usersStorage) AddToUserBalance(userId int, amount decimal.Decimal) (models.User, bool) {
	us.mu.Lock()
	user, ok := us.data[userId]
	if !ok {
		us.mu.Unlock()
		return models.User{}, ok
	}

	user.Balance = user.Balance.Add(amount)
	us.data[userId] = user
	us.mu.Unlock()

	return user, ok
}

func (us *usersStorage) SubFromUserBalance(userId int, amount decimal.Decimal) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, ok := us.data[userId]
	if !ok {
		return errors.New("user not found")
	}

	if user.Balance.LessThan(amount) {
		return errors.New("insufficient funds")
	}

	user.Balance = user.Balance.Sub(amount)
	us.data[userId] = user

	return nil
}

func (us *usersStorage) TransferBalance(from, to int, amount decimal.Decimal) error {
	us.mu.Lock()
	fromUser, ok := us.data[from]
	if !ok {
		us.mu.Unlock()
		return errors.New("sender not found")
	}

	toUser, ok := us.data[to]
	if !ok {
		us.mu.Unlock()
		return errors.New("receiver not found")
	}

	if fromUser.Balance.LessThan(amount) {
		us.mu.Unlock()
		return errors.New("insufficient funds")
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	us.data[from] = fromUser
	us.data[to] = toUser
	us.mu.Unlock()

	return nil
}
