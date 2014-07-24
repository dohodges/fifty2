package poker

import (
	. "github.com/dohodges/fifty2"
	"testing"
)

func BenchmarkGetHandStrength5(b *testing.B) {
	deck := NewDeck()
	for i := 0; i < b.N; i++ {
		Shuffle(deck)
		hand := deck[:5]
		GetHandStrength(hand)
	}
}

func BenchmarkGetHandStrength7(b *testing.B) {
	deck := NewDeck()
	for i := 0; i < b.N; i++ {
		Shuffle(deck)
		hand := deck[:7]
		GetHandStrength(hand)
	}
}

func BenchmarkGetHandStrength13(b *testing.B) {
	deck := NewDeck()
	for i := 0; i < b.N; i++ {
		Shuffle(deck)
		hand := deck[:13]
		GetHandStrength(hand)
	}
}

func TestStraightFlush(t *testing.T) {
	h := []Card{
		Card{Ace, Spades},
		Card{King, Spades},
		Card{Queen, Spades},
		Card{Jack, Spades},
		Card{Ten, Spades},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(StraightFlush, AceHigh, 0, 0))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x1E01))

	h = []Card{
		Card{Ace, Clubs},
		Card{Two, Clubs},
		Card{Three, Clubs},
		Card{Four, Clubs},
		Card{Five, Clubs},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(StraightFlush, CardStrength(Five), 0, 0))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x001F))
}

func TestQuads(t *testing.T) {
	h := []Card{
		Card{Seven, Clubs},
		Card{King, Hearts},
		Card{King, Clubs},
		Card{King, Diamonds},
		Card{King, Spades},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(Quads, CardStrength(King), 0, 0x0040))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(Quads, CardStrength(King), 0, 0x0040))
}

func TestFullHouse(t *testing.T) {
	h := []Card{
		Card{Five, Clubs},
		Card{Five, Diamonds},
		Card{Eight, Spades},
		Card{Five, Spades},
		Card{Eight, Clubs},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(FullHouse, CardStrength(Five), CardStrength(Eight), 0))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(FullHouse, CardStrength(Five), CardStrength(Eight), 0))
}

func TestFlush(t *testing.T) {
	h := []Card{
		Card{Two, Hearts},
		Card{Three, Hearts},
		Card{Ten, Hearts},
		Card{Four, Hearts},
		Card{Jack, Hearts},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(Flush, 0, 0, 0x060E))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x060E))
}

func TestStraight(t *testing.T) {
	h := []Card{
		Card{Eight, Hearts},
		Card{Nine, Spades},
		Card{Jack, Hearts},
		Card{Ten, Spades},
		Card{Seven, Diamonds},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(Straight, CardStrength(Jack), 0, 0))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x07C0))
}

func TestTrips(t *testing.T) {
	h := []Card{
		Card{Nine, Hearts},
		Card{Eight, Spades},
		Card{Nine, Spades},
		Card{Four, Clubs},
		Card{Nine, Clubs},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(Trips, CardStrength(Nine), 0, 0x088))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(Trips, CardStrength(Nine), 0, 0x088))
}

func TestTwoPair(t *testing.T) {
	h := []Card{
		Card{Seven, Hearts},
		Card{Two, Spades},
		Card{Two, Clubs},
		Card{King, Diamonds},
		Card{Seven, Clubs},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(TwoPair, CardStrength(Seven), CardStrength(Two), 0x1000))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(TwoPair, CardStrength(Seven), CardStrength(Two), 0x1000))
}

func TestPair(t *testing.T) {
	h := []Card{
		Card{Ace, Hearts},
		Card{King, Hearts},
		Card{Seven, Diamonds},
		Card{Ace, Clubs},
		Card{Two, Hearts},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(Pair, AceHigh, 0, 0x1042))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(Pair, AceLow, 0, 0x1042))
}

func TestHighCard(t *testing.T) {
	h := []Card{
		Card{Seven, Spades},
		Card{Queen, Diamonds},
		Card{Two, Hearts},
		Card{Ace, Spades},
		Card{Three, Clubs},
	}

	assertStrength(t, GetHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x2846))
	assertStrength(t, GetLowHandStrength(h), MakeHandStrength(HighCard, 0, 0, 0x0847))
}

func assertStrength(t *testing.T, actual, expect HandStrength) {
	if expect != actual {
		t.Errorf("expected - %#X\nactual - %#X", expect, actual)
	}
}
