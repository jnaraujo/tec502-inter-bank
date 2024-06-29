package transaction_processor

import (
	"fmt"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/constants"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func BackgroundJob() {
	go func() {
		for {
			time.Sleep(1 * time.Second) // espera 1 segundo para verificar as transações

			fmt.Println("Checking transactions...", storage.Token.HasToken())
			if storage.Token.HasToken() {
				bank := services.AskBankWithToken()
				if bank != nil && bank.Owner != storage.Token.Get().Owner {
					// O sistema tem o token, mas não é o dono
					storage.Token.Set(*bank)
					continue
				}

				processLocalTransactions()
				services.PassToken() // passa o token para o próximo banco

				continue
			}

			nextOwner := storage.Ring.Next(storage.Token.Get().Owner)
			if nextOwner == nil {
				continue
			}

			// se o próximo dono do token for o banco atual e o tempo de espera para o token interbancário for excedido
			if nextOwner.Id == config.Env.BankId && time.Since(storage.Token.Get().Ts) > constants.MaxWaitTimeForTokenInterBank {
				fmt.Println("Tempo de espera para token interbancário excedido. Solicitando token...")
				services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
			}
		}
	}()
}

func processLocalTransactions() {
	fmt.Println("Processing transactions")

	// processa as transações
	trIds := storage.TransactionQueue.List()
	for _, id := range trIds {
		if time.Since(storage.Token.Get().Ts) > constants.MaxTimeToProcessLocalTransactions {
			fmt.Println("Tempo para processar transações locais excedido")
			return
		}
		if !storage.Token.HasToken() {
			fmt.Println("Token não é mais do banco")
			return
		}
		// remove a transação da fila no final
		defer storage.TransactionQueue.Remove(id)

		tr := storage.Transactions.FindTransactionById(id)
		if tr == nil {
			continue
		}
		services.ProcessTransaction(*tr)
	}
}
