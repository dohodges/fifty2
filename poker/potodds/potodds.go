package main

import (
	"flag"
	"fmt"
	. "github.com/dohodges/fifty2"
	. "github.com/dohodges/fifty2/poker"
	"math"
	"os"
	"runtime/pprof"
	"strings"
)

var (
	game      Game
	board     []Card
	hands     [][]Card
	choose    []int
	fullBoard []Card
	fullHands [][]Card
)

type GameTally []*Tally

func NewGameTally(players int) GameTally {
	gt := make(GameTally, players)
	for i := 0; i < players; i++ {
		gt[i] = &Tally{}
	}
	return gt
}

func (gt GameTally) Add(gt2 GameTally) {
	for i := 0; i < len(gt); i++ {
		gt[i].Add(gt2[i])
	}
}

func (gt GameTally) Clone() GameTally {
	clone := make([]*Tally, len(gt))
	for i := range gt {
		clone[i] = gt[i].Clone()
	}
	return clone
}

func (gt GameTally) Delta(gt2 GameTally) float64 {
	var absDelta, deltas float64
	for i := range gt {
		absDelta += gt[i].Delta(gt2[i])
		deltas++
	}
	return absDelta / deltas
}

type Tally struct {
	Scoops int64
	HiWins int64
	LoWins int64
	HiTies int64
	LoTies int64
	Total  int64
}

func (t *Tally) Clone() *Tally {
	return &Tally{
		Scoops: t.Scoops,
		HiWins: t.HiWins,
		LoWins: t.LoWins,
		HiTies: t.HiTies,
		LoTies: t.LoTies,
		Total:  t.Total,
	}
}

func (t *Tally) Add(t2 *Tally) {
	t.Scoops += t2.Scoops
	t.HiWins += t2.HiWins
	t.LoWins += t2.LoWins
	t.HiTies += t2.HiTies
	t.LoTies += t2.LoTies
	t.Total += t2.Total
}

func (t *Tally) ScoopOdds() float64 {
	return 100. * float64(t.Scoops) / float64(t.Total)
}

func (t *Tally) HiWinOdds() float64 {
	return 100. * float64(t.HiWins) / float64(t.Total)
}

func (t *Tally) LoWinOdds() float64 {
	return 100. * float64(t.LoWins) / float64(t.Total)
}

func (t *Tally) HiTieOdds() float64 {
	return 100. * float64(t.HiTies) / float64(t.Total)
}

func (t *Tally) LoTieOdds() float64 {
	return 100. * float64(t.LoTies) / float64(t.Total)
}

func (t *Tally) Delta(t2 *Tally) float64 {
	return (math.Abs(t.HiWinOdds() - t2.HiWinOdds()) + math.Abs(t.LoWinOdds() - t2.LoWinOdds())) / 2.
}

