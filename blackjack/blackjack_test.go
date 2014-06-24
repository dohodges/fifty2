package blackjack

import (
	. "github.com/dohodges/fifty2"
	"testing"
)

func TestHandValues(t *testing.T) {

	assertHandValue(t, 4, []Card{
		Card{Two, Clubs},
		Card{Two, Diamonds},
	})

	assertHandValue(t, 12, []Card{
		Card{Ace, Spades},
		Card{Ace, Hearts},
	})

	assertHandValue(t, 13, []Card{
		Card{Ace, Spades},
		Card{Ace, Hearts},
		Card{Ace, Clubs},
	})

	assertHandValue(t, 14, []Card{
		Card{Ace, Spades},
		Card{Ace, Hearts},
		Card{Ace, Clubs},
		Card{Ace, Diamonds},
	})

	assertHandValue(t, 20, []Card{
		Card{Queen, Spades},
		Card{Ten, Clubs},
	})

	assertHandValue(t, 21, []Card{
		Card{Jack, Spades},
		Card{Ace, Hearts},
	})

	assertHandValue(t, 21, []Card{
		Card{Queen, Hearts},
		Card{Jack, Spades},
		Card{Ace, Hearts},
	})

	assertHandValue(t, 21, []Card{
		Card{Two, Hearts},
		Card{Four, Spades},
		Card{Nine, Hearts},
		Card{Three, Hearts},
		Card{Three, Diamonds},
	})

	assertHandValue(t, 29, []Card{
		Card{Queen, Hearts},
		Card{Jack, Spades},
		Card{Nine, Hearts},
	})

}

func assertHandValue(t *testing.T, expect uint, hand []Card) {
	if value := HandValue(hand); value != expect {
		t.Errorf("hand %v value %v != %v", String(hand), value, expect)
	}
}
