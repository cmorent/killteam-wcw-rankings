package scoring

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
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

// func getRankingScore(eventSize, rank int) float64 {
// 	return float64(eventSize - rank)
// }

func ComputeSeasonalRankings(eventsRankings map[string][]string) ([]*Player, error) {
	cfg := initConfig()

	// Parse Rankings, keep an index of players
	players := make(map[string]*Player)
	for eventName, rankings := range eventsRankings {
		eventSize := len(rankings)
		sizeFactor := getSizeFactor(cfg.sizeFactors, eventSize)
		slog.Info(fmt.Sprintf("%s %d players (%.2f)", eventName, eventSize, sizeFactor))

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
			// score += getRankingScore(eventSize, playerRanking)
			score *= sizeFactor

			p.Scores = append(
				p.Scores,
				score,
			)
			p.Tournaments = append(p.Tournaments, eventName)
		}
	}

	// Compute total scores
	rankings := []*Player{}
	for _, p := range players {
		p.ComputeTotalScore(cfg.totalScoreDiminishingRate, cfg.minDiminishedRate)
		rankings = append(rankings, p)
	}

	// Sort the table
	sort.SliceStable(rankings, func(i, j int) bool {
		return rankings[i].TotalScore > rankings[j].TotalScore
	})

	return rankings, nil
}
