package game

import (
	"reflect"
	"testing"

	"github.com/dkmccandless/cassino/card"
)

var actionTests = map[string]struct {
	g      game
	player int
	a      Action
	isErr  bool
	want   game
}{
	"invalid card": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{1: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{},
		},
		0,
		Action{Card: 0},
		true,
		game{},
	},
	"trail with owned build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{39: true, 50: true},
			},
			piles: map[int]Pile{
				14: Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 1},
			},
		},
		1,
		Action{Card: 50},
		true,
		game{},
	},
	"invalid ID": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{10: Pile{Cards: []card.Card{1}, Value: 1}},
		},
		0,
		Action{Card: 0, Sets: [][]int{{11}}},
		true,
		game{},
	},
	"duplicate ID in sets": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{10: Pile{Cards: []card.Card{1}, Value: 1}},
		},
		0,
		Action{Card: 0, Sets: [][]int{{10}, {10}}},
		true,
		game{},
	},
	"duplicate ID in add": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 8: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{10: Pile{Cards: []card.Card{1}, Value: 1}},
		},
		0,
		Action{Card: 0, Add: []int{10, 10}},
		true,
		game{},
	},
	"duplicate ID between add and sets": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 24: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				10: Pile{Cards: []card.Card{1}, Value: 1},
				11: Pile{Cards: []card.Card{7}, Value: 2},
				12: Pile{Cards: []card.Card{18}, Value: 5},
			},
		},
		0,
		Action{Card: 0, Add: []int{10, 12}, Sets: [][]int{{11, 12}}},
		true,
		game{},
	},
	"face card with add": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{17: true, 18: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{40}, Value: 0},
			},
		},
		0,
		Action{Card: 17, Add: []int{0}},
		true,
		game{},
	},
	"face build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true, 42: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{41}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{0}}, Build: true},
		true,
		game{},
	},
	"face capture invalid set": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{41}, Value: 0},
				1: Pile{Cards: []card.Card{42}, Value: 0},
				2: Pile{Cards: []card.Card{43}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{0, 1, 2}}},
		true,
		game{},
	},
	"face capture wrong rank": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{44}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{0}}},
		true,
		game{},
	},
	"compound add": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 37: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{32, 33}, Value: 9, Compound: true},
			},
		},
		0,
		Action{Card: 0, Add: []int{0}},
		true,
		game{},
	},
	"face card in add": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{5: true, 18: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{40}, Value: 0},
				1: Pile{Cards: []card.Card{9}, Value: 3},
			},
		},
		0,
		Action{Card: 5, Add: []int{0, 1}},
		true,
		game{},
	},
	"face set": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{17: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{40}, Value: 0},
				1: Pile{Cards: []card.Card{18}, Value: 5},
			},
		},
		0,
		Action{Card: 17, Sets: [][]int{{0, 1}}},
		true,
		game{},
	},
	"wrong set value": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true},
				map[card.Card]bool{20: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{28}, Value: 8},
			},
		},
		0,
		Action{Card: 36, Sets: [][]int{{0, 1}}},
		true,
		game{},
	},
	"build with no hand card": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{8: true, 15: true, 42: true, 45: true},
				map[card.Card]bool{20: true, 21: true, 22: true, 23: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{0}, Value: 1},
				1: Pile{Cards: []card.Card{7}, Value: 2},
			},
		},
		0,
		Action{Card: 8, Sets: [][]int{{0, 1}}, Build: true},
		true,
		game{},
	},
	"add build with no hand card": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 33: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				0: Pile{Cards: []card.Card{32}, Value: 9},
			},
		},
		0,
		Action{Card: 0, Add: []int{0}},
		true,
		game{},
	},
	"uncaptured build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{32: true, 36: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{4, 24}, Value: 9, Controller: 1},
				2: Pile{Cards: []card.Card{0}, Value: 1},
				3: Pile{Cards: []card.Card{8}, Value: 3},
				4: Pile{Cards: []card.Card{16}, Value: 5},
			},
			npiles: 4,
		},
		1,
		Action{Card: 32, Sets: [][]int{{2, 3, 4}}},
		true,
		game{},
	},

	"trail": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{39: true, 50: true},
			},
			piles: map[int]Pile{
				14: Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 0},
			},
			npiles: 14,
		},
		1,
		Action{Card: 50},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{39: true},
			},
			piles: map[int]Pile{
				14: Pile{Cards: []card.Card{2, 34, 38}, Value: 10, Controller: 0},
				15: Pile{Cards: []card.Card{50}},
			},
			npiles: 15,
		},
	},
	"face single": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{50}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{41}, Value: 0},
				2: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{1}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{50, 41, 40}, []card.Card{}},
			piles: map[int]Pile{
				2: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"face multiple": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{49, 50}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{41}, Value: 0},
				2: Pile{Cards: []card.Card{42}, Value: 0},
				3: Pile{Cards: []card.Card{43}, Value: 0},
				4: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{1}, {2}, {3}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{49, 50, 41, 42, 43, 40}, []card.Card{}},
			piles: map[int]Pile{
				4: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number single capture": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{1}, Value: 1},
				2: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 0, Sets: [][]int{{1}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{1, 0}, []card.Card{}},
			piles: map[int]Pile{
				2: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number multiple capture": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{1}, Value: 1},
				2: Pile{Cards: []card.Card{2}, Value: 1},
				3: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 0, Sets: [][]int{{1}, {2}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{1, 2, 0}, []card.Card{}},
			piles: map[int]Pile{
				3: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number sum capture": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{24}, Value: 7},
				4: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 36, Sets: [][]int{{1, 2, 3}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{0, 4, 24, 36}, []card.Card{}},
			piles: map[int]Pile{
				4: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number multiple sums capture": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{28}, Value: 8},
				4: Pile{Cards: []card.Card{32}, Value: 9},
				5: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 36, Sets: [][]int{{1, 4}, {2, 3}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{0, 32, 4, 28, 36}, []card.Card{}},
			piles: map[int]Pile{
				5: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number mixed capture": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{28}, Value: 8},
				4: Pile{Cards: []card.Card{32}, Value: 9},
				5: Pile{Cards: []card.Card{37}, Value: 10},
				6: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
		0,
		Action{Card: 36, Sets: [][]int{{5}, {1, 4}, {2, 3}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{37, 0, 32, 4, 28, 36}, []card.Card{}},
			piles: map[int]Pile{
				6: Pile{Cards: []card.Card{51}, Value: 0},
			},
		},
	},
	"number capture build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{32: true, 36: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{4, 24}, Value: 9, Controller: 1},
				2: Pile{Cards: []card.Card{0}, Value: 1},
				3: Pile{Cards: []card.Card{8}, Value: 3},
				4: Pile{Cards: []card.Card{16}, Value: 5},
			},
			npiles:      4,
			lastCapture: 0,
		},
		1,
		Action{Card: 32, Sets: [][]int{{1}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{36: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{4, 24, 32}},
			piles: map[int]Pile{
				2: Pile{Cards: []card.Card{0}, Value: 1},
				3: Pile{Cards: []card.Card{8}, Value: 3},
				4: Pile{Cards: []card.Card{16}, Value: 5},
			},
			npiles:      4,
			lastCapture: 1,
		},
	},
	"number single build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 2: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{1}, Value: 1},
			},
			npiles: 1,
		},
		0,
		Action{Card: 0, Sets: [][]int{{1}}, Build: true},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{2: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				2: Pile{Cards: []card.Card{1, 0}, Value: 1, Compound: true, Controller: 0},
			},
			npiles: 2,
		},
	},
	"number multiple build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 2: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{1}, Value: 1},
				2: Pile{Cards: []card.Card{2}, Value: 1},
			},
			npiles: 2,
		},
		0,
		Action{Card: 0, Sets: [][]int{{1}, {2}}, Build: true},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{2: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				3: Pile{Cards: []card.Card{1, 2, 0}, Value: 1, Compound: true, Controller: 0},
			},
			npiles: 3,
		},
	},
	"number sum build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{36: true, 39: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{24}, Value: 7},
			},
			npiles: 3,
		},
		1,
		Action{Card: 36, Sets: [][]int{{1, 2, 3}}, Build: true},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{39: true},
			},
			piles: map[int]Pile{
				4: Pile{Cards: []card.Card{0, 4, 24, 36}, Value: 10, Compound: true, Controller: 1},
			},
			npiles: 4,
		},
	},
	"number multiple sums build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true, 39: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{28}, Value: 8},
				4: Pile{Cards: []card.Card{32}, Value: 9},
			},
			npiles: 4,
		},
		0,
		Action{Card: 36, Sets: [][]int{{1, 4}, {2, 3}}, Build: true},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{39: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				5: Pile{Cards: []card.Card{0, 32, 4, 28, 36}, Value: 10, Compound: true, Controller: 0},
			},
			npiles: 5,
		},
	},
	"number mixed build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{36: true, 39: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{0}, Value: 1},
				2: Pile{Cards: []card.Card{4}, Value: 2},
				3: Pile{Cards: []card.Card{28}, Value: 8},
				4: Pile{Cards: []card.Card{32}, Value: 9},
				5: Pile{Cards: []card.Card{37}, Value: 10},
			},
			npiles: 5,
		},
		1,
		Action{Card: 36, Sets: [][]int{{5}, {1, 4}, {2, 3}}, Build: true},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{39: true},
			},
			piles: map[int]Pile{
				6: Pile{Cards: []card.Card{37, 0, 32, 4, 28, 36}, Value: 10, Compound: true, Controller: 1},
			},
			npiles: 6,
		},
	},
	"number add build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{0: true, 36: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{32}, Value: 9},
			},
			npiles: 1,
		},
		0,
		Action{Card: 0, Add: []int{1}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{36: true},
				map[card.Card]bool{20: true, 21: true},
			},
			piles: map[int]Pile{
				2: Pile{Cards: []card.Card{32, 0}, Value: 10, Compound: false, Controller: 0},
			},
			npiles: 2,
		},
	},
	"number add sets build": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{0: true, 36: true},
			},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{4}, Value: 2},
				2: Pile{Cards: []card.Card{28}, Value: 8},
				3: Pile{Cards: []card.Card{32}, Value: 9},
			},
			npiles: 3,
		},
		1,
		Action{Card: 0, Add: []int{3}, Sets: [][]int{{1, 2}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{36: true},
			},
			piles: map[int]Pile{
				4: Pile{Cards: []card.Card{4, 28, 32, 0}, Value: 10, Compound: true, Controller: 1},
			},
			npiles: 4,
		},
	},
	"uncaptured build with hand card": {
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{32: true, 33: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{4, 24}, Value: 9, Controller: 1},
				2: Pile{Cards: []card.Card{0}, Value: 1},
				3: Pile{Cards: []card.Card{8}, Value: 3},
				4: Pile{Cards: []card.Card{16}, Value: 5},
			},
			npiles:      4,
			lastCapture: 0,
		},
		1,
		Action{Card: 32, Sets: [][]int{{2, 3, 4}}},
		false,
		game{
			hand: []map[card.Card]bool{
				map[card.Card]bool{20: true},
				map[card.Card]bool{33: true},
			},
			keep: [][]card.Card{[]card.Card{}, []card.Card{0, 8, 16, 32}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{4, 24}, Value: 9, Controller: 1},
			},
			npiles:      4,
			lastCapture: 1,
		},
	},
	"sweep": {
		game{
			score: []int{0, 0},
			hand: []map[card.Card]bool{
				map[card.Card]bool{40: true},
				map[card.Card]bool{20: true},
			},
			keep: [][]card.Card{[]card.Card{50}, []card.Card{}},
			piles: map[int]Pile{
				1: Pile{Cards: []card.Card{41}, Value: 0},
			},
		},
		0,
		Action{Card: 40, Sets: [][]int{{1}}},
		false,
		game{
			score: []int{1, 0},
			hand: []map[card.Card]bool{
				map[card.Card]bool{},
				map[card.Card]bool{20: true},
			},
			keep:  [][]card.Card{[]card.Card{50, 41, 40}, []card.Card{}},
			piles: map[int]Pile{},
		},
	},
}

