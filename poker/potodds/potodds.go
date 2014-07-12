package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb"
	. "github.com/dohodges/fifty2"
	. "github.com/dohodges/fifty2/poker"
	"os"
	"reflect"
	"strings"
)

var (
	game Game
	board []Card
	hands [][]Card
	choose []int
	fullBoard []Card
	fullHands [][]Card
)

type Tally struct {
	Wins int64
	Ties int64
	Losses int64
}

func (t Tally) WinOdds() float64 {
	return 100. * float64(t.Wins) / float64(t.Total())
}

func (t Tally) TieOdds() float64 {
	return 100. * float64(t.Ties) / float64(t.Total())
}

func (t Tally) LossOdds() float64 {
	return 100. * float64(t.Losses) / float64(t.Total())
}

func (t Tally) Total() int64 {
	return t.Wins + t.Ties + t.Losses
}

func (t Tally) Add(t2 Tally) Tally {
	t.Wins += t2.Wins
	t.Ties += t2.Ties
	t.Losses += t2.Losses
	return t
}

func main() {

	var (
		gameFlag  string
		boardFlag string
	)

	flag.StringVar(&gameFlag, "game", string(Holdem), "game")
	flag.StringVar(&boardFlag, "board", "", "community cards")
	flag.Parse()

	game = GetGame(GameType(gameFlag))
	if game.Name == "" {
		fmt.Printf("potodds: unknown game - %s\n", gameFlag)
		os.Exit(1)
	}

	var err error
	board, err = NewCardReader(strings.NewReader(boardFlag)).ReadAll()
	if err != nil {
		fmt.Printf("potodds: invalid board - %v\n", err)
		os.Exit(1)
	} else if len(board) > game.BoardSize {
		fmt.Printf("potodds: %s has a maximum of %d community cards\n", game.Name, game.BoardSize)
		os.Exit(1)
	}

	hands = make([][]Card, flag.NArg())
	for i, arg := range flag.Args() {
		hand, err := NewCardReader(strings.NewReader(arg)).ReadAll()
		if err != nil {
			fmt.Printf("potodds: invalid hand - %v\n", err)
			os.Exit(1)
		} else if len(hand) > game.HandSize {
			fmt.Printf("potodds: %s has a maximum hand size of %d\n", game.Name, game.HandSize)
			os.Exit(1)
		}
		hands[i] = hand
	}

	if len(hands) < 2 {
		fmt.Printf("potodds: specify at least 2 hands\n")
		os.Exit(1)
	}

	deck := NewDeck()
	deck = Remove(deck, board...)
	for _, hand := range hands {
		deck = Remove(deck, hand...)
	}

	// copy known cards to full board and hands
	fullBoard = make([]Card, game.BoardSize)
	copy(fullBoard, board)
	fullHands = make([][]Card, len(hands))
	for i, hand := range hands {
		fullHands[i] = make([]Card, game.HandSize)
		copy(fullHands[i], hand)
	}

	// determine # cards to deal to board and each hand
	deckChoose := game.BoardSize - len(board)
	choose = make([]int, len(hands)+1)
	choose[0] = game.BoardSize - len(board)
	for i, hand := range hands {
		deckChoose += game.HandSize - len(hand)
		choose[i+1] = game.HandSize - len(hand)
	}

	progress := pb.New64(combination(len(deck), deckChoose))
	progress.Start()

	// tally each possible outcome
	tallys := make([]Tally, len(hands))
	for deal := range Combinations(deck, deckChoose) {
		dealTallys := TallyDeal(deal)
		for i := 0; i < len(tallys); i++ {
			tallys[i] = tallys[i].Add(dealTallys[i])
		}
		progress.Increment()
	}
	progress.Finish()

	// results
	fmt.Printf("Game - %s\n", game.Name)
	if game.BoardSize > 0 {
		fmt.Printf("Board %s\n", board)
	}
	for i, tally := range tallys {
		fmt.Printf("Player %2d - win: %6.2f%%  tie: %6.2f%%  lose: %6.2f%%  %s\n", i+1,
			tally.WinOdds(), tally.TieOdds(), tally.LossOdds(), hands[i])
	}

}

func TallyDeal(deal []Card) []Tally {
	tallys := make([]Tally, len(hands))

	// each possible deal
	for dealCombo := range MultipleCombinations(deal, choose) {
		hiStrengths := make([]HandStrength, len(fullHands))
		copy(fullBoard[len(board):], dealCombo[0])
		for i, fullHand := range fullHands {
			copy(fullHand[len(hands[i]):], dealCombo[i+1])
			strength, err := game.HiStrength(fullBoard, fullHand)
			if err != nil {
				strength = HandStrength{} // invalid hand
			}
			hiStrengths[i] = strength
		}

		// tally wins/losses/ties
		max := MaxHandStrength(hiStrengths)
		best := make([]int, 0, len(hiStrengths))
		for i, strength := range hiStrengths {
			if reflect.DeepEqual(strength, max) {
				best = append(best, i)
			} else {
				tallys[i].Losses++
			}
		}
		if len(best) > 1 {
			for i := range best {
				tallys[best[i]].Ties++
			}
		} else {
			tallys[best[0]].Wins++
		}
	}

	return tallys
}

func combination(n, k int) int64 {
	c := int64(n)
	for i := int64(1); i < int64(k); i++ {
		c *= (int64(n)-i)
		c /= i+1
	}
	return c
}
