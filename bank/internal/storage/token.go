package storage

import (
	"sync"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token"
)

type tokenStorage struct {
	mu    sync.RWMutex
	Token token.Token
}

var Token = &tokenStorage{}

func (ts *tokenStorage) Set(token token.Token) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Token = token
}

func (ts *tokenStorage) Get() token.Token {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.Token
}

func (ts *tokenStorage) HasValidToken() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.Token.IsOwnerInternal() && !ts.Token.HasExpired()
}

func (ts *tokenStorage) HasInvalidToken() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.Token.IsOwnerInternal() && ts.Token.HasExpired()
}
