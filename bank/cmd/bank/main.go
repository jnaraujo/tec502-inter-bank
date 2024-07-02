package main

import (
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

	slog.Info(fmt.Sprintf("Bank App - Id: %s", config.Env.BankId))

	if len(config.Env.Banks) == 0 {
		panic("No banks configured")
	}

	for idx, addr := range config.Env.Banks {
		storage.Ring.Add(interbank.NewBankId(uint16(idx)+1), addr)
	}

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
