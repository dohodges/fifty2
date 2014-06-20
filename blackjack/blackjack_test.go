package blackjack

import (
	. "github.com/dohodges/fifty2"
	"testing"
)

func TestHandValues(t *testing.T) {

	assertHandValue(t, 4, []Card{
		Card{Rank2, SuitClubs},
		Card{Rank2, SuitDiamonds},
	})

	assertHandValue(t, 12, []Card{
		Card{RankAce, SuitSpades},
		Card{RankAce, SuitHearts},
	})

	assertHandValue(t, 13, []Card{
		Card{RankAce, SuitSpades},
		Card{RankAce, SuitHearts},
		Card{RankAce, SuitClubs},
	})

	assertHandValue(t, 14, []Card{
		Card{RankAce, SuitSpades},
		Card{RankAce, SuitHearts},
		Card{RankAce, SuitClubs},
		Card{RankAce, SuitDiamonds},
	})

	assertHandValue(t, 20, []Card{
		Card{RankQueen, SuitSpades},
		Card{Rank10, SuitClubs},
	})

	assertHandValue(t, 21, []Card{
		Card{RankJack, SuitSpades},
		Card{RankAce, SuitHearts},
	})

	assertHandValue(t, 21, []Card{
		Card{RankQueen, SuitHearts},
		Card{RankJack, SuitSpades},
		Card{RankAce, SuitHearts},
	})

	assertHandValue(t, 21, []Card{
		Card{Rank2, SuitHearts},
		Card{Rank4, SuitSpades},
		Card{Rank9, SuitHearts},
		Card{Rank3, SuitHearts},
		Card{Rank3, SuitDiamonds},
	})

	assertHandValue(t, 29, []Card{
		Card{RankQueen, SuitHearts},
		Card{RankJack, SuitSpades},
		Card{Rank9, SuitHearts},
	})

}

func assertHandValue(t *testing.T, expect uint, hand []Card) {
	if value := HandValue(hand); value != expect {
		t.Errorf("hand %v value %v != %v", String(hand), value, expect)
	}
}
