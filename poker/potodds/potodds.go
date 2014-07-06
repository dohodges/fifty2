package main

import (
	"flag"
	"fmt"
	. "github.com/dohodges/fifty2"
	. "github.com/dohodges/fifty2/poker"
	"os"
	"strings"
)

type Game struct {
	Flag            string
	Name            string
	HiLo            bool
	MinHandSize     int
	MaxHandSize     int
	MaxBoardSize    int
	GetHandStrength func(board, hand []Card) HandStrength
}

var games []Game

func init() {

	Holdem := Game{
		Flag:         "holdem",
		Name:         "Texas Hold'em",
		HiLo:         false,
		MinHandSize:  2,
		MaxHandSize:  2,
		MaxBoardSize: 5,
		GetHandStrength: func(board, pocket []Card) HandStrength {
			hand := make([]Card, 7)
			copy(hand, pocket)
			copy(hand[2:], board)
			return GetHandStrength(hand)
		},
	}

	Omaha := Game{
		Flag:         "omaha",
		Name:         "Omaha",
		HiLo:         false,
		MinHandSize:  4,
		MaxHandSize:  4,
		MaxBoardSize: 5,
		GetHandStrength: func(board, pocket []Card) HandStrength {
			strengths := make([]HandStrength, 0, 6)
			for p := range Combinations(pocket, 2) {
				hand := make([]Card, 7)
				copy(hand, p)
				copy(hand[2:], board)
				strengths = append(strengths, GetHandStrength(hand))
			}
			return MaxHandStrength(strengths)
		},
	}

	Stud7 := Game{
		Flag:         "stud7",
		Name:         "7-card Stud",
		HiLo:         false,
		MinHandSize:  0,
		MaxHandSize:  7,
		MaxBoardSize: 0,
		GetHandStrength: func(board, hand []Card) HandStrength {
			return GetHandStrength(hand)
		},
	}

	Stud5 := Game{
		Flag:         "stud5",
		Name:         "5-card Stud",
		HiLo:         false,
		MinHandSize:  0,
		MaxHandSize:  5,
		MaxBoardSize: 0,
		GetHandStrength: func(board, hand []Card) HandStrength {
			return GetHandStrength(hand)
		},
	}

	games = []Game{Holdem, Omaha, Stud7, Stud5}
}

func GetGame(flag string) (Game, bool) {
	for _, game := range games {
		if game.Flag == flag {
			return game, true
		}
	}
	return Game{}, false
}

func main() {

	var (
		gameFlag  string
		boardFlag string
		hiLoFlag  bool
	)

	flag.StringVar(&gameFlag, "game", "holdem", "game")
	flag.StringVar(&boardFlag, "board", "", "community cards")
	flag.BoolVar(&hiLoFlag, "hilo", false, "hi/lo split")
	flag.Parse()

	game, ok := GetGame(gameFlag)
	if !ok {
		fmt.Printf("unknown game - %s", gameFlag)
		os.Exit(1)
	}

	if hiLoFlag && !game.HiLo {
		fmt.Printf("%s does not have a Hi/Lo variant\n", game.Name)
		os.Exit(1)
	}

	board, err := NewCardReader(strings.NewReader(boardFlag)).ReadAll()
	if err != nil {
		fmt.Printf("invalid board - %v\n", err)
		os.Exit(1)
	} else if len(board) > game.MaxBoardSize {
		fmt.Printf("%s has a maximum of %d community cards\n", game.Name, game.MaxBoardSize)
		os.Exit(1)
	}

	hand, err := NewCardReader(strings.NewReader(flag.Arg(0))).ReadAll()
	if err != nil {
		fmt.Printf("invalid hand - %v\n", err)
		os.Exit(1)
	} else if len(hand) < game.MinHandSize {
		fmt.Printf("%s has a minimum hand size of %d\n", game.Name, game.MinHandSize)
		os.Exit(1)
	} else if len(hand) > game.MaxHandSize {
		fmt.Printf("%s has a maximum hand size of %d\n", game.Name, game.MaxHandSize)
		os.Exit(1)
	}

	deck := NewDeck()
	for _, card := range board {
		deck = Remove(deck, card)
	}
	for _, card := range hand {
		deck = Remove(deck, card)
	}

	totalHits := 0
	rankHits := make(map[HandRank]int)

	choose := (game.MaxHandSize - len(hand)) + (game.MaxBoardSize - len(board))
	for combo := range Combinations(deck, choose) {
		var fullBoard, fullHand []Card
		if len(board) < game.MaxBoardSize {
			fullHand = hand
			fullBoard = make([]Card, game.MaxBoardSize)
			copy(fullBoard, board)
			copy(fullBoard[len(board):], combo)
		} else {
			fullBoard = board
			fullHand = make([]Card, game.MaxHandSize)
			copy(fullHand, hand)
			copy(fullHand[len(hand):], combo)
		}

		strength := game.GetHandStrength(fullBoard, fullHand)
		rankHits[strength.Rank]++
		totalHits++
	}

	for _, rank := range HandRanks() {
		if rankHits[rank] > 0 {
			fmt.Printf("%15s %.8f\n", rank.String(), float64(rankHits[rank])/float64(totalHits))
		}
	}

}
