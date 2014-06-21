package fifty2

import (
	"reflect"
	"testing"
)

func TestCombinations(t *testing.T) {
	hand := []Card{
		Card{Rank4, SuitSpades},
		Card{Rank5, SuitHearts},
		Card{Rank6, SuitDiamonds},
		Card{Rank7, SuitClubs},
	}

	combos := make([][]Card, 0, 6)
	for combo := range Combinations(hand, 2) {
		combos = append(combos, combo)
	}

	expect := [][]Card{
		[]Card{Card{Rank4, SuitSpades}, Card{Rank5, SuitHearts}},
		[]Card{Card{Rank4, SuitSpades}, Card{Rank6, SuitDiamonds}},
		[]Card{Card{Rank4, SuitSpades}, Card{Rank7, SuitClubs}},
		[]Card{Card{Rank5, SuitHearts}, Card{Rank6, SuitDiamonds}},
		[]Card{Card{Rank5, SuitHearts}, Card{Rank7, SuitClubs}},
		[]Card{Card{Rank6, SuitDiamonds}, Card{Rank7, SuitClubs}},
	}

	if !reflect.DeepEqual(combos, expect) {
		t.Errorf("missing combinations\nexpect - %v\nactual - %v", expect, combos)
	}

}
