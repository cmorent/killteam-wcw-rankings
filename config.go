package main

var (
	debugEnabled bool = true

	datasetFilePath = "./data/2024.json"

	nonFrenchPlayersSuffix = "_nfp"

	totalScoreDiminishingRate float64 = 0.20

	// true
	sizeFactorEnabled bool = true
	sizeFactors            = map[int]float64{
		0:  0.6,
		8:  0.8,
		16: 1,
		24: 1.2,
		32: 1.4,
	}

	rankingScoreSequences = [][]float64{
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
