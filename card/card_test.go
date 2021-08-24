package card

import (
	"testing"
)

var cardTests = []struct {
	rank                   int
	isAce, isFace, isSpade bool
	s                      string
}{
	{1, true, false, false, "♣A"},
	{1, true, false, false, "♦A"},
	{1, true, false, false, "♥A"},
	{1, true, false, true, "♠A"},
	{2, false, false, false, "♣2"},
	{2, false, false, false, "♦2"},
	{2, false, false, false, "♥2"},
	{2, false, false, true, "♠2"},
	{3, false, false, false, "♣3"},
	{3, false, false, false, "♦3"},
	{3, false, false, false, "♥3"},
	{3, false, false, true, "♠3"},
	{4, false, false, false, "♣4"},
	{4, false, false, false, "♦4"},
	{4, false, false, false, "♥4"},
	{4, false, false, true, "♠4"},
	{5, false, false, false, "♣5"},
	{5, false, false, false, "♦5"},
	{5, false, false, false, "♥5"},
	{5, false, false, true, "♠5"},
	{6, false, false, false, "♣6"},
	{6, false, false, false, "♦6"},
	{6, false, false, false, "♥6"},
	{6, false, false, true, "♠6"},
	{7, false, false, false, "♣7"},
	{7, false, false, false, "♦7"},
	{7, false, false, false, "♥7"},
	{7, false, false, true, "♠7"},
	{8, false, false, false, "♣8"},
	{8, false, false, false, "♦8"},
	{8, false, false, false, "♥8"},
	{8, false, false, true, "♠8"},
	{9, false, false, false, "♣9"},
	{9, false, false, false, "♦9"},
	{9, false, false, false, "♥9"},
	{9, false, false, true, "♠9"},
	{10, false, false, false, "♣T"},
	{10, false, false, false, "♦T"},
	{10, false, false, false, "♥T"},
	{10, false, false, true, "♠T"},
	{11, false, true, false, "♣J"},
	{11, false, true, false, "♦J"},
	{11, false, true, false, "♥J"},
	{11, false, true, true, "♠J"},
	{12, false, true, false, "♣Q"},
	{12, false, true, false, "♦Q"},
	{12, false, true, false, "♥Q"},
	{12, false, true, true, "♠Q"},
	{13, false, true, false, "♣K"},
	{13, false, true, false, "♦K"},
	{13, false, true, false, "♥K"},
	{13, false, true, true, "♠K"},
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
		if isFace := c.IsFace(); isFace != test.isFace {
			t.Errorf("IsFace(%v): got %v, expected %v",
				c, isFace, test.isFace,
			)
		}
		if isSpade := c.IsSpade(); isSpade != test.isSpade {
			t.Errorf("IsSpade(%v): got %v, expected %v",
				c, isSpade, test.isSpade,
			)
		}
		if s := c.String(); s != test.s {
			t.Errorf("String(%d): got %s, expected %s",
				i, s, test.s,
			)
		}
	}
}
