package game

import (
	"testing"

	"github.com/dkmccandless/cassino/card"
)

func mapCards(cards ...card.Card) map[card.Card]bool {
	m := make(map[card.Card]bool, len(cards))
	for _, c := range cards {
		m[c] = true
	}
	return m
}

func TestValidateCapture(t *testing.T) {
	for name, test := range map[string]struct {
		table    map[card.Card]bool
		c        card.Card
		captured [][]card.Card
		isErr    bool
	}{
		"face empty capture": {
			mapCards(41), 40, [][]card.Card{{}}, true,
		},
		"face captures with empty capture": {
			mapCards(41, 42, 43), 40, [][]card.Card{{41}, {42}, {}}, true,
		},
		"face captured card not on table": {
			mapCards(41), 40, [][]card.Card{{40}}, true,
		},
		"face captures with captured card not on table": {
			mapCards(41, 42), 40, [][]card.Card{{41}, {42}, {43}}, true,
		},
		"face duplicate capture": {
			mapCards(41), 40, [][]card.Card{{41}, {41}}, true,
		},
		"face captures with duplicate capture": {
			mapCards(41, 42), 40, [][]card.Card{{41}, {42}, {42}}, true,
		},
		"face single": {
			mapCards(41), 40, [][]card.Card{{41}}, false,
		},
		"face multiple": {
			mapCards(41, 42, 43), 40, [][]card.Card{{41}, {42}, {43}}, false,
		},
		"face multiple in same set": {
			mapCards(41, 42, 43), 40, [][]card.Card{{41, 42, 43}}, true,
		},
		"face captures with multiple in same set": {
			mapCards(41, 42, 43), 40, [][]card.Card{{41}, {42, 43}}, true,
		},
		"face wrong rank": {
			mapCards(44), 40, [][]card.Card{{44}}, true,
		},
		"face captures with wrong rank": {
			mapCards(41, 44), 40, [][]card.Card{{41}, {44}}, true,
		},

		"number empty capture": {
			mapCards(1), 0, [][]card.Card{{}}, true,
		},
		"number captures with empty capture": {
			mapCards(1, 2, 3), 0, [][]card.Card{{1}, {2}, {}}, true,
		},
		"number captured card not on table": {
			mapCards(1), 0, [][]card.Card{{0}}, true,
		},
		"number captures with captured card not on table": {
			mapCards(1, 2), 0, [][]card.Card{{1}, {2}, {3}}, true,
		},
		"number wrong rank": {
			mapCards(4), 0, [][]card.Card{{4}}, true,
		},
		"number captures with wrong rank": {
			mapCards(1, 4), 0, [][]card.Card{{1}, {4}}, true,
		},
		"number wrong sum": {
			mapCards(0, 28), 36, [][]card.Card{{0, 28}}, true,
		},
		"number captures with wrong sum": {
			mapCards(0, 28, 37), 36, [][]card.Card{{37}, {0, 28}}, true,
		},
		"number duplicate capture": {
			mapCards(1), 0, [][]card.Card{{1}, {1}}, true,
		},
		"number captures with duplicate capture": {
			mapCards(1, 2), 0, [][]card.Card{{1}, {2}, {2}}, true,
		},
		"number duplicate sum": {
			mapCards(0, 4, 24), 36, [][]card.Card{{0, 4, 24}, {0, 4, 24}}, true,
		},
		"number duplicate card in sums": {
			mapCards(0, 32, 33), 36, [][]card.Card{{0, 32}, {0, 33}}, true,
		},
		"number single": {
			mapCards(1), 0, [][]card.Card{{1}}, false,
		},
		"number multiple": {
			mapCards(1, 2, 3), 0, [][]card.Card{{1}, {2}, {3}}, false,
		},
		"number sum": {
			mapCards(0, 4, 24), 36, [][]card.Card{{0, 4, 24}}, false,
		},
		"number multiple sum": {
			mapCards(0, 4, 28, 32), 36, [][]card.Card{{0, 32}, {4, 28}}, false,
		},
		"number mixed": {
			mapCards(0, 4, 28, 32, 37), 36, [][]card.Card{{37}, {0, 32}, {4, 28}}, false,
		},
	} {
		err := (&game{table: test.table}).validateCapture(test.c, test.captured)
		if isErr := err != nil; isErr != test.isErr {
			t.Errorf("validateCapture(%q): got err=%v, expected %v",
				name, isErr, test.isErr,
			)
		}
	}
}
