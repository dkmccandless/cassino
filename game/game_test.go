package game

import (
	"reflect"
	"testing"

	"github.com/dkmccandless/cassino/card"
)

func TestValidateAction(t *testing.T) {
	for name, test := range map[string]struct {
		piles  map[int]Pile
		player int
		hand   map[card.Card]bool
		a      Action
		isErr  bool
	}{
		"invalid card": {
			map[int]Pile{},
			0,
			map[card.Card]bool{1: true},
			Action{Card: 0},
			true,
		},
		"trail with controlled build": {
			map[int]Pile{
				14: Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 1},
			},
			1,
			map[card.Card]bool{39: true, 50: true},
			Action{Card: 50},
			true,
		},
		"invalid ID": {
			map[int]Pile{10: Pile{Cards: []card.Card{1}, Value: 1}},
			0,
			map[card.Card]bool{0: true},
			Action{Card: 0, Sets: [][]int{{11}}},
			true,
		},
		"duplicate ID": {
			map[int]Pile{10: Pile{Cards: []card.Card{1}, Value: 1}},
			0,
			map[card.Card]bool{0: true},
			Action{Card: 0, Sets: [][]int{{10}, {10}}},
			true,
		},
		"face capture invalid set": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{41}, Value: 0},
				1: Pile{Cards: []card.Card{42}, Value: 0},
				2: Pile{Cards: []card.Card{43}, Value: 0},
			},
			0,
			map[card.Card]bool{40: true},
			Action{Card: 40, Sets: [][]int{{0, 1, 2}}},
			true,
		},
		"face capture wrong rank": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{44}, Value: 0},
			},
			0,
			map[card.Card]bool{40: true},
			Action{Card: 40, Sets: [][]int{{0}}},
			true,
		},
		"wrong set value": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{28}, Value: 8},
			},
			0,
			map[card.Card]bool{36: true},
			Action{Card: 36, Sets: [][]int{{0, 1}}},
			true,
		},
		"build with no hand card": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{7}, Value: 2},
			},
			0,
			map[card.Card]bool{8: true, 15: true, 42: true, 45: true},
			Action{Card: 8, Sets: [][]int{{0, 1}}, Build: true},
			true,
		},

		"trail": {
			map[int]Pile{
				14: Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 0},
			},
			1,
			map[card.Card]bool{39: true, 50: true},
			Action{Card: 50},
			false,
		},
		"face single": {
			map[int]Pile{0: Pile{Cards: []card.Card{41}, Value: 0}},
			0,
			map[card.Card]bool{40: true},
			Action{Card: 40, Sets: [][]int{{0}}},
			false,
		},
		"face multiple": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{41}, Value: 0},
				1: Pile{Cards: []card.Card{42}, Value: 0},
				2: Pile{Cards: []card.Card{43}, Value: 0},
			},
			0,
			map[card.Card]bool{40: true},
			Action{Card: 40, Sets: [][]int{{0}, {1}, {2}}},
			false,
		},
		"number single capture": {
			map[int]Pile{0: Pile{Cards: []card.Card{1}, Value: 1}},
			0,
			map[card.Card]bool{0: true},
			Action{Card: 0, Sets: [][]int{{0}}},
			false,
		},
		"number multiple capture": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{1}, Value: 1},
				1: Pile{Cards: []card.Card{2}, Value: 1},
			},
			0,
			map[card.Card]bool{0: true},
			Action{Card: 0, Sets: [][]int{{0}, {1}}},
			false,
		},
		"number sum capture": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{24}, Value: 7},
			},
			0,
			map[card.Card]bool{36: true},
			Action{Card: 36, Sets: [][]int{{0, 1, 2}}},
			false,
		},
		"number multiple sums capture": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{28}, Value: 8},
				3: Pile{Cards: []card.Card{32}, Value: 9},
			},
			0,
			map[card.Card]bool{36: true},
			Action{Card: 36, Sets: [][]int{{0, 3}, {1, 2}}},
			false,
		},
		"number mixed capture": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{28}, Value: 8},
				3: Pile{Cards: []card.Card{32}, Value: 9},
				4: Pile{Cards: []card.Card{37}, Value: 10},
			},
			0,
			map[card.Card]bool{36: true},
			Action{Card: 36, Sets: [][]int{{4}, {0, 3}, {1, 2}}},
			false,
		},
		"number single build": {
			map[int]Pile{0: Pile{Cards: []card.Card{1}, Value: 1}},
			0,
			map[card.Card]bool{0: true, 2: true},
			Action{Card: 0, Sets: [][]int{{0}}, Build: true},
			false,
		},
		"number multiple build": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{1}, Value: 1},
				1: Pile{Cards: []card.Card{2}, Value: 1},
			},
			0,
			map[card.Card]bool{0: true, 2: true},
			Action{Card: 0, Sets: [][]int{{0}, {1}}, Build: true},
			false,
		},
		"number sum build": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{24}, Value: 7},
			},
			0,
			map[card.Card]bool{36: true, 39: true},
			Action{Card: 36, Sets: [][]int{{0, 1, 2}}, Build: true},
			false,
		},
		"number multiple sums build": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{28}, Value: 8},
				3: Pile{Cards: []card.Card{32}, Value: 9},
			},
			0,
			map[card.Card]bool{36: true, 39: true},
			Action{Card: 36, Sets: [][]int{{0, 3}, {1, 2}}, Build: true},
			false,
		},
		"number mixed build": {
			map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{28}, Value: 8},
				3: Pile{Cards: []card.Card{32}, Value: 9},
				4: Pile{Cards: []card.Card{37}, Value: 10},
			},
			0,
			map[card.Card]bool{36: true, 39: true},
			Action{Card: 36, Sets: [][]int{{4}, {0, 3}, {1, 2}}, Build: true},
			false,
		},
	} {
		err := (&game{piles: test.piles}).validateAction(test.player, test.hand, test.a)
		if isErr := err != nil; isErr != test.isErr {
			switch {
			case isErr:
				t.Errorf("validateAction(%q): got %v, expected nil",
					name, err,
				)
			case test.isErr:
				t.Errorf("validateAction(%q): got nil, expected error", name)
			}
		}
	}
}

