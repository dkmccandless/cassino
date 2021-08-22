package game

import (
	"reflect"
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

func TestValidateAction(t *testing.T) {
	for name, test := range map[string]struct {
		table map[card.Card]bool
		a     Action
		isErr bool
	}{
		"face empty capture": {
			mapCards(41),
			Action{40, [][]card.Card{{}}},
			true,
		},
		"face captures with empty capture": {
			mapCards(41, 42, 43),
			Action{40, [][]card.Card{{41}, {42}, {}}},
			true,
		},
		"face captured card not on table": {
			mapCards(41),
			Action{40, [][]card.Card{{40}}},
			true,
		},
		"face captures with captured card not on table": {
			mapCards(41, 42),
			Action{40, [][]card.Card{{41}, {42}, {43}}},
			true,
		},
		"face duplicate capture": {
			mapCards(41),
			Action{40, [][]card.Card{{41}, {41}}},
			true,
		},
		"face captures with duplicate capture": {
			mapCards(41, 42),
			Action{40, [][]card.Card{{41}, {42}, {42}}},
			true,
		},
		"face single": {
			mapCards(41),
			Action{40, [][]card.Card{{41}}},
			false,
		},
		"face multiple": {
			mapCards(41, 42, 43),
			Action{40, [][]card.Card{{41}, {42}, {43}}},
			false,
		},
		"face multiple in same set": {
			mapCards(41, 42, 43),
			Action{40, [][]card.Card{{41, 42, 43}}},
			true,
		},
		"face captures with multiple in same set": {
			mapCards(41, 42, 43),
			Action{40, [][]card.Card{{41}, {42, 43}}},
			true,
		},
		"face wrong rank": {
			mapCards(44),
			Action{40, [][]card.Card{{44}}},
			true,
		},
		"face captures with wrong rank": {
			mapCards(41, 44),
			Action{40, [][]card.Card{{41}, {44}}},
			true,
		},

		"number empty capture": {
			mapCards(1),
			Action{0, [][]card.Card{{}}},
			true,
		},
		"number captures with empty capture": {
			mapCards(1, 2, 3),
			Action{0, [][]card.Card{{1}, {2}, {}}},
			true,
		},
		"number captured card not on table": {
			mapCards(1),
			Action{0, [][]card.Card{{0}}},
			true,
		},
		"number captures with captured card not on table": {
			mapCards(1, 2),
			Action{0, [][]card.Card{{1}, {2}, {3}}},
			true,
		},
		"number wrong rank": {
			mapCards(4),
			Action{0, [][]card.Card{{4}}},
			true,
		},
		"number captures with wrong rank": {
			mapCards(1, 4),
			Action{0, [][]card.Card{{1}, {4}}},
			true,
		},
		"number wrong sum": {
			mapCards(0, 28),
			Action{36, [][]card.Card{{0, 28}}},
			true,
		},
		"number captures with wrong sum": {
			mapCards(0, 28, 37),
			Action{36, [][]card.Card{{37}, {0, 28}}},
			true,
		},
		"number duplicate capture": {
			mapCards(1),
			Action{0, [][]card.Card{{1}, {1}}},
			true,
		},
		"number captures with duplicate capture": {
			mapCards(1, 2),
			Action{0, [][]card.Card{{1}, {2}, {2}}},
			true,
		},
		"number duplicate sum": {
			mapCards(0, 4, 24),
			Action{36, [][]card.Card{{0, 4, 24}, {0, 4, 24}}},
			true,
		},
		"number duplicate card in sums": {
			mapCards(0, 32, 33),
			Action{36, [][]card.Card{{0, 32}, {0, 33}}},
			true,
		},
		"number single": {
			mapCards(1),
			Action{0, [][]card.Card{{1}}},
			false,
		},
		"number multiple": {
			mapCards(1, 2, 3),
			Action{0, [][]card.Card{{1}, {2}, {3}}},
			false,
		},
		"number sum": {
			mapCards(0, 4, 24),
			Action{36, [][]card.Card{{0, 4, 24}}},
			false,
		},
		"number multiple sum": {
			mapCards(0, 4, 28, 32),
			Action{36, [][]card.Card{{0, 32}, {4, 28}}},
			false,
		},
		"number mixed": {
			mapCards(0, 4, 28, 32, 37),
			Action{36, [][]card.Card{{37}, {0, 32}, {4, 28}}},
			false,
		},
	} {
		err := (&game{table: test.table}).validateAction(test.a)
		if isErr := err != nil; isErr != test.isErr {
			t.Errorf("validateAction(%q): got err=%v, expected %v",
				name, isErr, test.isErr,
			)
		}
	}
}

func TestCapture(t *testing.T) {
	for name, test := range map[string]struct {
		g      *game
		player int
		c      card.Card
		want   *game
	}{
		"first": {
			&game{
				keep:  [][]card.Card{{}, {}},
				table: map[card.Card]bool{7: true, 16: true, 21: true},
			},
			0,
			7,
			&game{
				keep:  [][]card.Card{{7}, {}},
				table: map[card.Card]bool{16: true, 21: true},
			},
		},
		"sweep": {
			&game{
				keep: [][]card.Card{
					{31, 30, 43, 40, 8, 11},
					{7, 5, 0, 2, 50, 48},
				},
				table: map[card.Card]bool{25: true},
			},
			1,
			25,
			&game{
				keep: [][]card.Card{
					{31, 30, 43, 40, 8, 11},
					{7, 5, 0, 2, 50, 48, 25},
				},
				table: map[card.Card]bool{},
			},
		},
	} {
		test.g.capture(test.player, test.c)
		if !reflect.DeepEqual(test.g, test.want) {
			t.Errorf("capture(%q): got %+v, expected %+v",
				name, test.g, test.want,
			)
		}
	}
}

func TestScore(t *testing.T) {
	for _, test := range []struct {
		cards []card.Card
		n     int
	}{
		{[]card.Card{}, 0},
		{[]card.Card{0}, 1},
		{[]card.Card{1}, 1},
		{[]card.Card{2}, 1},
		{[]card.Card{3}, 1},
		{[]card.Card{7}, 1},
		{[]card.Card{37}, 2},
		{[]card.Card{0, 1, 2, 3, 7, 37}, 7},
		{[]card.Card{27, 31, 35, 39, 43, 47}, 0},
		{[]card.Card{27, 31, 35, 39, 43, 47, 51}, 1},
		{[]card.Card{
			14, 16, 17, 18, 20, 21, 22, 24, 25, 26, 28, 29, 30,
			32, 33, 34, 36, 38, 40, 41, 42, 44, 45, 46, 48, 49,
		}, 0},
		{[]card.Card{
			14, 16, 17, 18, 20, 21, 22, 24, 25, 26, 28, 29, 30,
			32, 33, 34, 36, 38, 40, 41, 42, 44, 45, 46, 48, 49,
			50,
		}, 3},
		{[]card.Card{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
			13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
			26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38,
			39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
		}, 11},
	} {
		if n := score(test.cards); n != test.n {
			t.Errorf("score(%v): got %v, expected %v", test.cards, n, test.n)
		}
	}
}
