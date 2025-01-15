package scoring

import "sort"

type Player struct {
	Name        string    `json:"name"`
	Scores      []float64 `json:"score"`
	TotalScore  float64   `json:"totalScore"`
	Tournaments []string  `json:"tournaments"`
}

func (p *Player) ComputeTotalScore(diminishingRate, minDiminishedRate float64) {
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
