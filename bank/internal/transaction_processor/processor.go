package transaction_processor

import (
	"fmt"
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
			time.Sleep(250 * time.Millisecond) // espera para verificar as transações

			if storage.Token.HasValidToken() {
				// verifica se o token ainda é do banco atual
				ok := verifyTokenOwnershipNetwork()
				if !ok {
					slog.Info("Verificação falhou")
					continue
				}

				processLocalTransactions()
				services.PassToken() // passa o token para o próximo banco
				continue
			}

			if storage.Token.HasInvalidToken() {
				slog.Info("Token expirado. Passando para o próximo banco...")
				continue
			}

			// se o tempo de espera para o token for excedido
			// o primeiro banco a perceber solicita o token
			// bancos com IDs menores têm prioridade
			bankTokenPriority := time.Duration(math.Pow(2, float64(config.Env.BankId-1))) * time.Second
			maxTokenWaitDuration := constants.MaxWaitTimeForTokenInterBank + bankTokenPriority
			if time.Since(storage.Token.Get().Ts) > maxTokenWaitDuration {
				slog.Info("Tempo de espera para token interbancário excedido. Solicitando token...", "since", time.Since(storage.Token.Get().Ts))
				services.BroadcastToken(config.Env.BankId) // faz um broadcast a todos os bancos avisando que o token agora é do banco atual
			}
		}
	}()
}

func verifyTokenOwnershipNetwork() bool {
	bank := services.RequestTokenFromBanks()
	if bank == nil {
		slog.Warn("Token não encontrado na rede...")
		return false
	}

	if bank.Owner != storage.Token.Get().Owner {
		// pode acontecer da rede estar desatualizada
		// espera um pouco e tenta novamente
		time.Sleep(constants.NetworkUpdateWaitDuration)
		bank = services.RequestTokenFromBanks()
		if bank == nil {
			slog.Warn("Token não encontrado na rede...")
			return false
		}
		// se após a segunda tentativa o token ainda não for do banco atual
		if bank.Owner != storage.Token.Get().Owner {
			return false
		}
	}
	return true
}

func processLocalTransactions() {
	trIds := storage.TransactionQueue.List()
	if len(trIds) == 0 {
		return
	}

	slog.Info("Processing transactions...", "hasToken", storage.Token.HasValidToken())

	start := time.Now()
	// processa as transações
	for _, id := range trIds {
		if time.Since(storage.Token.Get().Ts) > constants.MaxTimeToProcessLocalTransactions {
			slog.Info("Tempo para processar transações locais excedido")
			return
		}
		if !storage.Token.HasValidToken() {
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
	duration := time.Since(start)
	fmt.Println("Process Transactions:", duration, "Txs:", len(trIds))
}
