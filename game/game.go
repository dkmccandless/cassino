// Package game implements the game of Cassino.
package game

import (
	"fmt"
	"math/rand"

	"github.com/dkmccandless/cassino/card"
)

// A game administers a single complete game.
type game struct {
	// players records the Players in order.
	players []Player

	// deck contains the cards not yet dealt.
	deck []card.Card

	// table contains the cards on the table.
	table []card.Card
}

// Play plays a game of Cassino.
func Play(p0, p1 Player) {
	g := &game{players: []Player{p0, p1}}
	for _, v := range rand.Perm(52) {
		g.deck = append(g.deck, card.Card(v))
	}
	g.table = append(g.table, g.deck[:4]...)
	g.deck = g.deck[4:]

	for i := range g.players {
		g.players[i].Init(i, append([]card.Card{}, g.table...))
	}
	for len(g.deck) != 0 {
		g.playHand()
	}
}

// playHand deals and plays a single four-card hand.
func (g *game) playHand() {
	hand := []map[card.Card]bool{
		make(map[card.Card]bool, 4),
		make(map[card.Card]bool, 4),
	}
	for i, p := range g.players {
		for _, c := range g.deck[:4] {
			hand[i][c] = true
		}
		p.Hand(append([]card.Card{}, g.deck[:4]...))
		g.deck = g.deck[4:]
	}

	for len(hand[0]) != 0 {
		for i, p := range g.players {
			c := p.Play()
			if !hand[i][c] {
				panic(fmt.Sprintf("invalid card %v", c))
			}
			g.table = append(g.table, c)
			delete(hand[i], c)
		}
	}
}
