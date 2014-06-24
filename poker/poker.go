package poker

import (
	. "github.com/dohodges/fifty2"
	"sort"
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

const (
	AceLow Rank = Ace
	AceHigh Rank = King + 1
)


type HandStrength struct {
	HandRank HandRank
	Strength RankSlice
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
	straightRanks := make([]Rank, 0, 4)
	for _, bitSet := range suitBitSet {
		if rank, found := findStraight(bitSet); found {
			straightRanks = append(straightRanks, rank)
		}
	}
	if len(straightRanks) > 0 {
		return HandStrength{StraightFlush, []Rank{MaxRank(straightRanks)}}
	}

	// quads
	for rank := AceHigh; rank >= Two; rank-- {
		if rankCount[rank % 13] >= 4 {
			kickers := maxRankSet(bitSet &^ rank.Mask(), 1)
			return HandStrength{Quads, append([]Rank{rank}, kickers...)}
		}
	}

	// full house
	for hiRank := AceHigh; hiRank >= Two; hiRank-- {
		if rankCount[hiRank % 13] >= 3 {
			for loRank := AceHigh; loRank >= Two; loRank-- {
				if loRank != hiRank && rankCount[loRank % 13] >= 2 {
					return HandStrength{FullHouse, []Rank{hiRank, loRank}}
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
		return HandStrength{Flush, maxRankSet(flushBitSet, 5)}
	}

	// straight
	if rank, found := findStraight(bitSet); found {
		return HandStrength{Straight, []Rank{rank}}
	}

	// trips
	for rank := AceHigh; rank >= Two; rank-- {
		if rankCount[rank % 13] >= 3 {
			kickers := maxRankSet(bitSet &^ rank.Mask(), 2)
			return HandStrength{Trips, append([]Rank{rank}, kickers...)}
		}
	}

	// two pair / pair
	for hiRank := AceHigh; hiRank >= Two; hiRank-- {
		if rankCount[hiRank % 13] >= 2 {
			for loRank := AceHigh; loRank >= Two; loRank-- {
				if loRank != hiRank && rankCount[loRank % 13] >= 2 {
					kickers := maxRankSet(bitSet &^ (hiRank.Mask() | loRank.Mask()), 1)
					return HandStrength{TwoPair, append([]Rank{hiRank, loRank}, kickers...)}
				}
			}
			kickers := maxRankSet(bitSet &^ hiRank.Mask(), 3)
			return HandStrength{Pair, append([]Rank{hiRank}, kickers...)}
		}
	}

	// high card
	return HandStrength{HighCard, maxRankSet(bitSet, 5)}
}

func findStraight(bitSet uint16) (Rank, bool) {
	// ace high straight - 0001 1110 0000 0001
	mask := uint16(0x1E01)
	if (bitSet & mask) == mask {
		return AceHigh, true
	}

	for r := King; r >= Five; r-- {
		mask = uint16(0x001F) << (r - Five)
		if (bitSet & mask) == mask {
			return r, true
		}
	}

	return 0, false
}

func maxRankSet(bitSet uint16, length int) []Rank {
	maxSet := RankSet(bitSet)
	sort.Sort(sort.Reverse(RankSlice(maxSet)))
	return maxSet[0:length]
}
