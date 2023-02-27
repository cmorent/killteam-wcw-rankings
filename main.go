package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

var (
	avgAttendees    float64 = 0
	medianAttendees float64 = 0
)

func getTotalScore(pScore []float64) float64 {
	if totalScoreCountBestFirst {
		sort.Sort(sort.Reverse(sort.Float64Slice(pScore)))
	}

	if totalScoreMaxEvents > 0 && len(pScore) >= totalScoreMaxEvents {
		pScore = pScore[:totalScoreMaxEvents]
	}

	var c float64 = 1
	var res float64
	for _, s := range pScore {
		res += s * c
		c *= totalScoreDiminishingRate
	}
	return res
}

func getSizeFactor(s int) float64 {
	if !sizeFactorEnabled {
		return 1
	}

	d := avgAttendees
	if sizeFactorUseMedian {
		d = medianAttendees
	}

	return float64(s) * 1 / float64(d)
}

func getRankingScore(r int) float64 {
	f := rankingScoreSequences[rankingScoreSequencesIdx]
	idx := len(f) - r
	if idx < 0 {
		idx = 0
	}
	return float64(f[idx])
}

type Player struct {
	Name        string    `json:"name"`
	Scores      []float64 `json:"score"`
	TotalScore  float64   `json:"totalScore"`
	Tournaments []string  `json:"tournaments"`
}

func (p *Player) ComputeScore() {
	p.TotalScore = getTotalScore(p.Scores)
	if totalScoreMaxEvents > 0 && len(p.Tournaments) >= totalScoreMaxEvents {
		p.TotalScore += totalScoreMaxEventBonus
	}
}

func main() {
	fs, err := os.Open(datasetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	var tData map[string][]string
	err = json.NewDecoder(fs).Decode(&tData)
	if err != nil {
		log.Fatal(err)
	}

	// Compute both average and median number of attendees
	nbPlayers := 0
	var tAttendees []int
	for _, rankings := range tData {
		nbPlayers += len(rankings)
		tAttendees = append(tAttendees, len(rankings))
	}
	sort.Ints(tAttendees)
	avgAttendees = float64(nbPlayers) / float64(len(tData))
	medianAttendees = float64(tAttendees[len(tAttendees)/2-1])
	if len(tAttendees)%2 != 0 {
		medianAttendees = float64(tAttendees[(len(tAttendees)+1)/2-1])
	}
	dbg("Average number of players: %.2f\n", avgAttendees)
	dbg("Median number of players: %.2f\n\n", medianAttendees)

	players := make(map[string]*Player)

	// Parse Rankings
	for tName, rankings := range tData {
		sizeFactor := getSizeFactor(len(rankings))
		dbg("tournament %s - %d players (%.2f)\n", tName, len(rankings), sizeFactor)
		for idx, pName := range rankings {
			if strings.HasSuffix(pName, nonFrenchPlayersSuffix) {
				if nonFrenchPlayersIgnore {
					continue
				}
				pName = strings.TrimSuffix(pName, nonFrenchPlayersSuffix)
			}
			p, ok := players[pName]
			if !ok {
				p = &Player{pName, []float64{}, 0.0, []string{}}
				players[pName] = p
			}
			ranking := idx + 1
			p.Scores = append(p.Scores, getRankingScore(ranking)*sizeFactor)
			p.Tournaments = append(p.Tournaments, tName)
		}
	}
	dbg("")

	// Compute Scores
	rankings := make([]string, 0, len(players))
	for pName, p := range players {
		rankings = append(rankings, pName)
		p.ComputeScore()
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return players[rankings[i]].TotalScore > players[rankings[j]].TotalScore
	})

	w := new(tabwriter.Writer).Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", "RANK", "PLAYER", "SCORE", "EVENTS")
	for idx, pName := range rankings {
		p := players[pName]
		fmt.Fprintf(w, "\n %.2d\t%s\t%.2f\t%d/%d\t", idx+1, pName, p.TotalScore, len(p.Tournaments), totalScoreMaxEvents)
		if idx == 7 {
			if topEightOnly {
				break
			}
			fmt.Fprint(w, "\n --\t------------\t-----\t---\t")
		}
	}
	w.Write([]byte("\n"))
}

func dbg(f string, a ...any) {
	if debug {
		fmt.Printf(f, a...)
	}
}
