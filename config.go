package main

type Config struct {
	datasetFilePath           string
	nonFrenchPlayersSuffix    string
	nonOrganizerPlayerSuffix  string
	totalScoreDiminishingRate float64
	minDiminishedRate         float64
	sizeFactors               map[int]float64
	scoreSequence             []float64
}

func initConfig() Config {
	return Config{
		datasetFilePath:           "./data/2024.json",
		nonFrenchPlayersSuffix:    "_fp",
		nonOrganizerPlayerSuffix:  "_to",
		totalScoreDiminishingRate: 0.20,
		minDiminishedRate:         0.20,
		sizeFactors:               map[int]float64{0: 0.6, 8: 0.8, 16: 1, 24: 1.2, 32: 1.4},
		scoreSequence:             []float64{1, 2, 3, 5, 8, 13, 21, 34, 55},
	}
}
