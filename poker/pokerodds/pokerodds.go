package main

import (
	"fmt"
	. "github.com/dohodges/fifty2"
	. "github.com/dohodges/fifty2/poker"
)

func main() {

	combos := 0
	rankCombos := make(map[HandRank]int)

	for hand := range Combinations(NewDeck(), 5) {
		strength := GetHandStrength(hand)
		rankCombos[strength.Rank] += 1
		combos += 1
	}

	for _, rank := range HandRanks() {
		fmt.Printf("%15s %.8f\n", rank.String(), float64(rankCombos[rank]) / float64(combos))
	}

}
