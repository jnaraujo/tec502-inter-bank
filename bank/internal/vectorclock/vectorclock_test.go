package vectorclock

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorClock(t *testing.T) {
	bank1 := NewVectorClock()
	bank2 := NewVectorClock()
	bank3 := NewVectorClock()

	// Cria um evento no banco 1 e no banco 3
	bank1.Increment("bank1")
	bank3.Increment("bank3")

	// banco 1 e banco 2 enviam seus eventos
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		// banco 2 e 3 recebem o evento do banco 1
		bank2.Update(bank1)
		bank3.Update(bank1)
		wg.Done()
	}()
	go func() {
		// banco 3 e 1 recebem o evento do banco 3
		bank1.Update(bank3)
		bank2.Update(bank3)
		wg.Done()
	}()
	go func() {
		// banco 2 cria um evento
		bank2.Increment("bank2")

		// banco 1 e 3 recebem o evento do banco 2
		bank1.Update(bank2)
		bank3.Update(bank2)
		wg.Done()
	}()
	wg.Wait()

	// banco 2 cria dois novos evento
	bank2.Increment("bank2")
	bank2.Increment("bank2")

	// banco 1 e 3 recebem o evento do banco 2
	bank1.Update(bank2)
	bank3.Update(bank2)
	expected := map[string]int{
		"bank1": 1,
		"bank2": 3,
		"bank3": 1,
	}

	assert.Equal(t, expected, bank1.clock)
	assert.Equal(t, expected, bank2.clock)
	assert.Equal(t, expected, bank3.clock)
}

func TestVectorClockUpdate(t *testing.T) {
	vc1 := NewVectorClock()
	vc2 := NewVectorClock()

	vc1.Increment("nodeA")
	vc1.Increment("nodeA")

	vc2.Increment("nodeA")
	vc2.Increment("nodeB")

	vc1.Update(vc2)

	assert.Equal(t, 2, vc1.clock["nodeA"])
	assert.Equal(t, 1, vc1.clock["nodeB"])
	assert.Len(t, vc1.clock, 2)
}

func TestVectorClockCompare(t *testing.T) {
	vc1 := NewVectorClock()
	vc2 := NewVectorClock()

	vc1.Increment("nodeA")
	vc1.Increment("nodeA")

	vc2.Increment("nodeA")
	vc2.Increment("nodeB")

	assert.Equal(t, 2, vc1.Compare(vc2))
}
