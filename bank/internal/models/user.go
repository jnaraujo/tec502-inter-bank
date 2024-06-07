package models

import (
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

type User struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	CreatedAt    time.Time       `json:"created_at"`
	Balance      decimal.Decimal `json:"balance"`
	InterBankKey interbank.IBK   `json:"ibk"`
}
