package main

import (
	"flag"
	"fmt"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/http"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/transaction_processor"
)

func main() {
	err := config.InitEnv()
	if err != nil {
		panic(err)
	}

	var port int
	var bankId int

	flag.IntVar(&port, "port", config.Env.ServerPort, "Server port")
	flag.IntVar(&bankId, "id", int(config.Env.BankId), "Bank id")
	flag.Parse()

	config.Env.ServerPort = port
	config.Env.BankId = interbank.NewBankId(uint16(bankId))

	fmt.Printf("Bank App - Id: %s\n", config.Env.BankId)

	storage.Ring.Add(interbank.NewBankId(1), "localhost:3001")
	storage.Ring.Add(interbank.NewBankId(2), "localhost:3002")

	if storage.Ring.FindBankWithLowestId().Id == config.Env.BankId {
		fmt.Println("I'm the bank with the lowest id")
		go func() {
			transaction_processor.Process()
		}()
	}

	err = http.NewServer(config.Env.ServerPort)
	if err != nil {
		panic(err)
	}
}
