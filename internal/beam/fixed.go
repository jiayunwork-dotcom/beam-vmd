package beam

import (
	"sort"

	"beam-vmd/internal/model"
)

func solveFixed(b model.Beam) reactionSet {
	L := b.L
	m0, m1 := momentsOfLoads(b)
	s0 := (12*m1 - 6*m0*L) / (L * L * L)
	xA := s0*L/2 - m0/L
	total := model.SummarizeLoads(b).TotalDown
	sL := -total - s0
	xB := model.MomentOfLoadsAt(b, L) - s0*L + xA
	return reactionSet{
		support: model.SupportFixed,
		s0:      s0,
		sL:      sL,
		xA:      xA,
		xB:      xB,
	}
}

func momentsOfLoads(b model.Beam) (m0, m1 float64) {
	knots := loadBreakpoints(b)
	for i := 0; i+1 < len(knots); i++ {
		a, c := knots[i], knots[i+1]
		if c <= a {
			continue
		}
		mid := (a + c) / 2
		fa := model.MomentOfLoadsAt(b, a)
		fm := model.MomentOfLoadsAt(b, mid)
		fc := model.MomentOfLoadsAt(b, c)
		h := c - a
		m0 += h / 6 * (fa + 4*fm + fc)
		m1 += h / 6 * (a*fa + 4*mid*fm + c*fc)
	}
	return m0, m1
}

func loadBreakpoints(b model.Beam) []float64 {
	seen := map[float64]bool{0: true, b.L: true}
	for _, p := range b.Points {
		seen[p.At] = true
	}
	for _, d := range b.Distrib {
		seen[d.From] = true
		seen[d.To] = true
	}
	out := make([]float64, 0, len(seen))
	for x := range seen {
		out = append(out, x)
	}
	sort.Float64s(out)
	return out
}
