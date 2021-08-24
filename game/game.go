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

	// Sets lists any IDs of Piles to be captured or built with.
	// For face cards, each Pile must be in its own set.
	// For number cards, the values of the Piles in each set must sum to the
	// value of the hand card.
	Sets [][]int

	// Build reports whether the Action creates or modifies a build.
	Build bool
}

// A Pile contains cards on the table.
type Pile struct {
	// Cards lists the cards in the Pile.
	// A build contains two or more cards.
	Cards []card.Card

	// Value is the Pile's numerical value.
	// Face cards have no value.
	Value int

	// Controller is the player who last played onto the Pile, if it is a build.
	Controller int
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
		g.addCardPile(c)
	}
	g.deck = g.deck[4:]

	for i := range g.players {
		piles := make(map[int]Pile, len(g.piles))
		for id, p := range g.piles {
			piles[id] = copyPile(p)
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
				piles[id] = copyPile(p)
			}
			a := p.Play(piles)
			if err := g.validateAction(i, hand[i], a); err != nil {
				panic(err)
			}
			switch {
			case len(a.Sets) == 0:
				// Trail
				delete(hand[i], a.Card)
				g.addCardPile(a.Card)
			case a.Build:
				p := Pile{Value: a.Card.Rank(), Controller: i}
				for _, set := range a.Sets {
					for _, id := range set {
						p.Cards = append(p.Cards, g.piles[id].Cards...)
						delete(g.piles, id)
					}
				}
				p.Cards = append(p.Cards, a.Card)
				delete(hand[i], a.Card)
				g.addPile(p)
			default:
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
}

// validateAction checks whether an Action is valid.
func (g *game) validateAction(player int, hand map[card.Card]bool, a Action) error {
	if !hand[a.Card] {
		return fmt.Errorf("invalid card %v", a.Card)
	}
	if len(a.Sets) == 0 {
		// Trail
		for _, p := range g.piles {
			if len(p.Cards) > 1 && p.Controller == player {
				return fmt.Errorf("cannot trail while building")
			}
		}
		return nil
	}

	// Each ID must be valid and used only once
	ids := make(map[int]bool)
	for _, set := range a.Sets {
		for _, id := range set {
			if _, ok := g.piles[id]; !ok {
				return fmt.Errorf("invalid pile %v", id)
			}
			if ids[id] {
				return fmt.Errorf("duplicate pile %v", id)
			}
			ids[id] = true
		}
	}

	// Face card sets must have exactly one card of matching rank
	if a.Card.IsFace() {
		for _, set := range a.Sets {
			if len(set) != 1 {
				return fmt.Errorf("invalid set %v using %v", set, a.Card)
			}
			if c := g.piles[set[0]].Cards[0]; c.Rank() != a.Card.Rank() {
				return fmt.Errorf("invalid capture %v using %v", c, a.Card)
			}
		}
		return nil
	}

	// Number card sets must have the correct sum and contain no face cards
	for _, set := range a.Sets {
		var sum int
		for _, id := range set {
			value := g.piles[id].Value
			if value == 0 {
				return fmt.Errorf("invalid pile %v using %v", g.piles[id], a.Card)
			}
			sum += value
		}
		if sum != a.Card.Rank() {
			return fmt.Errorf("invalid set %v (sum %v) using %v", set, sum, a.Card)
		}
	}

	if !a.Build {
		// Valid capture
		return nil
	}

	// Builds must have a card in hand that can capture
	for c := range hand {
		if c == a.Card {
			continue
		}
		if c.Rank() == a.Card.Rank() {
			// Valid build
			return nil
		}
	}
	return fmt.Errorf("no card to capture build")
}

// addCardPile adds a new Pile containing a single card to the table.
func (g *game) addCardPile(c card.Card) {
	var value int
	if !c.IsFace() {
		value = c.Rank()
	}
	g.addPile(Pile{Cards: []card.Card{c}, Value: value})
}

// addPile adds a new Pile to the table.
func (g *game) addPile(p Pile) {
	g.npiles++
	g.piles[g.npiles] = p
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

// copyPile returns a Pile deeply equal to p that does not share memory with p.
func copyPile(p Pile) Pile {
	return Pile{
		append([]card.Card{}, p.Cards...),
		p.Value,
		p.Controller,
	}
}
