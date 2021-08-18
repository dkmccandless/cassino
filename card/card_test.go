package card

import (
	"testing"
)

var cardTests = []struct {
	rank           int
	isAce, isSpade bool
}{
	{1, true, false},
	{1, true, false},
	{1, true, false},
	{1, true, true},
	{2, false, false},
	{2, false, false},
	{2, false, false},
	{2, false, true},
	{3, false, false},
	{3, false, false},
	{3, false, false},
	{3, false, true},
	{4, false, false},
	{4, false, false},
	{4, false, false},
	{4, false, true},
	{5, false, false},
	{5, false, false},
	{5, false, false},
	{5, false, true},
	{6, false, false},
	{6, false, false},
	{6, false, false},
	{6, false, true},
	{7, false, false},
	{7, false, false},
	{7, false, false},
	{7, false, true},
	{8, false, false},
	{8, false, false},
	{8, false, false},
	{8, false, true},
	{9, false, false},
	{9, false, false},
	{9, false, false},
	{9, false, true},
	{10, false, false},
	{10, false, false},
	{10, false, false},
	{10, false, true},
	{11, false, false},
	{11, false, false},
	{11, false, false},
	{11, false, true},
	{12, false, false},
	{12, false, false},
	{12, false, false},
	{12, false, true},
	{13, false, false},
	{13, false, false},
	{13, false, false},
	{13, false, true},
}

func TestCard(t *testing.T) {
	for i, test := range cardTests {
		c := Card(i)
		if rank := c.Rank(); rank != test.rank {
			t.Errorf("Rank(%v): got %v, expected %v",
				c, rank, test.rank,
			)
		}
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
