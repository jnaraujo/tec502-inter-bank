package vectorclock

import (
	"fmt"
	"sync"
)

// VectorClock representa um relógio vetorial.
type VectorClock struct {
	clock map[string]int
	mu    sync.Mutex
}

// NewVectorClock cria um novo relógio vetorial.
func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make(map[string]int),
	}
}

// Increment incrementa o contador para um determinado nó.
func (vc *VectorClock) Increment(nodeID string) {
	vc.mu.Lock()
	vc.clock[nodeID]++
	vc.mu.Unlock()
}

// Update atualiza o relógio vetorial com outro relógio vetorial.
func (vc *VectorClock) Update(other *VectorClock) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for nodeID, count := range other.clock {
		if current, exists := vc.clock[nodeID]; !exists || count > current {
			vc.clock[nodeID] = count
		}
	}
}

// Compare compara este relógio vetorial com outro.
func (vc *VectorClock) Compare(other *VectorClock) int {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	less, greater := false, false

	for nodeID, count := range vc.clock {
		otherCount, exists := other.clock[nodeID]
		if !exists {
			greater = true
			continue
		}
		if count < otherCount {
			less = true
		} else if count > otherCount {
			greater = true
		}
	}

	for nodeID := range other.clock {
		if _, exists := vc.clock[nodeID]; !exists {
			less = true
		}
	}

	if less && !greater {
		return -1 // vc < other
	} else if greater && !less {
		return 1 // vc > other
	} else if !less && !greater {
		return 0 // vc == other
	} else {
		return 2 // vc e other são concorrentes
	}
}

// String retorna uma representação em string do relógio vetorial.
func (vc *VectorClock) String() string {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return fmt.Sprintf("%v", vc.clock)
}
