package game

import "github.com/dkmccandless/cassino/card"

// A Player can participate in a game of Cassino.
type Player interface {
	// Init informs the Player of the initial state of the game.
	// pos is the Player's position in the order of play.
	Init(pos int, piles map[int]Pile)

	// Hand supplies a new hand of four cards.
	Hand(hand []card.Card)

	// Note informs the Player of cards their opponent plays and captures.
	// If captured is not empty, it includes the played card.
	Note(played card.Card, captured []card.Card)

	// Play reports the Action the player takes on their turn.
	Play(piles map[int]Pile) Action
}
