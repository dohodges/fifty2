package fifty2

import (
	"bufio"
	"io"
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

func ParseSuit(r rune) (Suit, error) {
	switch r {
	case 'c', 'C', '♣':
		return Clubs, nil
	case 'd', 'D', '♦':
		return Diamonds, nil
	case 'h', 'H', '♥':
		return Hearts, nil
	case 's', 'S', '♠':
		return Spades, nil
	}
	return 0, fmt.Errorf("fifty2: unknown suit[%c]", r)
}

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

func ParseRank(r rune) (Rank, error) {
	switch r {
	case 'a', 'A':
		return Ace, nil
	case '2':
		return Two, nil
	case '3':
		return Three, nil
	case '4':
		return Four, nil
	case '5':
		return Five, nil
	case '6':
		return Six, nil
	case '7':
		return Seven, nil
	case '8':
		return Eight, nil
	case '9':
		return Nine, nil
	case 't', 'T':
		return Ten, nil
	case 'j', 'J':
		return Jack, nil
	case 'q', 'Q':
		return Queen, nil
	case 'k', 'K':
		return King, nil
	}
	return 0, fmt.Errorf("fifty2: unknown rank[%c]", r)
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

func (r Rank) Mask() uint16 {
	return uint16(1) << r
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

type CardReader struct {
	reader *bufio.Reader
}

func NewCardReader(reader io.Reader) *CardReader {
	return &CardReader{bufio.NewReader(reader)}
}

func (cr *CardReader) Read() (Card, error) {
	r, _, err := cr.reader.ReadRune()
	if err != nil {
		return Card{}, err
	}

	rank, err := ParseRank(r)
	if err != nil {
		return Card{}, err
	}

	r, _, err = cr.reader.ReadRune()
	if err != nil {
		return Card{}, err
	}

	suit, err := ParseSuit(r)
	if err != nil {
		return Card{}, err
	}

	return Card{rank, suit}, nil
}

func (cr *CardReader) ReadAll() ([]Card, error) {
	cards := make([]Card, 0)
	for  {
		card, err := cr.Read()
		if err == nil {
			cards = append(cards, card)
		} else if err == io.EOF {
			return cards, nil
		} else {
			return cards, err
		}
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

func Index(slice []Card, card Card) int {
	for i, c := range slice {
		if c == card {
			return i
		}
	}
	return -1
}

func Remove(slice []Card, card Card) []Card {
	index := Index(slice, card)
	if index < 0 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
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