func TestValidateAction(t *testing.T) {
	for name, test := range actionTests {
		err := test.g.validateAction(test.player, test.a)
		if isErr := err != nil; isErr != test.isErr {
			switch {
			case isErr:
				t.Errorf("validateAction(%q): got %v, expected nil", name, err)
			case test.isErr:
				t.Errorf("validateAction(%q): got nil, expected error", name)
			}
		}
	}
}

func TestDo(t *testing.T) {
	for name, test := range actionTests {
		if test.isErr {
			continue
		}
		test.g.do(test.player, test.a)
		if !reflect.DeepEqual(test.g, test.want) {
			t.Errorf("do(%q): got %+v, expected %+v", name, test.g, test.want)
		}
	}
}

func TestIsBuild(t *testing.T) {
	for _, test := range []struct {
		a    Action
		want bool
	}{
		{Action{Card: 20}, false},
		{Action{Card: 20, Add: []int{7}}, true},
		{Action{Card: 20, Sets: [][]int{{9}}}, false},
		{Action{Card: 20, Sets: [][]int{{9}}, Build: true}, true},
		{Action{Card: 20, Add: []int{7}, Sets: [][]int{{9}}, Build: true}, true},
	} {
		if isBuild := test.a.isBuild(); isBuild != test.want {
			t.Errorf("isBuild(%+v): got %v, expected %v",
				test.a, isBuild, test.want,
			)
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

func TestHaveSameRank(t *testing.T) {
	for _, test := range []struct {
		hand    map[card.Card]bool
		rank    int
		exclude card.Card
		want    bool
	}{
		{map[card.Card]bool{}, 9, 32, false},
		{map[card.Card]bool{32: true}, 9, 32, false},
		{map[card.Card]bool{32: true, 33: true, 34: true}, 9, 32, true},
		{map[card.Card]bool{32: true, 33: true, 37: true}, 9, 32, true},
		{map[card.Card]bool{32: true, 36: true, 37: true}, 9, 32, false},
	} {
		got := haveSameRank(test.hand, test.rank, test.exclude)
		if got != test.want {
			t.Errorf("haveSameRank(%v, %v, %v): got %v, expected %v",
				test.hand, test.rank, test.exclude, got, test.want,
			)
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
