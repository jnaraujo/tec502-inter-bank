package transaction_processor

import (
	"fmt"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/token_ring"
)

func BackgroundJob() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Checking transactions...", storage.Token.HasToken())
			if storage.Token.HasToken() {
				// TODO: pergunta se tem realmente o token
				processLocalTransactions()
				token_ring.PassToken()
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
			fmt.Printf("Transação %s não existe.\n", id)
			continue
		}

		processTransaction(*tr)
	}
}
