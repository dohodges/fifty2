package poker

import (
	. "github.com/dohodges/fifty2"
	"reflect"
	"testing"
)

func TestStraightFlush(t *testing.T) {

	s := GetHandStrength([]Card{
		Card{Ace, Spades},
		Card{King, Spades},
		Card{Queen, Spades},
		Card{Jack, Spades},
		Card{Ten, Spades},
	})

	assertStrength(t, s, HandStrength{StraightFlush, []CardStrength{AceHigh}})

	s = GetHandStrength([]Card{
		Card{Ace, Clubs},
		Card{Two, Clubs},
		Card{Three, Clubs},
		Card{Four, Clubs},
		Card{Five, Clubs},
	})

	assertStrength(t, s, HandStrength{StraightFlush, []CardStrength{CardStrength(Five)}})
}

func TestQuads(t *testing.T) {

	s := GetHandStrength([]Card{
		Card{Seven, Clubs},
		Card{King, Hearts},
		Card{King, Clubs},
		Card{King, Diamonds},
		Card{King, Spades},
	})

	assertStrength(t, s, HandStrength{Quads, []CardStrength{CardStrength(King), CardStrength(Seven)}})
}

func TestFullHouse(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Five, Clubs},
		Card{Five, Diamonds},
		Card{Eight, Spades},
		Card{Five, Spades},
		Card{Eight, Clubs},
	})

	assertStrength(t, s, HandStrength{FullHouse, []CardStrength{CardStrength(Five), CardStrength(Eight)}})
}

func TestFlush(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Two, Hearts},
		Card{Three, Hearts},
		Card{Ten, Hearts},
		Card{Four, Hearts},
		Card{Jack, Hearts},
	})

	assertStrength(t, s, HandStrength{Flush, []CardStrength{
		CardStrength(Jack),
		CardStrength(Ten),
		CardStrength(Four),
		CardStrength(Three),
		CardStrength(Two),
	}})
}

func TestStraight(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Eight, Hearts},
		Card{Nine, Spades},
		Card{Jack, Hearts},
		Card{Ten, Spades},
		Card{Seven, Diamonds},
	})

	assertStrength(t, s, HandStrength{Straight, []CardStrength{CardStrength(Jack)}})
}

func TestTrips(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Nine, Hearts},
		Card{Eight, Spades},
		Card{Nine, Spades},
		Card{Four, Clubs},
		Card{Nine, Clubs},
	})

	assertStrength(t, s, HandStrength{Trips, []CardStrength{
		CardStrength(Nine),
		CardStrength(Eight),
		CardStrength(Four),
	}})
}

func TestTwoPair(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Seven, Hearts},
		Card{Two, Spades},
		Card{Two, Clubs},
		Card{King, Diamonds},
		Card{Seven, Clubs},
	})

	assertStrength(t, s, HandStrength{TwoPair, []CardStrength{
		CardStrength(Seven),
		CardStrength(Two),
		CardStrength(King),
	}})
}

func TestPair(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Ace, Hearts},
		Card{King, Hearts},
		Card{Seven, Diamonds},
		Card{Ace, Clubs},
		Card{Two, Hearts},
	})

	assertStrength(t, s, HandStrength{Pair, []CardStrength{
		AceHigh,
		CardStrength(King),
		CardStrength(Seven),
		CardStrength(Two),
	}})
}

func TestHighCard(t *testing.T) {
	s := GetHandStrength([]Card{
		Card{Seven, Spades},
		Card{Queen, Diamonds},
		Card{Two, Hearts},
		Card{Ace, Spades},
		Card{Three, Clubs},
	})

	assertStrength(t, s, HandStrength{HighCard, []CardStrength{
		AceHigh,
		CardStrength(Queen),
		CardStrength(Seven),
		CardStrength(Three),
		CardStrength(Two),
	}})
}

func assertStrength(t *testing.T, actual, expect HandStrength) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("expected - %v\nactual - %v", expect, actual)
	}
}
