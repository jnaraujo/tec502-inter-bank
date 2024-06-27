package main

import (
	"flag"
	"fmt"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/http"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
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

	signal := make(chan bool)

	go func() {
		// run http server on background
		err = http.NewServer(config.Env.ServerPort)
		if err != nil {
			signal <- false
		}
	}()

	if storage.Ring.FindBankWithLowestId().Id == config.Env.BankId {
		fmt.Println("I'm the bank with the lowest id")
		// verifica se o token já esta na rede.
		if !services.IsTokenOnRing() {
			// se não estiver, cria o token
			services.BroadcastToken(config.Env.BankId)
		}
	}

	bank := services.AskBankWithToken()
	if bank != nil && bank.Owner == config.Env.BankId {
		storage.Token.Set(*bank)
	}

	// inicia o processamento de transações em background
	transaction_processor.BackgroundJob()

	// aguarda o sinal de encerramento
	<-signal
}
