package models

import (
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

type Account struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	Document     string          `json:"document"`
	InterBankKey interbank.IBK   `json:"ibk"`
	CreatedAt    time.Time       `json:"created_at"`
	Balance      decimal.Decimal `json:"balance"`
}
