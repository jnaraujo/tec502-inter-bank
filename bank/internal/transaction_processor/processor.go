package transaction_processor

import (
	"fmt"
	"time"

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
			}
		}
	}()
}

func processLocalTransactions() {
	fmt.Println("Processing transactions")

	// processa as transações
	trIds := storage.TransactionQueue.List()
	for _, id := range trIds {
		fmt.Printf("Processing transaction %s\n", id)
		// remove a transação da fila no final
		defer storage.TransactionQueue.Remove(id)

		tr := storage.Transactions.FindTransactionById(id)
		if tr == nil {
			continue
		}
		processTransaction(*tr)
	}
}
