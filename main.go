package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getSizeFactor(factors map[int]float64, s int) float64 {
	k := s / 8 * 8
	if k > 32 {
		k = 32
	}

	return factors[k]
}

func getRankingScore(seq []float64, r int) float64 {
	idx := len(seq) - r
	if idx < 0 {
		idx = 0
	}
	return seq[idx]
}

type Player struct {
	Name        string    `json:"name"`
	Scores      []float64 `json:"score"`
	TotalScore  float64   `json:"totalScore"`
	Tournaments []string  `json:"tournaments"`
}

func (p *Player) ComputeScore(dr float64) {
	sort.Sort(sort.Reverse(sort.Float64Slice(p.Scores)))

	var c float64 = 1
	for _, score := range p.Scores {
		p.TotalScore += score * c
		c -= dr
	}
}

func main() {
	cfg := initConfig()

	fs, err := os.Open(cfg.datasetFilePath)
	if err != nil {
		slog.Error("failed to open dataset", slog.Any("error", err))
	}
	defer fs.Close()

	var data map[string][]string
	err = json.NewDecoder(fs).Decode(&data)
	if err != nil {
		panic(err)
	}

	// Parse Rankings

	players := make(map[string]*Player)
	for name, rankings := range data {
		sizeFactor := getSizeFactor(cfg.sizeFactors, len(rankings))

		debug("%s - %d players (%.2f)\n", name, len(rankings), sizeFactor)

		for idx, pName := range rankings {
			if strings.HasSuffix(pName, cfg.nonFrenchPlayersSuffix) {
				continue
			}

			p, ok := players[pName]
			if !ok {
				p = &Player{pName, []float64{}, 0.0, []string{}}
				players[pName] = p
			}

			p.Scores = append(
				p.Scores,
				getRankingScore(cfg.scoreSequence, idx+1)*sizeFactor,
			)
			p.Tournaments = append(p.Tournaments, name)
		}
	}

	// Compute Scores
	rankings := make([]string, 0, len(players))
	for pName, p := range players {
		rankings = append(rankings, pName)
		p.ComputeScore(cfg.totalScoreDiminishingRate)
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return players[rankings[i]].TotalScore > players[rankings[j]].TotalScore
	})

	// Format console display
	w := new(tabwriter.Writer).Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t", "RANK", "PLAYER", "SCORE", "EVENTS")

	var rank int
	prevRank := -1
	prevScore := -1.0
	for idx, pName := range rankings {
		p := players[pName]

		rank = idx + 1
		if p.TotalScore == prevScore {
			rank = prevRank
		}
		prevRank = rank
		prevScore = p.TotalScore

		fmt.Fprintf(w, "\n %.2d\t%s\t%.2f\t%d/%d\t", rank, formatName(pName), p.TotalScore, len(p.Tournaments), len(data))
		if idx == 7 {
			fmt.Fprint(w, "\n --\t------------\t-----\t---\t")
		}
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}
}

func formatName(n string) string {
	caser := cases.Title(language.French)
	return caser.String(strings.ReplaceAll(n, ".", " "))
}

func debug(f string, a ...any) {
	fmt.Printf(f, a...)
}
