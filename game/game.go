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

	// keep contains the cards captured by each player.
	keep [][]card.Card

	// deck contains the cards not yet dealt.
	deck []card.Card

	// table contains the cards on the table.
	table map[card.Card]bool

	// lastCapture records who played the most recent capture.
	lastCapture int
}

// An Action describes the action a player takes on their turn.
type Action struct {
	// Card is the player's hand card.
	Card card.Card

	// Sets lists any sets of cards on the table to be captured.
	// For face cards, each must be in its own set.
	// For cards with numerical value, the cards in each set must sum to the
	// value of the capturing card.
	Sets [][]card.Card
}

// Play plays a game of Cassino.
func Play(p0, p1 Player) (score []int) {
	g := &game{
		players: []Player{p0, p1},
		keep:    make([][]card.Card, 2),
		table:   make(map[card.Card]bool),
	}
	for _, v := range rand.Perm(52) {
		g.deck = append(g.deck, card.Card(v))
	}
	for _, c := range g.deck[:4] {
		g.table[c] = true
	}
	g.deck = g.deck[4:]

	for i := range g.players {
		table := make([]card.Card, 0, len(g.table))
		for c := range g.table {
			table = append(table, c)
		}
		g.players[i].Init(i, table)
	}
	for len(g.deck) != 0 {
		g.playHand()
	}
	for c := range g.table {
		g.capture(g.lastCapture, c)
	}

	score = make([]int, 2)
	for i := range g.players {
		if len(g.keep[i]) > 26 {
			score[i] += 3
		}
		var spades int
		for _, c := range g.keep[i] {
			if c.IsSpade() {
				spades++
			}
			switch {
			case c == card.BigCassino:
				score[i] += 2
			case c == card.LittleCassino, c.IsAce():
				score[i]++
			}
		}
		if spades >= 7 {
			score[i]++
		}
	}
	return score
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
			table := make([]card.Card, 0, len(g.table))
			for c := range g.table {
				table = append(table, c)
			}
			a := p.Play(table)
			if !hand[i][a.Card] {
				panic(fmt.Sprintf("invalid card %v", a.Card))
			}
			if len(a.Sets) == 0 {
				// Trail
				delete(hand[i], a.Card)
				g.table[a.Card] = true
				continue
			}
			if err := g.validateAction(a); err != nil {
				panic(err)
			}
			for _, set := range a.Sets {
				for _, cc := range set {
					g.capture(i, cc)
				}
			}
			delete(hand[i], a.Card)
			g.keep[i] = append(g.keep[i], a.Card)
			g.lastCapture = i
		}
	}
}

// validateAction checks whether an Action is valid.
func (g *game) validateAction(a Action) error {
	// cards records the table cards involved the Action.
	cards := make(map[card.Card]bool)
	switch {
	case a.Card.IsFace():
		for _, set := range a.Sets {
			if len(set) != 1 {
				return fmt.Errorf("invalid capture %v using %v", set, a.Card)
			}
			cc := set[0]
			if !g.table[cc] || a.Card.Rank() != cc.Rank() {
				return fmt.Errorf("invalid capture %v using %v", set, a.Card)
			}
			if cards[cc] {
				return fmt.Errorf("duplicate card %v", cc)
			}
			cards[cc] = true
		}
	default:
		for _, set := range a.Sets {
			if len(set) == 0 {
				return fmt.Errorf("invalid capture %v using %v", set, a.Card)
			}
			var sum int
			for _, cc := range set {
				if !g.table[cc] || cc.IsFace() {
					return fmt.Errorf("invalid capture %v using %v", set, a.Card)
				}
				if cards[cc] {
					return fmt.Errorf("duplicate card %v", cc)
				}
				cards[cc] = true
				sum += cc.Rank()
			}
			if sum != a.Card.Rank() {
				return fmt.Errorf("invalid capture %v using %v", set, a.Card)
			}
		}
	}
	return nil
}

// capture moves a card on the table into a player's keep.
func (g *game) capture(player int, c card.Card) {
	delete(g.table, c)
	g.keep[player] = append(g.keep[player], c)
}
