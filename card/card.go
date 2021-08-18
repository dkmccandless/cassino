// Package card defines cards for playing Cassino.
package card

// A Card is a playing card in a standard 52-card deck. Cards are ordered first
// by rank, then by suit in the order clubs, diamonds, hearts, spades (e.g. ace
// of clubs = 0, ace of diamonds = 1, king of spades = 51).
type Card int

const (
	// LittleCassino is the two of spades.
	LittleCassino Card = 7

	// BigCassino is the ten of diamonds.
	BigCassino Card = 37
)

// Rank returns a card's rank. Face cards have ranks 11-13.
func (c Card) Rank() int { return int(c/4 + 1) }

// IsAce reports whether a card is an ace.
func (c Card) IsAce() bool { return c.Rank() == 1 }

// IsSpade reports whether a card is a spade.
func (c Card) IsSpade() bool { return c%4 == 3 }

// String returns a string representation of a card.
func (c Card) String() string {
	return `♣♦♥♠`[c%4*3:c%4*3+3] + `A23456789TJQK`[c.Rank()-1:c.Rank()]
}
