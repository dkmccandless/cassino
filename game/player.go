package game

import "github.com/dkmccandless/cassino/card"

// A Player can participate in a game of Cassino.
type Player interface {
	// Init informs the Player of the initial state of the game.
	// pos is the Player's position in the order of play.
	Init(pos int, table []card.Card)

	// Hand supplies a new hand of four cards.
	Hand(hand []card.Card)

	// Play reports the card the Player plays and any captured from the table,
	// grouped into subsets of equal value.
	Play(table []card.Card) (card card.Card, captured [][]card.Card)
}
