package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	Id        int             `json:"id" validate:"required"`
	Name      string          `json:"name" validate:"required,lte=255"`
	CreatedAt time.Time       `json:"created_at"`
	Balance   decimal.Decimal `json:"balance"`
}
