package storage

import (
	"sync"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token"
)

type tokenStorage struct {
	mu    sync.RWMutex
	token token.Token
}

var Token = &tokenStorage{}

func (ts *tokenStorage) Set(token token.Token) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.token = token
}

func (ts *tokenStorage) Get() token.Token {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.token
}

func (ts *tokenStorage) HasToken() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.token.IsOwnerInternal() && ts.token.IsValid()
}
