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
	keep []map[card.Card]bool

	// deck contains the cards not yet dealt.
	deck []card.Card

	// table contains the cards on the table.
	table map[card.Card]bool

	// lastCapture records who played the most recent capture.
	lastCapture int
}

// Play plays a game of Cassino.
func Play(p0, p1 Player) (score []int) {
	g := &game{
		players: []Player{p0, p1},
		keep: []map[card.Card]bool{
			make(map[card.Card]bool),
			make(map[card.Card]bool),
		},
		table: make(map[card.Card]bool),
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
		delete(g.table, c)
		g.keep[g.lastCapture][c] = true
	}

	score = make([]int, 2)
	for i := range g.players {
		if len(g.keep[i]) > 26 {
			score[i] += 3
		}
		var spades int
		for c := range g.keep[i] {
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
			c, captured := p.Play(table)
			if !hand[i][c] {
				panic(fmt.Sprintf("invalid card %v", c))
			}
			if len(captured) == 0 {
				// Trail
				delete(hand[i], c)
				g.table[c] = true
				continue
			}
			if err := g.validateCapture(c, captured); err != nil {
				panic(err)
			}
			for _, set := range captured {
				for _, cc := range set {
					delete(g.table, cc)
					g.keep[i][cc] = true
				}
			}
			delete(hand[i], c)
			g.keep[i][c] = true
			g.lastCapture = i
		}
	}
}

// validateCapture checks whether a capture is valid.
func (g *game) validateCapture(c card.Card, captured [][]card.Card) error {
	// capCards records the cards involved the capture.
	capCards := make(map[card.Card]bool)
	switch {
	case c.IsFace():
		for _, set := range captured {
			if len(set) != 1 {
				return fmt.Errorf("invalid capture %v using %v", set, c)
			}
			cc := set[0]
			if !g.table[cc] || c.Rank() != cc.Rank() {
				return fmt.Errorf("invalid capture %v using %v", set, c)
			}
			if capCards[cc] {
				return fmt.Errorf("duplicate capture %v", cc)
			}
			capCards[cc] = true
		}
	default:
		for _, set := range captured {
			if len(set) == 0 {
				return fmt.Errorf("invalid capture %v using %v", set, c)
			}
			var sum int
			for _, cc := range set {
				if !g.table[cc] || cc.IsFace() {
					return fmt.Errorf("invalid capture %v using %v", set, c)
				}
				if capCards[cc] {
					return fmt.Errorf("duplicate capture %v", cc)
				}
				capCards[cc] = true
				sum += cc.Rank()
			}
			if sum != c.Rank() {
				return fmt.Errorf("invalid capture %v using %v", set, c)
			}
		}
	}
	return nil
}
