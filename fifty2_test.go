package fifty2

import (
	"reflect"
	"testing"
)

func TestCombinations(t *testing.T) {
	hand := []Card{
		Card{Four, Spades},
		Card{Five, Hearts},
		Card{Six, Diamonds},
		Card{Seven, Clubs},
	}

	combos := make([][]Card, 0, 6)
	for combo := range Combinations(hand, 2) {
		combos = append(combos, combo)
	}

	expect := [][]Card{
		[]Card{Card{Four, Spades}, Card{Five, Hearts}},
		[]Card{Card{Four, Spades}, Card{Six, Diamonds}},
		[]Card{Card{Four, Spades}, Card{Seven, Clubs}},
		[]Card{Card{Five, Hearts}, Card{Six, Diamonds}},
		[]Card{Card{Five, Hearts}, Card{Seven, Clubs}},
		[]Card{Card{Six, Diamonds}, Card{Seven, Clubs}},
	}

	if !reflect.DeepEqual(combos, expect) {
		t.Errorf("missing combinations\nexpect - %v\nactual - %v", expect, combos)
	}

}
