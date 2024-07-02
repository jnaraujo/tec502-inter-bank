package transaction_processor

import (
	"log/slog"
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

			token := storage.Token.Get()
			if token.Owner == config.Env.BankId && token.HasExpired() {
				slog.Info("Token expirado. Passando para o próximo banco...")
				services.PassToken() // passa o token para o próximo banco
			}

			nextOwner := storage.Ring.Next(storage.Token.Get().Owner)
			if nextOwner == nil {
				continue
			}

			// se o próximo dono do token for o banco atual e o tempo de espera para o token interbancário for excedido
			if nextOwner.Id == config.Env.BankId && time.Since(storage.Token.Get().Ts) > constants.MaxWaitTimeForTokenInterBank {
				slog.Info("Tempo de espera para token interbancário excedido. Solicitando token...")
				services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
				continue
			}

			before := storage.Ring.Before(storage.Token.Get().Owner) // banco anterior ao dono do token
			if before == nil {
				continue
			}

			// se o tempo de espera para o token interbancário for excedido
			// e o banco anterior ao dono do token for o banco atual
			// o banco atual é responsável passar o token novamente
			maxDuration := time.Duration(float64(constants.MaxWaitTimeForTokenInterBank) * 1.5) // ~22.5s
			if time.Since(storage.Token.Get().Ts) > maxDuration && before.Id == config.Env.BankId {
				slog.Info("Tempo de espera para token interbancário excedido. Solicitando token...")
				services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
				continue
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
