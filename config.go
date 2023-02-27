package main

var (
	debug bool = false

	datasetFilePath = "./data.json"

	topEightOnly = true

	nonFrenchPlayersIgnore = true
	nonFrenchPlayersSuffix = "_nfp"

	totalScoreCountBestFirst  bool    = true
	totalScoreDiminishingRate float64 = 0.8
	totalScoreMaxEvents       int     = 0
	totalScoreMaxEventBonus   float64 = 0

	// true
	sizeFactorEnabled   bool = true
	sizeFactorUseMedian bool = true

	// MaxRounds ?
	// MaxRoundsUndefeated ?

	rankingScoreSequences = [][]int{
		// 0 - Custom
		{1, 2, 3, 4, 5, 7, 8, 10, 12},
		// 1 - Fibonacci
		{0, 1, 2, 3, 5, 8, 13, 21, 34},
		// 2 - Fibonacci - XXX
		{1, 2, 3, 5, 8, 13, 21, 34, 55},
		// 3 - Lucas
		{1, 3, 4, 7, 11, 18, 29, 47, 76},
		// 4 - Power of 2
		{1, 2, 4, 8, 16, 32, 64, 128, 256},
	}
	rankingScoreSequencesIdx int = 2
)