func TestAddCardPile(t *testing.T) {
	for name, test := range map[string]struct {
		g    *game
		c    card.Card
		want *game
	}{
		"empty": {
			&game{
				piles:  map[int]Pile{},
				npiles: 9,
			},
			11,
			&game{
				piles: map[int]Pile{
					10: Pile{Cards: []card.Card{11}, Value: 3},
				},
				npiles: 10,
			},
		},
		"non-empty": {
			&game{
				piles: map[int]Pile{
					6:  Pile{Cards: []card.Card{50}, Value: 0},
					9:  Pile{Cards: []card.Card{16}, Value: 5},
					11: Pile{Cards: []card.Card{33}, Value: 9},
					12: Pile{Cards: []card.Card{22}, Value: 6},
				},
				npiles: 15,
			},
			45,
			&game{
				piles: map[int]Pile{
					6:  Pile{Cards: []card.Card{50}, Value: 0},
					9:  Pile{Cards: []card.Card{16}, Value: 5},
					11: Pile{Cards: []card.Card{33}, Value: 9},
					12: Pile{Cards: []card.Card{22}, Value: 6},
					16: Pile{Cards: []card.Card{45}, Value: 0},
				},
				npiles: 16,
			},
		},
	} {
		if test.g.addCardPile(test.c); !reflect.DeepEqual(test.g, test.want) {
			t.Errorf("addCardPile(%q): got %+v, expected %+v",
				name, test.g, test.want,
			)
		}
	}
}

func TestAddPile(t *testing.T) {
	for name, test := range map[string]struct {
		g    *game
		p    Pile
		want *game
	}{
		"single": {
			&game{
				piles:  map[int]Pile{},
				npiles: 9,
			},
			Pile{Cards: []card.Card{11}, Value: 3},
			&game{
				piles: map[int]Pile{
					10: Pile{Cards: []card.Card{11}, Value: 3},
				},
				npiles: 10,
			},
		},
		"empty": {
			&game{
				piles:  map[int]Pile{},
				npiles: 9,
			},
			Pile{Cards: []card.Card{0, 32, 36}, Value: 10},
			&game{
				piles: map[int]Pile{
					10: Pile{Cards: []card.Card{0, 32, 36}, Value: 10},
				},
				npiles: 10,
			},
		},
		"non-empty": {
			&game{
				piles: map[int]Pile{
					6:  Pile{Cards: []card.Card{50}, Value: 0},
					9:  Pile{Cards: []card.Card{16}, Value: 5},
					11: Pile{Cards: []card.Card{33}, Value: 9},
					12: Pile{Cards: []card.Card{22}, Value: 6},
				},
				npiles: 15,
			},
			Pile{Cards: []card.Card{19, 15, 3, 31, 32}, Value: 9},
			&game{
				piles: map[int]Pile{
					6:  Pile{Cards: []card.Card{50}, Value: 0},
					9:  Pile{Cards: []card.Card{16}, Value: 5},
					11: Pile{Cards: []card.Card{33}, Value: 9},
					12: Pile{Cards: []card.Card{22}, Value: 6},
					16: Pile{Cards: []card.Card{19, 15, 3, 31, 32}, Value: 9},
				},
				npiles: 16,
			},
		},
	} {
		if test.g.addPile(test.p); !reflect.DeepEqual(test.g, test.want) {
			t.Errorf("addPile(%q): got %+v, expected %+v",
				name, test.g, test.want,
			)
		}
	}
}

func TestCapture(t *testing.T) {
	for name, test := range map[string]struct {
		g      *game
		player int
		id     int
		want   *game
	}{
		"first": {
			&game{
				keep: [][]card.Card{{}, {}},
				piles: map[int]Pile{
					0: {Cards: []card.Card{7}, Value: 2},
					1: {Cards: []card.Card{16}, Value: 5},
					2: {Cards: []card.Card{21}, Value: 6},
				},
			},
			0,
			1,
			&game{
				keep: [][]card.Card{{16}, {}},
				piles: map[int]Pile{
					0: {Cards: []card.Card{7}, Value: 2},
					2: {Cards: []card.Card{21}, Value: 6},
				},
			},
		},
		"sweep": {
			&game{
				keep: [][]card.Card{
					{31, 30, 43, 40, 8, 11},
					{7, 5, 0, 2, 50, 48},
				},
				piles: map[int]Pile{15: {Cards: []card.Card{25}, Value: 7}},
			},
			1,
			15,
			&game{
				keep: [][]card.Card{
					{31, 30, 43, 40, 8, 11},
					{7, 5, 0, 2, 50, 48, 25},
				},
				piles: map[int]Pile{},
			},
		},
	} {
		test.g.capture(test.player, test.id)
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

func TestCopyPile(t *testing.T) {
	for _, p := range []Pile{
		Pile{Cards: []card.Card{51}},
		Pile{Cards: []card.Card{0}, Value: 1},
		Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 1},
	} {
		c := copyPile(p)
		if &c.Cards == &p.Cards {
			t.Errorf("copyPile(%+v): copy shares memory", p)
		}
		if !reflect.DeepEqual(c, p) {
			t.Errorf("copyPile(%+v): copy is not deeply equal", p)
		}
	}
}
