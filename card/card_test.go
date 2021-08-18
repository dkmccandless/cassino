package card

import (
	"testing"
)

var cardTests = []struct {
	rank           int
	isAce, isSpade bool
	s              string
}{
	{1, true, false, "♣A"},
	{1, true, false, "♦A"},
	{1, true, false, "♥A"},
	{1, true, true, "♠A"},
	{2, false, false, "♣2"},
	{2, false, false, "♦2"},
	{2, false, false, "♥2"},
	{2, false, true, "♠2"},
	{3, false, false, "♣3"},
	{3, false, false, "♦3"},
	{3, false, false, "♥3"},
	{3, false, true, "♠3"},
	{4, false, false, "♣4"},
	{4, false, false, "♦4"},
	{4, false, false, "♥4"},
	{4, false, true, "♠4"},
	{5, false, false, "♣5"},
	{5, false, false, "♦5"},
	{5, false, false, "♥5"},
	{5, false, true, "♠5"},
	{6, false, false, "♣6"},
	{6, false, false, "♦6"},
	{6, false, false, "♥6"},
	{6, false, true, "♠6"},
	{7, false, false, "♣7"},
	{7, false, false, "♦7"},
	{7, false, false, "♥7"},
	{7, false, true, "♠7"},
	{8, false, false, "♣8"},
	{8, false, false, "♦8"},
	{8, false, false, "♥8"},
	{8, false, true, "♠8"},
	{9, false, false, "♣9"},
	{9, false, false, "♦9"},
	{9, false, false, "♥9"},
	{9, false, true, "♠9"},
	{10, false, false, "♣T"},
	{10, false, false, "♦T"},
	{10, false, false, "♥T"},
	{10, false, true, "♠T"},
	{11, false, false, "♣J"},
	{11, false, false, "♦J"},
	{11, false, false, "♥J"},
	{11, false, true, "♠J"},
	{12, false, false, "♣Q"},
	{12, false, false, "♦Q"},
	{12, false, false, "♥Q"},
	{12, false, true, "♠Q"},
	{13, false, false, "♣K"},
	{13, false, false, "♦K"},
	{13, false, false, "♥K"},
	{13, false, true, "♠K"},
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
		if s := c.String(); s != test.s {
			t.Errorf("String(%d): got %s, expected %s",
				i, s, test.s,
			)
		}
	}
}
