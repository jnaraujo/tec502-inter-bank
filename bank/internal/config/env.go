package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvSchema struct {
	ServerPort int `envconfig:"SERVER_PORT" required:"true"`
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
