package models

import (
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeIndividual AccountType = "individual" // Pessoa física
	AccountTypeLegal      AccountType = "legal"      // Pessoa jurídica
	AccountTypeJoint      AccountType = "joint"      // Conta conjunta
)

type Account struct {
	Id             int             `json:"id"`
	Name           string          `json:"name"`
	Documents      []string        `json:"documents"`
	Type           AccountType     `json:"type"`
	InterBankKey   interbank.IBK   `json:"ibk"`
	CreatedAt      time.Time       `json:"created_at"`
	Balance        decimal.Decimal `json:"balance"`
	PendingBalance decimal.Decimal `json:"-"`
	BlockedBalance decimal.Decimal `json:"-"`
}
