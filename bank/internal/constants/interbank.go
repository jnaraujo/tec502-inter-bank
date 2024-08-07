package constants

import "time"

// Constantes para o tempo de espera de operações interbancárias
const (
	OperationTimeout                  = 250 * time.Millisecond // Tempo máximo que o banco espera por uma resposta de outro banco para uma operação
	MaxTimeToProcessLocalTransactions = 10 * time.Second       // Tempo máximo que o banco tem para processar transações locais
	MaxWaitTimeForTokenInterBank      = 15 * time.Second       // Tempo máximo que o banco espera por um token de outro banco quando ele deveria ser o próximo
	MaxInterBankRequestTimeout        = 250 * time.Millisecond // Tempo máximo que o banco espera por uma resposta de outro banco
)

// Constantes para o tempo de espera de tokens
const (
	MaxTimeToRequestToken = 250 * time.Millisecond
)

// Constantes para o tempo de espera de atualização de rede
const (
	NetworkUpdateWaitDuration = 50 * time.Millisecond
)
