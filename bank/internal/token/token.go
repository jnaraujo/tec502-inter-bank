package token

import (
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
)

const (
	MAX_TOKEN_AGE = 30 * time.Second // O tempo máximo que o token pode ser considerado válido. Um tempo maior que isso indica que o token está preso em algum lugar ou que esse banco está fora do ar.
)

type Token struct {
	Owner interbank.BankId `json:"owner"`
	Ts    time.Time        `json:"ts"`
}

func (t *Token) IsZero() bool {
	return t.Owner == 0 && t.Ts.IsZero()
}

func (t *Token) HasExpired() bool {
	return time.Since(t.Ts) < MAX_TOKEN_AGE
}

func (t *Token) IsOwnerInternal() bool {
	return t.Owner == config.Env.BankId
}
