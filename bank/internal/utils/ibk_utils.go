package utils

import (
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
)

func IsLocalUserIBK(ibk interbank.IBK) bool {
	return ibk.BankId == config.Env.BankId
}
