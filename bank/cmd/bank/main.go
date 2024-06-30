package main

import (
	"flag"
	"fmt"
	"log/slog"

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

	slog.Info(fmt.Sprintf("Bank App - Id: %s", config.Env.BankId))

	storage.Ring.Add(interbank.NewBankId(1), "localhost:3001")
	storage.Ring.Add(interbank.NewBankId(2), "localhost:3002")
	storage.Ring.Add(interbank.NewBankId(3), "localhost:3003")

	signal := make(chan bool)

	// O servidor deve ser rodado antes do token ring
	// pois é necessário acessar algumas rotas para
	// configurar corretamente
	go func() {
		// run http server on background
		err = http.NewServer(config.Env.ServerPort)
		if err != nil {
			signal <- false
		}
	}()

	// Inicia e configura o token ring
	services.SetupTokenRing()

	// inicia o processamento de transações em background
	transaction_processor.BackgroundJob()

	// aguarda o sinal de encerramento
	<-signal
}
