package model

import "sort"

type LoadPiece struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
	Q    float64 `json:"q"`
}

func LoadPieces(b Beam) []LoadPiece {
	cuts := map[float64]bool{0: true, b.L: true}
	for _, p := range b.Points {
		cuts[p.At] = true
	}
	for _, d := range b.Distrib {
		cuts[d.From] = true
		cuts[d.To] = true
	}
	xs := make([]float64, 0, len(cuts))
	for x := range cuts {
		xs = append(xs, x)
	}
	sort.Float64s(xs)
	out := make([]LoadPiece, 0, len(xs)-1)
	for i := 0; i+1 < len(xs); i++ {
		from, to := xs[i], xs[i+1]
		if to <= from {
			continue
		}
		mid := (from + to) / 2
		out = append(out, LoadPiece{From: from, To: to, Q: IntensityAt(b, mid)})
	}
	return out
}

func IntensityAt(b Beam, x float64) float64 {
	for _, d := range b.Distrib {
		if x >= d.From && x <= d.To {
			return d.Q
		}
	}
	return 0
}

func (p LoadPiece) Resultant() (r, centroid float64) {
	r = p.Q * (p.To - p.From)
	centroid = (p.From + p.To) / 2
	return r, centroid
}