func main() {

	var (
		gameFlag  string
		boardFlag string
		approx    bool
		profile   string
	)

	flag.StringVar(&gameFlag, "game", string(Holdem), "game")
	flag.StringVar(&boardFlag, "board", "", "community cards")
	flag.BoolVar(&approx, "approx", false, "approximate")
	flag.StringVar(&profile, "profile", "", "create cpu profile")
	flag.Parse()

	if profile != "" {
		f, err := os.Create(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	gameTally := NewGameTally(len(hands))

	if approx {
		iterations := 0
		for {
			lastTally := gameTally.Clone()
			for i := 0; i < 100; i++ {
				Shuffle(deck)
				deal := deck[:deckChoose]
				gameTally.Add(TallyDeal(deal))
				iterations++
			}
			if iterations > 100 && gameTally.Delta(lastTally) < .001 {
				fmt.Printf("Iterations - %d\n", iterations)
				break
			}
		}
	} else {
		// tally each possible outcome
		fmt.Printf("Combinations - %d\n", combination(len(deck), deckChoose))
		for itr := Combinations(deck, deckChoose); itr.HasNext(); {
			gameTally.Add(TallyDeal(itr.Next()))
		}
	}

	// results
	fmt.Printf("Game - %s\n", game.Name)
	if game.BoardSize > 0 {
		fmt.Printf("Board %s\n", board)
	}
	for i, tally := range gameTally {
		if game.IsHiLo() {
			fmt.Printf("Player %2d - Scoop: %6.2f%%  HiWin: %6.2f%%  LoWin: %6.2f%% HiTie: %6.2f%%  LoTie: %6.2f%%  %s\n",
			i+1, tally.ScoopOdds(), tally.HiWinOdds(), tally.LoWinOdds(), tally.HiTieOdds(), tally.LoTieOdds(), hands[i])
		} else if game.HasHiHand() {
			fmt.Printf("Player %2d - win: %6.2f%%  tie: %6.2f%%  %s\n", i+1,
				tally.HiWinOdds(), tally.HiTieOdds(), hands[i])
		} else if game.HasLoHand() {
			fmt.Printf("Player %2d - win: %6.2f%%  tie: %6.2f%%  %s\n", i+1,
				tally.LoWinOdds(), tally.LoTieOdds(), hands[i])
		}
	}

}

func TallyDeal(deal []Card) GameTally {
	tally := NewGameTally(len(hands))

	var hiStrengths, loStrengths []HandStrength
	if game.HasHiHand() {
		hiStrengths = make([]HandStrength, len(fullHands))
	}
	if game.HasLoHand() {
		loStrengths = make([]HandStrength, len(fullHands))
	}

	// each possible deal
	for itr := MultipleCombinations(deal, choose); itr.HasNext(); {
		dealCombo := itr.Next()
		copy(fullBoard[len(board):], dealCombo[0])
		for i, fullHand := range fullHands {
			copy(fullHand[len(hands[i]):], dealCombo[i+1])
			if game.HasHiHand() {
				hiStrengths[i] = game.HiStrength(fullBoard, fullHand)
			}
			if game.HasLoHand() {
				loStrengths[i] = game.LoStrength(fullBoard, fullHand)
			}
		}

		// tally wins & ties
		bestHi := GetBestHiHands(hiStrengths)
		bestLo := GetBestLoHands(loStrengths)

		if len(bestHi) == 1 && len(bestLo) == 1 && bestHi[0] == bestLo[0] {
			tally[bestHi[0]].Scoops++
		} else {
			if len(bestHi) == 1 {
				tally[bestHi[0]].HiWins++
			} else if len(bestHi) > 1 {
				for _, h := range bestHi {
					tally[h].HiTies++
				}
			}
			if len(bestLo) == 1 {
				tally[bestLo[0]].LoWins++
			} else if len(bestLo) > 1 {
				for _, h := range bestLo {
					tally[h].LoTies++
				}
			}
		}

		for _, t := range tally {
			t.Total++
		}
	}

	return tally
}

func GetBestHiHands(strengths []HandStrength) []int {
	max := MakeHandStrength(NoHand, 0, 0, 0)
	best := make([]int, 0, len(strengths))
	for i, strength := range strengths {
		if strength > max {
			max = strength
			best = append(best[0:0], i)
		} else if strength == max {
			best = append(best, i)
		}
	}
	return best
}

func GetBestLoHands(strengths []HandStrength) []int {
	min := MakeHandStrength(StraightFlush, 0, 0, 0)
	best := make([]int, 0, len(strengths))
	for i, strength := range strengths {
		if strength.Rank() != NoHand && strength < min {
			min = strength
			best = append(best[0:0], i)
		} else if strength == min {
			best = append(best, i)
		}
	}
	return best
}

func combination(n, k int) int64 {
	c := int64(n)
	for i := int64(1); i < int64(k); i++ {
		c *= (int64(n) - i)
		c /= i + 1
	}
	return c
}
