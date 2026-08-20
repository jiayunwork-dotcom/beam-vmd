package model

import "sort"

// LoadPiece is one load-invariant interval of the beam: within it no point load
// and no distributed-segment boundary occurs, so the distributed intensity is
// constant. Point loads act exactly at the piece boundaries.
type LoadPiece struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
	Q    float64 `json:"q"` // constant intensity (downward-positive) on the piece
}

// LoadPieces returns the load-invariant intervals covering [0, L] in ascending
// order. Intervals outside any distributed segment carry Q = 0.
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

// IntensityAt returns the distributed intensity at position x, or zero when x
// lies outside every distributed segment.
func IntensityAt(b Beam, x float64) float64 {
	for _, d := range b.Distrib {
		if x >= d.From && x <= d.To {
			return d.Q
		}
	}
	return 0
}

// Resultant returns the resultant downward load of the piece and its centroid.
func (p LoadPiece) Resultant() (r, centroid float64) {
	r = p.Q * (p.To - p.From)
	centroid = (p.From + p.To) / 2
	return r, centroid
}
