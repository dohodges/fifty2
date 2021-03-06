package fifty2

import (
	"reflect"
	"strings"
	"testing"
)

func BenchmarkCombinations(b *testing.B) {
	deck := NewDeck()
	for i := 0; i < b.N; i++ {
		for itr := Combinations(deck, 7); itr.HasNext(); itr.Next() {
		}
	}
}

func BenchmarkMultipleCombinations(b *testing.B) {
	deck := NewDeck()
	for i := 0; i < b.N; i++ {
		for itr := MultipleCombinations(deck, []int{3, 2}); itr.HasNext(); itr.Next() {
		}
	}
}

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
	for itr := Combinations(hand, 2); itr.HasNext(); {
		combos = append(combos, itr.Next())
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

func TestMultipleCombinations(t *testing.T) {
	deck := []Card{
		Card{Four, Spades},
		Card{Five, Hearts},
		Card{Six, Diamonds},
		Card{Seven, Clubs},
	}

	comboSets := make([][][]Card, 0, 12)
	for itr := MultipleCombinations(deck, []int{2, 1}); itr.HasNext(); {
		comboSets = append(comboSets, itr.Next())
	}

	expect := [][][]Card{
		[][]Card{
			[]Card{Card{Four, Spades}, Card{Five, Hearts}},
			[]Card{Card{Six, Diamonds}},
		},
		[][]Card{
			[]Card{Card{Four, Spades}, Card{Five, Hearts}},
			[]Card{Card{Seven, Clubs}},
		},

		[][]Card{
			[]Card{Card{Four, Spades}, Card{Six, Diamonds}},
			[]Card{Card{Five, Hearts}},
		},
		[][]Card{
			[]Card{Card{Four, Spades}, Card{Six, Diamonds}},
			[]Card{Card{Seven, Clubs}},
		},

		[][]Card{
			[]Card{Card{Four, Spades}, Card{Seven, Clubs}},
			[]Card{Card{Five, Hearts}},
		},
		[][]Card{
			[]Card{Card{Four, Spades}, Card{Seven, Clubs}},
			[]Card{Card{Six, Diamonds}},
		},

		[][]Card{
			[]Card{Card{Five, Hearts}, Card{Six, Diamonds}},
			[]Card{Card{Four, Spades}},
		},
		[][]Card{
			[]Card{Card{Five, Hearts}, Card{Six, Diamonds}},
			[]Card{Card{Seven, Clubs}},
		},

		[][]Card{
			[]Card{Card{Five, Hearts}, Card{Seven, Clubs}},
			[]Card{Card{Four, Spades}},
		},
		[][]Card{
			[]Card{Card{Five, Hearts}, Card{Seven, Clubs}},
			[]Card{Card{Six, Diamonds}},
		},

		[][]Card{
			[]Card{Card{Six, Diamonds}, Card{Seven, Clubs}},
			[]Card{Card{Four, Spades}},
		},
		[][]Card{
			[]Card{Card{Six, Diamonds}, Card{Seven, Clubs}},
			[]Card{Card{Five, Hearts}},
		},
	}

	if !reflect.DeepEqual(comboSets, expect) {
		t.Errorf("missing combinations\nexpect - %v\nactual - %v", expect, comboSets)
	}

}
