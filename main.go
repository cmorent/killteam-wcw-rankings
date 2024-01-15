package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getSizeFactor(s int) float64 {
	if !sizeFactorEnabled {
		return 1
	}
	k := s / 8 * 8
	if k > 32 {
		k = 32
	}

	return sizeFactors[k]
}

func getRankingScore(r int) float64 {
	f := rankingScoreSequences[rankingScoreSequencesIdx]
	idx := len(f) - r
	if idx < 0 {
		idx = 0
	}
	return f[idx]
}

type Player struct {
	Name        string    `json:"name"`
	Scores      []float64 `json:"score"`
	TotalScore  float64   `json:"totalScore"`
	Tournaments []string  `json:"tournaments"`
}

func (p *Player) ComputeScore() {
	sort.Sort(sort.Reverse(sort.Float64Slice(p.Scores)))

	var c float64 = 1
	for _, score := range p.Scores {
		p.TotalScore += score * c
		c -= totalScoreDiminishingRate
	}
}

func main() {
	fs, err := os.Open(datasetFilePath)
	if err != nil {
		log.Fatalf("failed to open dataset: %s", err)
	}
	defer fs.Close()

	var data map[string][]string
	err = json.NewDecoder(fs).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	// Parse Rankings

	players := make(map[string]*Player)
	for name, rankings := range data {
		sizeFactor := getSizeFactor(len(rankings))

		debug("tournament %s - %d players (%.2f)\n", name, len(rankings), sizeFactor)

		for idx, pName := range rankings {
			if strings.HasSuffix(pName, nonFrenchPlayersSuffix) {
				continue
			}

			p, ok := players[pName]
			if !ok {
				p = &Player{pName, []float64{}, 0.0, []string{}}
				players[pName] = p
			}

			p.Scores = append(p.Scores, getRankingScore(idx+1)*sizeFactor)
			p.Tournaments = append(p.Tournaments, name)
		}
	}
	debug("\n")

	// Compute Scores
	rankings := make([]string, 0, len(players))
	for pName, p := range players {
		rankings = append(rankings, pName)
		p.ComputeScore()
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return players[rankings[i]].TotalScore > players[rankings[j]].TotalScore
	})

	// ctx := context.Background()
	// sa := option.WithCredentialsFile("sa.json")

	// app, err := firebase.NewApp(ctx, nil, sa)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// cl, err := app.Firestore(ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// defer cl.Close()

	// Format console display
	w := new(tabwriter.Writer).Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", "RANK", "PLAYER", "SCORE", "EVENTS")

	rank, prevRank := 0, -1
	prevScore := -1.0
	for idx, pName := range rankings {
		p := players[pName]

		rank = idx + 1
		if p.TotalScore == prevScore {
			rank = prevRank
		}
		prevRank = rank
		prevScore = p.TotalScore

		// 	// Firestore
		// 	pts := strings.Split(pName, ".")
		// 	if len(pts) == 1 {
		// 		pts = append(pts, "")
		// 	}

		// 	iter := cl.Collection("rankings").Where("firstname", "==", pts[0]).Where("lastname", "==", pts[1]).Documents(ctx)
		// 	for {
		// 		doc, err := iter.Next()
		// 		if err == iterator.Done {
		// 			break
		// 		}

		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}

		// 		_, err = cl.Collection("rankings").Doc(doc.Ref.ID).Update(ctx, []firestore.Update{
		// 			{
		// 				Path:  "rank",
		// 				Value: rank,
		// 			},
		// 			{
		// 				Path:  "score",
		// 				Value: p.TotalScore,
		// 			},
		// 		})
		// 		if err != nil {
		// 			log.Fatalf("update error: %s", err)
		// 		}
		// 	}

		fmt.Fprintf(w, "\n %.2d\t%s\t%.2f\t%d/%d\t", rank, formatName(pName), p.TotalScore, len(p.Tournaments), len(data))
		if idx == 7 {
			fmt.Fprint(w, "\n --\t------------\t-----\t---\t")
		}
	}
	w.Write([]byte("\n"))
}

func formatName(n string) string {
	caser := cases.Title(language.French)
	return caser.String(strings.ReplaceAll(n, ".", " "))
}

func debug(f string, a ...any) {
	if debugEnabled {
		fmt.Printf(f, a...)
	}
}
