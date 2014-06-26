package poker

import (
	. "github.com/dohodges/fifty2"
)

type HandRank uint16

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	Quads
	StraightFlush
)

func HandRanks() []HandRank {
	return []HandRank{StraightFlush, Quads, FullHouse, Flush, Straight, Trips, TwoPair, Pair, HighCard}
}

func (hr HandRank) String() string {
	switch hr {
	case StraightFlush:
		return "Straight Flush"
	case Quads:
		return "Quads"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case Trips:
		return "Trips"
	case TwoPair:
		return "Two Pair"
	case Pair:
		return "Pair"
	case HighCard:
		return "High Card"
	}
	return ""
}

type CardStrength Rank

const (
	AceLow  CardStrength = CardStrength(Ace)
	AceHigh CardStrength = CardStrength(King + 1)
)

func (cs CardStrength) Rank() Rank {
	return Rank(cs % 13)
}

func MaxCardStrength(strengths []CardStrength) CardStrength {
	if len(strengths) == 0 {
		panic("fifty2: cannot find max strength from an empty slice")
	}
	max := strengths[0]
	for i := 1; i < len(strengths); i++ {
		if strengths[i] > max {
			max = strengths[i]
		}
	}
	return max
}

type HandStrength struct {
	Rank     HandRank
	Strength []CardStrength
}

func GetHandStrength(hand []Card) HandStrength {

	var (
		bitSet uint16
		suitBitSet [4]uint16
		rankCount  [13]uint8
		suitCount  [4]uint8
	)

	for _, card := range hand {
		rankCount[card.Rank]++
		suitCount[card.Suit]++
		suitBitSet[card.Suit] |= card.Rank.Mask()
		bitSet |= card.Rank.Mask()
	}

	// straight flush
	straights := make([]CardStrength, 0, 4)
	for _, bitSet := range suitBitSet {
		if strength, found := findStraight(bitSet); found {
			straights = append(straights, strength)
		}
	}
	if len(straights) > 0 {
		strength := MaxCardStrength(straights)
		return HandStrength{StraightFlush, []CardStrength{strength}}
	}

	// quads
	for strength := AceHigh; strength > AceLow; strength-- {
		if rankCount[strength.Rank()] >= 4 {
			kickers := getKickers(bitSet &^ strength.Rank().Mask(), 1)
			return HandStrength{Quads, append([]CardStrength{strength}, kickers...)}
		}
	}

	// full house
	for hiStrength := AceHigh; hiStrength > AceLow; hiStrength-- {
		if rankCount[hiStrength.Rank()] >= 3 {
			for loStrength := AceHigh; loStrength > AceLow; loStrength-- {
				if loStrength != hiStrength && rankCount[loStrength.Rank()] >= 2 {
					return HandStrength{FullHouse, []CardStrength{hiStrength, loStrength}}
				}
			}
		}
	}

	// flush
	flushBitSet := uint16(0)
	for suit, count := range suitCount {
		if count >= 5  && suitBitSet[suit] > flushBitSet {
			flushBitSet = suitBitSet[suit]
		}
	}
	if flushBitSet > 0 {
		return HandStrength{Flush, getKickers(flushBitSet, 5)}
	}

	// straight
	if strength, found := findStraight(bitSet); found {
		return HandStrength{Straight, []CardStrength{strength}}
	}

	// trips
	for strength := AceHigh; strength > AceLow; strength-- {
		if rankCount[strength.Rank()] >= 3 {
			kickers := getKickers(bitSet &^ strength.Rank().Mask(), 2)
			return HandStrength{Trips, append([]CardStrength{strength}, kickers...)}
		}
	}

	// two pair / pair
	for hiStrength := AceHigh; hiStrength > AceLow; hiStrength-- {
		if rankCount[hiStrength.Rank()] >= 2 {
			for loStrength := AceHigh; loStrength > AceLow; loStrength-- {
				if loStrength != hiStrength && rankCount[loStrength.Rank()] >= 2 {
					kickers := getKickers(bitSet &^ (hiStrength.Rank().Mask() | loStrength.Rank().Mask()), 1)
					return HandStrength{TwoPair, append([]CardStrength{hiStrength, loStrength}, kickers...)}
				}
			}
			kickers := getKickers(bitSet &^ hiStrength.Rank().Mask(), 3)
			return HandStrength{Pair, append([]CardStrength{hiStrength}, kickers...)}
		}
	}

	// high card
	return HandStrength{HighCard, getKickers(bitSet, 5)}
}

func findStraight(bitSet uint16) (CardStrength, bool) {
	// ace high straight - 0001 1110 0000 0001
	mask := uint16(0x1E01)
	if (bitSet & mask) == mask {
		return AceHigh, true
	}

	for r := King; r >= Five; r-- {
		mask = uint16(0x001F) << (r - Five)
		if (bitSet & mask) == mask {
			return CardStrength(r), true
		}
	}

	return 0, false
}

func getKickers(bitSet uint16, max int) []CardStrength {
	kickers := make([]CardStrength, 0, max)
	for strength := AceHigh; strength > AceLow; strength-- {
		if (bitSet & strength.Rank().Mask()) != 0 {
			kickers = append(kickers, strength)
			if len(kickers) == max {
				return kickers
			}
		}
	}
	return kickers
}
