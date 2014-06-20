package fifty2

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Suit uint8

const (
	SuitClubs Suit = 1 << iota
	SuitDiamonds
	SuitHearts
	SuitSpades
)

func (s Suit) Rune() rune {
	switch s {
	case SuitClubs:
		return '♣'
	case SuitDiamonds:
		return '♦'
	case SuitHearts:
		return '♥'
	case SuitSpades:
		return '♠'
	}
	return 0
}

func Suits() []Suit {
	return []Suit{SuitClubs, SuitDiamonds, SuitHearts, SuitSpades}
}

type Rank uint16

const (
	RankAce Rank = 1 << iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	Rank9
	Rank10
	RankJack
	RankQueen
	RankKing
)

func Ranks() []Rank {
	return []Rank{RankAce, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8, Rank9, Rank10, RankJack, RankQueen, RankKing}
}

func (r Rank) Rune() rune {
	switch r {
	case RankAce:
		return 'A'
	case Rank2:
		return '2'
	case Rank3:
		return '3'
	case Rank4:
		return '4'
	case Rank5:
		return '5'
	case Rank6:
		return '6'
	case Rank7:
		return '7'
	case Rank8:
		return '8'
	case Rank9:
		return '9'
	case Rank10:
		return 'T'
	case RankJack:
		return 'J'
	case RankQueen:
		return 'Q'
	case RankKing:
		return 'K'
	}
	return 0
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank.Rune(), c.Suit.Rune())
}

func String(set []Card) string {
	cards := make([]string, 0, len(set))
	for _, c := range set {
		cards = append(cards, c.String())
	}
	return `{ ` + strings.Join(cards, ` `) + ` }`
}

func Shuffle(set []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for src := 0; src < len(set); src++ {
		dest := r.Intn(len(set))
		swap := set[src]
		set[src] = set[dest]
		set[dest] = swap
	}
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

func NewDeckSet(decks uint8) []Card {
	deckSet := make([]Card, 52*decks)
	for d := uint8(0); d < decks; d++ {
		a := d * 52
		b := a + 52
		copy(deckSet[a:b], NewDeck())
	}
	return deckSet
}
