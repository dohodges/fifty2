package fifty2

import (
	"reflect"
	"strings"
	"testing"
)

func TestCardReader(t *testing.T) {
	card, _ := NewCardReader(strings.NewReader("7♠")).Read()
	if !reflect.DeepEqual(card, Card{Seven, Spades}) {
		t.Errorf("incorrect card read - %s", card)
	}

	hand, _ := NewCardReader(strings.NewReader("3C4D")).ReadAll()
	if !reflect.DeepEqual(hand, []Card{Card{Three, Clubs}, Card{Four, Diamonds}}) {
		t.Errorf("incorrect hand read - %s", hand)
	}
}

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
