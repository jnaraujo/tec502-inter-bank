package transaction_processor

import (
	"log/slog"
	"math"
	"time"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/constants"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/services"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/storage"
)

func BackgroundJob() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond) // espera para verificar as transações

			slog.Info("Checking transactions...", "hasToken", storage.Token.HasToken())
			if storage.Token.HasToken() {
				bank := services.RequestTokenFromBanks()
				if bank != nil && bank.Owner != storage.Token.Get().Owner {
					// O sistema tem o token, mas não é o dono
					storage.Token.Set(*bank)
					continue
				}

				processLocalTransactions()
				services.PassToken() // passa o token para o próximo banco
				continue
			}

			token := storage.Token.Get()
			if token.Owner == config.Env.BankId && token.HasExpired() {
				slog.Info("Token expirado. Passando para o próximo banco...")
				services.PassToken() // passa o token para o próximo banco
			}

			// se o tempo de espera para o token for excedido
			// o primeiro banco a perceber solicita o token
			// bancos com IDs menores têm prioridade
			bankTokenPriority := math.Pow(2, float64(config.Env.BankId))
			maxTokenWaitDuration := time.Duration(float64(constants.MaxWaitTimeForTokenInterBank) + bankTokenPriority)
			if time.Since(storage.Token.Get().Ts) > maxTokenWaitDuration {
				slog.Info("Tempo de espera para token interbancário excedido. Solicitando token...")
				services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
			}
		}
	}()
}

func processLocalTransactions() {
	slog.Info("Processing transactions...", "hasToken", storage.Token.HasToken())

	// processa as transações
	trIds := storage.TransactionQueue.List()
	for _, id := range trIds {
		if time.Since(storage.Token.Get().Ts) > constants.MaxTimeToProcessLocalTransactions {
			slog.Info("Tempo para processar transações locais excedido")
			return
		}
		if !storage.Token.HasToken() {
			slog.Info("Token não está mais no banco")
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
