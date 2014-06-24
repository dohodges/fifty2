package fifty2

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Suit uint8

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func (s Suit) Rune() rune {
	switch s {
	case Clubs:
		return '♣'
	case Diamonds:
		return '♦'
	case Hearts:
		return '♥'
	case Spades:
		return '♠'
	}
	return 0
}

func (s Suit) Mask() uint8 {
	return uint8(1) << s
}

func Suits() []Suit {
	return []Suit{Clubs, Diamonds, Hearts, Spades}
}

func SuitSet(bitSet uint8) []Suit {
	suitSet := make([]Suit, 0, 4)
	for _, suit := range Suits() {
		if (bitSet & suit.Mask()) != 0 {
			suitSet = append(suitSet, suit)
		}
	}
	return suitSet
}

type Rank uint8

const (
	Ace Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) Mask() uint16 {
	return uint16(1) << r
}

func (r Rank) Rune() rune {
	switch r {
	case Ace:
		return 'A'
	case Two:
		return '2'
	case Three:
		return '3'
	case Four:
		return '4'
	case Five:
		return '5'
	case Six:
		return '6'
	case Seven:
		return '7'
	case Eight:
		return '8'
	case Nine:
		return '9'
	case Ten:
		return 'T'
	case Jack:
		return 'J'
	case Queen:
		return 'Q'
	case King:
		return 'K'
	}
	return 0
}

func Ranks() []Rank {
	return []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
}

func RankSet(bitSet uint16) []Rank {
	rankSet := make([]Rank, 0, 13)
	for _, rank := range Ranks() {
		if (bitSet & rank.Mask()) != 0 {
			rankSet = append(rankSet, rank)
		}
	}
	return rankSet
}

func MaxRank(ranks []Rank) Rank {
	if len(ranks) == 0 {
		panic("fifty2: cannot find max rank from empty slice")
	}
	max := ranks[0]
	for i := 1; i < len(ranks); i++ {
		if ranks[i] > max {
			max = ranks[i]
		}
	}
	return max
}

type RankSlice []Rank

func (rs RankSlice) Len() int {
	return len(rs)
}

func (rs RankSlice) Less(i, j int) bool {
	return rs[i] < rs[j]
}

func (rs RankSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%c%c", c.Rank.Rune(), c.Suit.Rune())
}

func NewDeck() []Card {
	deck := make([]Card, 52)
	index := 0
	for _, suit := range Suits() {
		for _, rank := range Ranks() {
			deck[index] = Card{Rank: rank, Suit: suit}
			index++
		}
	}
	return deck
}

func NewDeckSet(decks uint) []Card {
	deckSet := make([]Card, 52*decks)
	for d := uint(0); d < decks; d++ {
		a := d * 52
		b := a + 52
		copy(deckSet[a:b], NewDeck())
	}
	return deckSet
}

func String(slice []Card) string {
	cards := make([]string, 0, len(slice))
	for _, c := range slice {
		cards = append(cards, c.String())
	}
	return `{ ` + strings.Join(cards, ` `) + ` }`
}

func Shuffle(slice []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for src := 0; src < len(slice); src++ {
		dest := r.Intn(len(slice))
		swap := slice[src]
		slice[src] = slice[dest]
		slice[dest] = swap
	}
}

func Combinations(slice []Card, choose int) chan []Card {
	if choose > len(slice) {
		panic("fifty2: cannot produce combinations larger than given card slice")
	}

	ch := make(chan []Card)
	go func() {
		findCombinations(slice, []Card{}, choose, ch)
		close(ch)
	}()

	return ch
}

func findCombinations(slice []Card, combo []Card, choose int, ch chan []Card) {
	for i, c := range slice {
		nextCombo := make([]Card, len(combo) + 1)
		copy(nextCombo, combo)
		nextCombo[len(combo)] = c
		if choose == 1 {
			ch <- nextCombo
		} else {
			findCombinations(slice[i+1:], nextCombo, choose-1, ch)
		}
	}
}
