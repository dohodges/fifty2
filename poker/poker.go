package poker

import (
	. "github.com/dohodges/fifty2"
	"github.com/golang/groupcache/lru"
)

var cache *lru.Cache

func init() {
	cache = lru.New(2598960)
}

type HandRank uint8

const (
	NoHand HandRank = iota
	HighCard
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
	case NoHand:
		return "No Hand"
	}
	return ""
}

type CardStrength uint8

const (
	AceLow  CardStrength = CardStrength(Ace)
	AceHigh CardStrength = CardStrength(King + 1)
)

func (cs CardStrength) Rank() Rank { return Rank(cs % 13) }

func (cs CardStrength) Mask() uint16 { return uint16(1) << cs }

func MaxCardStrength(strengths []CardStrength) CardStrength {
	if len(strengths) == 0 {
		panic("fifty2/poker: cannot find max strength from an empty slice")
	}
	max := strengths[0]
	for i := 1; i < len(strengths); i++ {
		if strengths[i] > max {
			max = strengths[i]
		}
	}
	return max
}

type HandStrength uint32

func MakeHandStrength(rank HandRank, strength1, strength2 CardStrength, kickers uint16) HandStrength {
	return HandStrength(uint32(rank)<<24 | uint32(strength1)<<20 | uint32(strength2)<<16 | uint32(kickers))
}

func (hs HandStrength) Rank() HandRank {
	return HandRank(hs >> 24)
}

func MaxHandStrength(strengths []HandStrength) HandStrength {
	if len(strengths) == 0 {
		panic("fifty2/poker: cannot find max strength from an empty slice")
	}
	max := strengths[0]
	for i := 1; i < len(strengths); i++ {
		if strengths[i] > max {
			max = strengths[i]
		}
	}
	return max
}

func MinHandStrength(strengths []HandStrength) HandStrength {
	if len(strengths) == 0 {
		panic("fifty2/poker: cannot find min strength from an empty slice")
	}
	min := strengths[0]
	for i := 1; i < len(strengths); i++ {
		if strengths[i] < min {
			min = strengths[i]
		}
	}
	return min
}

func GetHandStrength(hand []Card) HandStrength {
	if strength, hit := cache.Get(Mask(hand)); hit {
		return strength.(HandStrength)
	}

	strength := calculateHandStrength(hand)
	cache.Add(Mask(hand), strength)

	return strength
}

func calculateHandStrength(hand []Card) HandStrength {

	var (
		bitSet     uint16
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
		return MakeHandStrength(StraightFlush, strength, 0, 0)
	}

	// quads
	for strength := AceHigh; strength > AceLow; strength-- {
		if rankCount[strength.Rank()] >= 4 {
			kickers := getKickers(bitSet&^strength.Rank().Mask(), 1)
			return MakeHandStrength(Quads, strength, 0, kickers)
		}
	}

	// full house
	for hiStrength := AceHigh; hiStrength > AceLow; hiStrength-- {
		if rankCount[hiStrength.Rank()] >= 3 {
			for loStrength := AceHigh; loStrength > AceLow; loStrength-- {
				if loStrength != hiStrength && rankCount[loStrength.Rank()] >= 2 {
					return MakeHandStrength(FullHouse, hiStrength, loStrength, 0)
				}
			}
		}
	}

	// flush
	flushBitSet := uint16(0)
	for suit, count := range suitCount {
		if count >= 5 && suitBitSet[suit] > flushBitSet {
			flushBitSet = suitBitSet[suit]
		}
	}
	if flushBitSet > 0 {
		return MakeHandStrength(Flush, 0, 0, getKickers(flushBitSet, 5))
	}

	// straight
	if strength, found := findStraight(bitSet); found {
		return MakeHandStrength(Straight, strength, 0, 0)
	}

	// trips
	for strength := AceHigh; strength > AceLow; strength-- {
		if rankCount[strength.Rank()] >= 3 {
			kickers := getKickers(bitSet&^strength.Rank().Mask(), 2)
			return MakeHandStrength(Trips, strength, 0, kickers)
		}
	}

	// two pair / pair
	for hiStrength := AceHigh; hiStrength > AceLow; hiStrength-- {
		if rankCount[hiStrength.Rank()] >= 2 {
			for loStrength := hiStrength - 1; loStrength > AceLow; loStrength-- {
				if rankCount[loStrength.Rank()] >= 2 {
					kickers := getKickers(bitSet&^(hiStrength.Rank().Mask()|loStrength.Rank().Mask()), 1)
					return MakeHandStrength(TwoPair, hiStrength, loStrength, kickers)
				}
			}
			kickers := getKickers(bitSet&^hiStrength.Rank().Mask(), 3)
			return MakeHandStrength(Pair, hiStrength, 0, kickers)
		}
	}

	// high card
	return MakeHandStrength(HighCard, 0, 0, getKickers(bitSet, 5))
}

