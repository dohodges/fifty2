package poker

import (
	. "github.com/dohodges/fifty2"
)

type GameType string

const (
	Holdem    GameType = "holdem"
	Omaha     GameType = "omaha"
	OmahaHiLo GameType = "omahahl"
	Stud7     GameType = "stud7"
	Stud7HiLo GameType = "stud7hl"
	Stud5     GameType = "stud5"
	Razz      GameType = "razz"
)

type GameStrengthFunc func(board, hand []Card) HandStrength

type Game struct {
	Name       string
	HandSize   int
	BoardSize  int
	HiStrength GameStrengthFunc
	LoStrength GameStrengthFunc
}

func (g Game) IsHiLo() bool {
	return g.LoStrength != nil
}

var games map[GameType]Game

func init() {

	games = map[GameType]Game{

		Holdem: Game{
			Name:       "Texas Hold'em",
			HandSize:   2,
			BoardSize:  5,
			HiStrength: GetHoldemHandStrength,
			LoStrength: nil,
		},

		Omaha: Game{
			Name:       "Omaha",
			HandSize:   4,
			BoardSize:  5,
			HiStrength: GetOmahaHandStrength,
			LoStrength: nil,
		},

		OmahaHiLo: Game{
			Name:       "Omaha Hi/Lo",
			HandSize:   4,
			BoardSize:  5,
			HiStrength: GetOmahaHandStrength,
			LoStrength: GetOmahaLowHandStrength,
		},

		Stud7: Game{
			Name:       "7-card Stud",
			HandSize:   7,
			BoardSize:  0,
			HiStrength: func(board, hand []Card) HandStrength { return GetHandStrength(hand) },
			LoStrength: nil,
		},

		Stud7HiLo: Game{
			Name:       "7-card Stud Hi/Lo",
			HandSize:   7,
			BoardSize:  0,
			HiStrength: func(board, hand []Card) HandStrength { return GetHandStrength(hand) },
			LoStrength: func(board, hand []Card) HandStrength { return GetLowHandStrength(hand, true) },
		},

		Stud5: Game{
			Name:       "5-card Stud",
			HandSize:   5,
			BoardSize:  0,
			HiStrength: func(board, hand []Card) HandStrength { return GetHandStrength(hand) },
			LoStrength: nil,
		},

		Razz: Game{
			Name:       "Razz",
			HandSize:   7,
			BoardSize:  0,
			HiStrength: nil,
			LoStrength: func(board, hand []Card) HandStrength { return GetLowHandStrength(hand, false) },
		},
	}
}

func GetGame(gt GameType) Game {
	return games[gt]
}

func GetHoldemHandStrength(board, pocket []Card) HandStrength {
	hand := make([]Card, 7)
	copy(hand, pocket)
	copy(hand[2:], board)
	return GetHandStrength(hand)
}

func GetOmahaHandStrength(board, pocket []Card) HandStrength {
	strengths := make([]HandStrength, 0, 6)
	for itr := Combinations(pocket, 2); itr.HasNext(); {
		hand := make([]Card, 7)
		copy(hand, itr.Next())
		copy(hand[2:], board)
		strengths = append(strengths, GetHandStrength(hand))
	}
	return MaxHandStrength(strengths)
}

func GetOmahaLowHandStrength(board, pocket []Card) HandStrength {
	strengths := make([]HandStrength, 0, 6)
	for itr := Combinations(pocket, 2); itr.HasNext(); {
		hand := make([]Card, 7)
		copy(hand, itr.Next())
		copy(hand[2:], board)
		strength := GetLowHandStrength(hand, true)
		if strength.Rank() != NoHand {
			strengths = append(strengths, strength)
		}
	}

	if len(strengths) > 0 {
		return MinHandStrength(strengths)
	}

	return MakeHandStrength(NoHand, 0, 0, 0)
}
