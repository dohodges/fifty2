package poker

import (
	. "github.com/dohodges/fifty2"
)

type GameType string

const (
	Holdem GameType = "holdem"
	Omaha  GameType = "omaha"
	Stud7  GameType = "stud7"
	Stud5  GameType = "stud5"
)

type GameStrengthFunc func(board, hand []Card) (HandStrength, error)

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

		Stud7: Game{
			Name:       "7-card Stud",
			HandSize:   7,
			BoardSize:  0,
			HiStrength: func(board, hand []Card) (HandStrength, error) { return GetHandStrength(hand), nil },
			LoStrength: nil,
		},

		Stud5: Game{
			Name:       "5-card Stud",
			HandSize:   5,
			BoardSize:  0,
			HiStrength: func(board, hand []Card) (HandStrength, error) { return GetHandStrength(hand), nil },
			LoStrength: nil,
		},
	}
}

func GetGame(gt GameType) Game {
	return games[gt]
}

func GetHoldemHandStrength(board, pocket []Card) (HandStrength, error) {
	hand := make([]Card, 7)
	copy(hand, pocket)
	copy(hand[2:], board)
	return GetHandStrength(hand), nil
}

func GetOmahaHandStrength(board, pocket []Card) (HandStrength, error) {
	strengths := make([]HandStrength, 0, 6)
	for p := range Combinations(pocket, 2) {
		hand := make([]Card, 7)
		copy(hand, p)
		copy(hand[2:], board)
		strengths = append(strengths, GetHandStrength(hand))
	}
	return MaxHandStrength(strengths), nil
}

func GetOmahaLowHandStrength(board, pocket []Card) (HandStrength, error) {
	strengths := make([]HandStrength, 0, 6)
	for p := range Combinations(pocket, 2) {
		hand := make([]Card, 7)
		copy(hand, p)
		copy(hand[2:], board)
		strengths = append(strengths, GetHandStrength(hand))
	}
	return MinHandStrength(strengths), nil
}
