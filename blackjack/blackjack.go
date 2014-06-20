package blackjack

import (
	. "github.com/dohodges/fifty2"
)

func CardValue(c Card) uint {
	switch c.Rank {
	case Rank2:
		return 2
	case Rank3:
		return 3
	case Rank4:
		return 4
	case Rank5:
		return 5
	case Rank6:
		return 6
	case Rank7:
		return 7
	case Rank8:
		return 8
	case Rank9:
		return 9
	case Rank10, RankJack, RankQueen, RankKing:
		return 10
	case RankAce:
		return 11
	}
	return 0
}

func HandValue(hand []Card) uint {
	aces := 0
	value := uint(0)
	for _, c := range hand {
		value += CardValue(c)
		if c.Rank == RankAce {
			aces++
		}
	}

	for ; value > 21 && aces > 0; aces-- {
		value -= 10
	}

	return value
}
