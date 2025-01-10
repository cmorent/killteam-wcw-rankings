package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func getSizeFactor(factors map[int]float64, s int) float64 {
	k := s / 8 * 8
	if k > 32 {
		k = 32
	}

	return factors[k]
}

func getSequenceScore(seq []float64, r int) float64 {
	idx := len(seq) - r
	if idx < 0 {
		idx = 0
	}

	return seq[idx]
}

func getRankingScore(eventSize, rank int) float64 {
	return float64(eventSize - rank)
}

type Player struct {
	Name        string    `json:"name"`
	Scores      []float64 `json:"score"`
	TotalScore  float64   `json:"totalScore"`
	Tournaments []string  `json:"tournaments"`
}

func (p *Player) ComputeSeasonalScore(diminishingRate, minDiminishedRate float64) {
	sort.Sort(sort.Reverse(sort.Float64Slice(p.Scores)))

	var factor float64 = 1
	for _, score := range p.Scores {
		p.TotalScore += score * factor
		factor -= diminishingRate
		if factor < minDiminishedRate {
			factor = minDiminishedRate
		}
	}
}

func displayRankings(rankings []string, players map[string]*Player, eventsData map[string][]string) error {
	w := new(tabwriter.Writer).Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t", "RANK", "PLAYER", "SCORE", "EVENTS")

	var rank int
	prevRank := -1
	prevScore := -1.0
	for idx, playerName := range rankings {
		p := players[playerName]

		rank = idx + 1
		if p.TotalScore == prevScore {
			rank = prevRank
		}
		prevRank = rank
		prevScore = p.TotalScore

		fmt.Fprintf(w, "\n %.2d\t%s\t%.2f\t%d/%d\t", rank, formatName(playerName), p.TotalScore, len(p.Tournaments), len(eventsData))
		if idx == 7 {
			fmt.Fprint(w, "\n --\t------------\t-----\t---\t")
		}
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func main() {
	cfg := initConfig()

	// Load dataset
	fs, err := os.Open(cfg.datasetFilePath)
	if err != nil {
		slog.Error("failed to open dataset", slog.Any("error", err))
	}
	defer fs.Close()

	// Decode the dataset
	// map[event][]rankings{P1, P2, ..., Pn}
	var eventsData map[string][]string
	err = json.NewDecoder(fs).Decode(&eventsData)
	if err != nil {
		panic(err)
	}

	// Parse Rankings, keep an index of players
	players := make(map[string]*Player)
	for eventName, rankings := range eventsData {
		eventSize := len(rankings)
		sizeFactor := getSizeFactor(cfg.sizeFactors, eventSize)
		debug("%s - %d players (%.2f)\n", eventName, eventSize, sizeFactor)

		for idx, playerName := range rankings {
			playerRanking := idx + 1

			// Ignore foreigners and TO if any
			if strings.HasSuffix(playerName, cfg.nonFrenchPlayersSuffix) ||
				strings.HasSuffix(playerName, cfg.nonOrganizerPlayerSuffix) {
				continue
			}

			// Attempt to retrieve cached player, create a new one if cache miss
			// and then add score data for the event.
			p, ok := players[playerName]
			if !ok {
				p = &Player{playerName, []float64{}, 0.0, []string{}}
				players[playerName] = p
			}

			score := getSequenceScore(cfg.scoreSequence, playerRanking)
			score += getRankingScore(eventSize, playerRanking)
			score *= sizeFactor

			p.Scores = append(
				p.Scores,
				score,
			)
			p.Tournaments = append(p.Tournaments, eventName)
		}
	}

	// Compute Scores
	rankings := make([]string, 0, len(players))
	for playerName, p := range players {
		rankings = append(rankings, playerName)
		p.ComputeSeasonalScore(cfg.totalScoreDiminishingRate, cfg.minDiminishedRate)
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return players[rankings[i]].TotalScore > players[rankings[j]].TotalScore
	})

	if err = displayRankings(rankings, players, eventsData); err != nil {
		log.Fatal(err)
	}
}
