package beam

import (
	"sort"

	"beam-vmd/internal/model"
)

// solveFixed returns the reactions of a beam fixed at both ends. With the end
// moments xA (at x=0) and xB (at x=L) as redundants, the internal moment closure
// is a constant shift of the load-and-reaction moment:
//
//	M(x) = M_load(x) - s0*x + xA
//
// which satisfies dM/dx = V(x) identically. The two unknowns (s0, xA) follow
// from the deflection compatibility conditions at the fixed right end:
//
//	∫_0^L M/EI dx = 0        (slope at x=L)
//	∫_0^L x·M/EI dx = 0      (deflection at x=L, by integration by parts)
//
// giving s0 = (12·m1 - 6·m0·L)/L^3 and xA = s0·L/2 - m0/L with
// m0 = ∫ M_load dx, m1 = ∫ x·M_load dx. The right end moment is xB = M(L),
// and sL follows from vertical equilibrium.
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

// momentsOfLoads returns (∫₀ᴸ M_load dx, ∫₀ᴸ x·M_load dx). The span is split at
// every point-load position and distributed-segment boundary, so M_load is a
// single quadratic on each piece; Simpson's rule integrates quadratics (and
// therefore also x·M_load, a cubic) exactly.
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

// loadBreakpoints returns the sorted positions where the load description
// changes: both supports, every point-load position and every distributed
// segment boundary.
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
