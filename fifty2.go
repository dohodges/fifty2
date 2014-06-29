package fifty2

import (
	"fmt"
	"math/rand"
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

func Shuffle(slice []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for src := 0; src < len(slice); src++ {
		dest := r.Intn(len(slice))
		slice[dest], slice[src] = slice[src], slice[dest]
	}
}

func Combinations(slice []Card, choose int) chan []Card {
	if choose > len(slice) {
		panic("fifty2: cannot produce combinations larger than given card slice")
	}

	ch := make(chan []Card)
	go func() {
		findCombinations(slice, make([]Card, choose), choose, ch)
		close(ch)
	}()

	return ch
}

func findCombinations(slice []Card, combo []Card, choose int, ch chan []Card) {
	for i, c := range slice {
		combo[len(combo)-choose] = c
		if choose == 1 {
			result := make([]Card, len(combo))
			copy(result, combo)
			ch <- result
		} else {
			findCombinations(slice[i+1:], combo, choose-1, ch)
		}
	}
}
