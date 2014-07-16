package poker

import (
	. "github.com/dohodges/fifty2"
	"reflect"
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

	assertStrength(t, GetHandStrength(h), HandStrength{StraightFlush, []CardStrength{AceHigh}})
	assertStrength(t, GetLowHandStrength(h), HandStrength{HighCard, []CardStrength{
		CardStrength(King),
		CardStrength(Queen),
		CardStrength(Jack),
		CardStrength(Ten),
		AceLow,
	}})

	h = []Card{
		Card{Ace, Clubs},
		Card{Two, Clubs},
		Card{Three, Clubs},
		Card{Four, Clubs},
		Card{Five, Clubs},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{StraightFlush, []CardStrength{CardStrength(Five)}})
	assertStrength(t, GetLowHandStrength(h), HandStrength{HighCard, []CardStrength{
		CardStrength(Five),
		CardStrength(Four),
		CardStrength(Three),
		CardStrength(Two),
		AceLow,
	}})
}

func TestQuads(t *testing.T) {
	h := []Card{
		Card{Seven, Clubs},
		Card{King, Hearts},
		Card{King, Clubs},
		Card{King, Diamonds},
		Card{King, Spades},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{Quads, []CardStrength{
		CardStrength(King),
		CardStrength(Seven),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{Quads, []CardStrength{
		CardStrength(King),
		CardStrength(Seven),
	}})
}

func TestFullHouse(t *testing.T) {
	h := []Card{
		Card{Five, Clubs},
		Card{Five, Diamonds},
		Card{Eight, Spades},
		Card{Five, Spades},
		Card{Eight, Clubs},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{FullHouse, []CardStrength{
		CardStrength(Five),
		CardStrength(Eight),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{FullHouse, []CardStrength{
		CardStrength(Five),
		CardStrength(Eight),
	}})
}

func TestFlush(t *testing.T) {
	h := []Card{
		Card{Two, Hearts},
		Card{Three, Hearts},
		Card{Ten, Hearts},
		Card{Four, Hearts},
		Card{Jack, Hearts},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{Flush, []CardStrength{
		CardStrength(Jack),
		CardStrength(Ten),
		CardStrength(Four),
		CardStrength(Three),
		CardStrength(Two),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{HighCard, []CardStrength{
		CardStrength(Jack),
		CardStrength(Ten),
		CardStrength(Four),
		CardStrength(Three),
		CardStrength(Two),
	}})
}

func TestStraight(t *testing.T) {
	h := []Card{
		Card{Eight, Hearts},
		Card{Nine, Spades},
		Card{Jack, Hearts},
		Card{Ten, Spades},
		Card{Seven, Diamonds},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{Straight, []CardStrength{CardStrength(Jack)}})
	assertStrength(t, GetLowHandStrength(h), HandStrength{HighCard, []CardStrength{
		CardStrength(Jack),
		CardStrength(Ten),
		CardStrength(Nine),
		CardStrength(Eight),
		CardStrength(Seven),
	}})
}

func TestTrips(t *testing.T) {
	h := []Card{
		Card{Nine, Hearts},
		Card{Eight, Spades},
		Card{Nine, Spades},
		Card{Four, Clubs},
		Card{Nine, Clubs},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{Trips, []CardStrength{
		CardStrength(Nine),
		CardStrength(Eight),
		CardStrength(Four),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{Trips, []CardStrength{
		CardStrength(Nine),
		CardStrength(Eight),
		CardStrength(Four),
	}})
}

func TestTwoPair(t *testing.T) {
	h := []Card{
		Card{Seven, Hearts},
		Card{Two, Spades},
		Card{Two, Clubs},
		Card{King, Diamonds},
		Card{Seven, Clubs},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{TwoPair, []CardStrength{
		CardStrength(Seven),
		CardStrength(Two),
		CardStrength(King),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{TwoPair, []CardStrength{
		CardStrength(Seven),
		CardStrength(Two),
		CardStrength(King),
	}})
}

func TestPair(t *testing.T) {
	h := []Card{
		Card{Ace, Hearts},
		Card{King, Hearts},
		Card{Seven, Diamonds},
		Card{Ace, Clubs},
		Card{Two, Hearts},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{Pair, []CardStrength{
		AceHigh,
		CardStrength(King),
		CardStrength(Seven),
		CardStrength(Two),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{Pair, []CardStrength{
		AceLow,
		CardStrength(King),
		CardStrength(Seven),
		CardStrength(Two),
	}})
}

func TestHighCard(t *testing.T) {
	h := []Card{
		Card{Seven, Spades},
		Card{Queen, Diamonds},
		Card{Two, Hearts},
		Card{Ace, Spades},
		Card{Three, Clubs},
	}

	assertStrength(t, GetHandStrength(h), HandStrength{HighCard, []CardStrength{
		AceHigh,
		CardStrength(Queen),
		CardStrength(Seven),
		CardStrength(Three),
		CardStrength(Two),
	}})

	assertStrength(t, GetLowHandStrength(h), HandStrength{HighCard, []CardStrength{
		CardStrength(Queen),
		CardStrength(Seven),
		CardStrength(Three),
		CardStrength(Two),
		AceLow,
	}})
}

func assertStrength(t *testing.T, actual, expect HandStrength) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expected - %v\nactual - %v", expect, actual)
	}
}
