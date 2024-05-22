package config

import (
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvSchema struct {
	ServerPort int              `envconfig:"SERVER_PORT" required:"true"`
	BankId     interbank.BankId `envconfig:"BANK_ID" required:"true"`
}

var Env EnvSchema

func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = envconfig.Process("", &Env)
	if err != nil {
		return err
	}

	return nil
}