func GetLowHandStrength(hand []Card, eightOrBetter bool) HandStrength {

	var (
		bitSet    uint16
		rankCount [13]uint8
	)

	for _, card := range hand {
		rankCount[card.Rank]++
		bitSet |= card.Rank.Mask()
	}

	// high card
	kickers, found := getLowKickers(bitSet, 5)
	if found == 5 || found == len(hand) {
		if !eightOrBetter || kickers < CardStrength(9).Mask() {
			return MakeHandStrength(HighCard, 0, 0, kickers)
		}
	}

	if eightOrBetter {
		return MakeHandStrength(NoHand, 0, 0, 0)
	}

	// pair / two pair
	for loStrength := AceLow; loStrength < AceHigh; loStrength++ {
		if rankCount[loStrength.Rank()] >= 2 {
			kickers, found = getLowKickers(bitSet&^loStrength.Rank().Mask(), 3)
			if found == 3 || found == (len(hand)-2) {
				return MakeHandStrength(Pair, loStrength, 0, kickers)
			}
			for hiStrength := loStrength + 1; hiStrength < AceHigh; hiStrength++ {
				if rankCount[hiStrength.Rank()] >= 2 {
					kickers, found = getLowKickers(bitSet&^(loStrength.Rank().Mask()|hiStrength.Rank().Mask()), 1)
					if found == 1 || len(hand) == 4 {
						return MakeHandStrength(TwoPair, hiStrength, loStrength, kickers)
					}
				}
			}
		}
	}

	// trips / full house
	for hiStrength := AceLow; hiStrength < AceHigh; hiStrength++ {
		if rankCount[hiStrength.Rank()] >= 3 {
			kickers, found = getLowKickers(bitSet&^hiStrength.Rank().Mask(), 2)
			if found == 2 || found == (len(hand)-3) {
				return MakeHandStrength(Trips, hiStrength, 0, kickers)
			}
			for loStrength := AceLow; loStrength < AceHigh; loStrength++ {
				if loStrength != hiStrength && rankCount[loStrength.Rank()] >= 2 {
					return MakeHandStrength(FullHouse, hiStrength, loStrength, 0)
				}
			}
		}
	}

	// quads
	for strength := AceLow; strength < AceHigh; strength++ {
		if rankCount[strength.Rank()] >= 4 {
			kickers, _ = getLowKickers(bitSet&^strength.Rank().Mask(), 1)
			return MakeHandStrength(Quads, strength, 0, kickers)
		}
	}

	panic("fifty2/poker: impossible low hand")
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

func getKickers(bitSet uint16, max int) uint16 {
	kickers := uint16(0)
	found := 0
	for strength := AceHigh; strength > AceLow; strength-- {
		if (bitSet & strength.Rank().Mask()) != 0 {
			kickers |= strength.Mask()
			found++
			if found == max {
				return kickers
			}
		}
	}
	return kickers
}

func getLowKickers(bitSet uint16, max int) (uint16, int) {
	kickers := uint16(0)
	found := 0
	for strength := AceLow; strength < AceHigh; strength++ {
		if (bitSet & strength.Rank().Mask()) != 0 {
			kickers |= strength.Mask()
			found++
			if found == max {
				return kickers, found
			}
		}
	}
	return kickers, found
}
