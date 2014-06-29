package main

import (
	"flag"
	"fmt"
	. "github.com/dohodges/fifty2"
	. "github.com/dohodges/fifty2/poker"
)

func main() {

	var size int

	flag.IntVar(&size, "s", 5, "hand size")
	flag.Parse()

	totalHits := 0
	rankHits := make(map[HandRank]int)

	for hand := range Combinations(NewDeck(), size) {
		rank := GetHandStrength(hand).Rank
		rankHits[rank] += 1
		totalHits += 1
	}

	for _, rank := range HandRanks() {
		fmt.Printf("%15s %.8f\n", rank.String(), float64(rankHits[rank]) / float64(totalHits))
	}

}
