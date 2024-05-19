package main

type Config struct {
	datasetFilePath           string
	nonFrenchPlayersSuffix    string
	totalScoreDiminishingRate float64
	sizeFactors               map[int]float64
	scoreSequence             []float64
	rankingScoreSequencesIdx  int
}

func initConfig() Config {
	return Config{
		datasetFilePath:           "./data/2024.json",
		nonFrenchPlayersSuffix:    "_fp",
		totalScoreDiminishingRate: 0.20,
		sizeFactors:               map[int]float64{0: 0.6, 8: 0.8, 16: 1, 24: 1.2, 32: 1.4},
		scoreSequence:             []float64{1, 2, 3, 5, 8, 13, 21, 34, 55},
		rankingScoreSequencesIdx:  2,
	}
}
