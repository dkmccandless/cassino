package card

import (
	"testing"
)

var cardTests = []struct {
	isAce, isSpade bool
}{
	{true, false},
	{true, false},
	{true, false},
	{true, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
	{false, false},
	{false, false},
	{false, false},
	{false, true},
}

func TestCard(t *testing.T) {
	for i, test := range cardTests {
		c := Card(i)
		if isAce := c.IsAce(); isAce != test.isAce {
			t.Errorf("IsAce(%v): got %v, expected %v",
				c, isAce, test.isAce,
			)
		}
		if isSpade := c.IsSpade(); isSpade != test.isSpade {
			t.Errorf("IsSpade(%v): got %v, expected %v",
				c, isSpade, test.isSpade,
			)
		}
	}
}
