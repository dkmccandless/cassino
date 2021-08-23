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

	// piles contains the cards on the table.
	piles map[int]Pile

	// npiles records how many Piles have been added.
	npiles int

	// lastCapture records who played the most recent capture.
	lastCapture int
}

// An Action describes the action a player takes on their turn.
type Action struct {
	// Card is the player's hand card.
	Card card.Card

	// Sets lists any IDs of Piles to be captured.
	// For face cards, each Pile must be in its own set.
	// For number cards, the values of the Piles in each set must sum to the
	// value of the capturing card.
	Sets [][]int
}

// A Pile contains cards on the table.
type Pile struct {
	// Cards records the cards in the Pile.
	Cards []card.Card

	// Value is the Pile's numerical value.
	// For a number card, this is the card's value.
	// Face cards have no value.
	Value int
}

// Play plays a game of Cassino and returns the final score.
func Play(p0, p1 Player) []int {
	g := &game{
		players: []Player{p0, p1},
		keep:    make([][]card.Card, 2),
		piles:   make(map[int]Pile),
	}
	for _, v := range rand.Perm(52) {
		g.deck = append(g.deck, card.Card(v))
	}
	for _, c := range g.deck[:4] {
		g.addPile(c)
	}
	g.deck = g.deck[4:]

	for i := range g.players {
		piles := make(map[int]Pile, len(g.piles))
		for id, p := range g.piles {
			piles[id] = Pile{
				Cards: append([]card.Card{}, p.Cards...),
				Value: p.Value,
			}
		}
		g.players[i].Init(i, piles)
	}
	for len(g.deck) != 0 {
		g.playHand()
	}
	for id := range g.piles {
		g.capture(g.lastCapture, id)
	}

	return []int{score(g.keep[0]), score(g.keep[1])}
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
			piles := make(map[int]Pile, len(g.piles))
			for id, p := range g.piles {
				piles[id] = Pile{
					Cards: append([]card.Card{}, p.Cards...),
					Value: p.Value,
				}
			}
			a := p.Play(piles)
			if !hand[i][a.Card] {
				panic(fmt.Sprintf("invalid card %v", a.Card))
			}
			if len(a.Sets) == 0 {
				// Trail
				delete(hand[i], a.Card)
				g.addPile(a.Card)
				continue
			}
			if err := g.validateAction(a); err != nil {
				panic(err)
			}
			for _, set := range a.Sets {
				for _, id := range set {
					g.capture(i, id)
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
	// ids records the IDs of Piles involved the Action.
	ids := make(map[int]bool)
	switch {
	case a.Card.IsFace():
		for _, set := range a.Sets {
			if len(set) != 1 {
				return fmt.Errorf("invalid set %v using %v", set, a.Card)
			}
			id := set[0]
			p, ok := g.piles[id]
			if !ok {
				return fmt.Errorf("invalid Pile ID %v", id)
			}
			if ids[id] {
				return fmt.Errorf("duplicate Pile %v", id)
			}
			ids[id] = true
			if len(p.Cards) != 1 || p.Value != 0 {
				return fmt.Errorf("invalid Pile %+v", p)
			}
			if a.Card.Rank() != p.Cards[0].Rank() {
				return fmt.Errorf("invalid capture %v using %v", p.Cards[0], a.Card)
			}
		}
	default:
		for _, set := range a.Sets {
			if len(set) == 0 {
				return fmt.Errorf("invalid set %v using %v", set, a.Card)
			}
			var sum int
			for _, id := range set {
				p, ok := g.piles[id]
				if !ok {
					return fmt.Errorf("invalid Pile ID %v", id)
				}
				if ids[id] {
					return fmt.Errorf("duplicate Pile %v", id)
				}
				ids[id] = true
				if len(p.Cards) == 0 || p.Value == 0 {
					return fmt.Errorf("invalid Pile %+v", p)
				}
				sum += p.Value
			}
			if sum != a.Card.Rank() {
				return fmt.Errorf("invalid capture %v using %v", set, a.Card)
			}
		}
	}
	return nil
}

// addPile adds a new Pile to the table.
func (g *game) addPile(c card.Card) {
	g.npiles++
	g.piles[g.npiles] = Pile{Cards: []card.Card{c}, Value: c.Value()}
}

// capture moves a Pile into a player's keep.
func (g *game) capture(player int, id int) {
	g.keep[player] = append(g.keep[player], g.piles[id].Cards...)
	delete(g.piles, id)
}

// score returns the score of a slice of cards.
func score(cards []card.Card) int {
	// Most cards: 3
	// Most spades: 1
	// Big Cassino: 2
	// Little Cassino: 1
	// Each ace: 1
	var n, spades int
	if len(cards) > 26 {
		n += 3
	}
	for _, c := range cards {
		if c.IsSpade() {
			spades++
		}
		switch {
		case c == card.BigCassino:
			n += 2
		case c == card.LittleCassino, c.IsAce():
			n++
		}
	}
	if spades >= 7 {
		n++
	}
	return n
}
