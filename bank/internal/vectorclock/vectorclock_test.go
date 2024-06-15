package vectorclock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorClock(t *testing.T) {
	vc1 := NewVectorClock()

	vc1.Increment("nodeA")
	vc1.Increment("nodeA")
	vc1.Increment("nodeB")

	assert.Equal(t, 2, vc1.clock["nodeA"])
	assert.Equal(t, 1, vc1.clock["nodeB"])
	assert.Len(t, vc1.clock, 2)
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
