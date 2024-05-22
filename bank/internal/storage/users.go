package storage

import (
	"sync"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
	"github.com/shopspring/decimal"
)

type UsersStorage struct {
	sync.RWMutex
	data map[int]models.User
}

var users = &UsersStorage{
	RWMutex: sync.RWMutex{},
	data:    make(map[int]models.User),
}

func CreateUser(name string) models.User {
	users.Lock()
	user := &models.User{
		Id:        len(users.data) + 1,
		Name:      name,
		CreatedAt: time.Now(),
		Balance:   decimal.NewFromInt(0),
	}
	users.data[user.Id] = *user
	users.Unlock()

	return *user
}

func FindUserById(id int) (models.User, bool) {
	users.RLock()
	user, ok := users.data[id]
	users.RUnlock()
	return user, ok
}
