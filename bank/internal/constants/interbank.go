package constants

import "time"

const (
	OperationTimeout                  = 250 * time.Millisecond // Tempo máximo que o banco espera por uma resposta de outro banco para uma operação
	MaxTimeToProcessLocalTransactions = 10 * time.Second       // Tempo máximo que o banco tem para processar transações locais
	MaxWaitTimeForTokenInterBank      = 15 * time.Second       // Tempo máximo que o banco espera por um token de outro banco quando ele deveria ser o próximo
)
