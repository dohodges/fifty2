package blackjack

import (
	. "github.com/dohodges/fifty2"
)

func CardValue(c Card) uint {
	switch c.Rank {
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	case Ace:
		return 11
	}
	return 0
}

func HandValue(hand []Card) uint {
	aces := 0
	value := uint(0)
	for _, c := range hand {
		value += CardValue(c)
		if c.Rank == Ace {
			aces++
		}
	}

	for ; value > 21 && aces > 0; aces-- {
		value -= 10
	}

	return value
}
