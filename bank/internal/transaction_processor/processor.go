package transaction_processor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank/ibt"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func Process() {
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

		ibt.Process(*tr)
	}

	// envia a transação para o próximo banco
	nextBank := storage.Ring.Next(config.Env.BankId)
	if nextBank == nil {
		panic("Não tem banco no Token Ring!")
	}

	// simula um tempo de processamento
	// TODO: remove this
	time.Sleep(2 * time.Second)

	fmt.Println("Sending ring to", nextBank.Id)
	sendRing(nextBank.Id)
}

func sendRing(to interbank.BankId) {
	// se o banco atual for o banco destino
	// então o anel já deu a volta.
	// talvez todos os bancos estejam fora do ar
	// ou o banco atual seja o único banco no anel (talvez ele mesmo esteja fora do ar)
	if to == config.Env.BankId {
		// não tem mais para onde enviar
		fmt.Println("O token ring deu a volta! Provavelmente todos os bancos estão fora do ar. (Ou eu mesmo estou fora do ar)")
		return
	}

	bank := storage.Ring.Find(to)
	if bank == nil {
		panic("Banco não encontrado no Token Ring!")
	}

	res, err := http.Post("http://"+bank.Addr+"/interbank/token", "application/json", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Println("Error sending ring to", bank.Id)
		// banco esta fora do ar
		// tenta enviar para o próximo banco
		nextBank := storage.Ring.Next(to)
		if nextBank == nil {
			panic("Não tem banco no Token Ring!")
		}
		sendRing(nextBank.Id)
		return
	}
	fmt.Println("Ring sent to", bank.Addr)
}
