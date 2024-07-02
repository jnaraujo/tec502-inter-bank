package storage

import (
	"sync"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/interbank"
)

type ringData struct {
	Id   interbank.BankId
	Addr string
}

// implementação de um token ring para
// comunicação entre os bancos
type ringStorage struct {
	mu   sync.RWMutex
	ring []ringData
}

var Ring = &ringStorage{}

// Adiciona um banco ao anel de comunicação
func (r *ringStorage) Add(bankId interbank.BankId, addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ring = append(r.ring, ringData{Id: bankId, Addr: addr})
}

func (r *ringStorage) Find(bankId interbank.BankId) *ringData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, bank := range r.ring {
		if bank.Id == bankId {
			return &bank
		}
	}

	return nil
}

// Retorna o próximo banco no anel de comunicação
func (r *ringStorage) Next(bankId interbank.BankId) *ringData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i, bank := range r.ring {
		if bank.Id == bankId {
			if i+1 == len(r.ring) {
				return &r.ring[0]
			}
			return &r.ring[i+1]
		}
	}

	return nil
}

// Retorna o banco anterior no anel de comunicação
func (r *ringStorage) Before(bankId interbank.BankId) *ringData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i, bank := range r.ring {
		if bank.Id == bankId {
			if i == 0 {
				return &r.ring[len(r.ring)-1]
			}
			return &r.ring[i-1]
		}
	}

	return nil
}

// Retorna o anel de comunicação
func (r *ringStorage) List() []ringData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ring
}

func (r *ringStorage) FindBankWithLowestId() *ringData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lowest := r.ring[0]

	for _, bank := range r.ring {
		if int(bank.Id) < int(lowest.Id) {
			lowest = bank
		}
	}

	return &lowest
}
