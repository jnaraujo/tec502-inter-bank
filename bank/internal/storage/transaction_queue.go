package storage

import (
	"sync"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/models"
)

type tqStorage struct {
	mu    sync.RWMutex
	queue []models.TransactionId
}

var TransactionQueue = &tqStorage{}

// Adiciona uma transação na fila de transações
func (tq *tqStorage) Add(id models.TransactionId) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	tq.queue = append(tq.queue, id)
}

// Retorna a primeira transação da fila de transações
func (tq *tqStorage) Remove(id models.TransactionId) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	for i, t := range tq.queue {
		if t == id {
			tq.queue = append(tq.queue[:i], tq.queue[i+1:]...)
			return
		}
	}
}

// Retorna a fila de transações
func (tq *tqStorage) List() []models.TransactionId {
	tq.mu.RLock()
	defer tq.mu.RUnlock()

	return tq.queue
}
