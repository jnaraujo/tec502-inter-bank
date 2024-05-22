package main

import (
	"fmt"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/http"
)

func main() {
	err := config.InitEnv()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Bank App - Id: %s\n", config.Env.BankId)

	err = http.NewServer()
	if err != nil {
		panic(err)
	}
}
