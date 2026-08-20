package beam

import "beam-vmd/internal/model"

// refineMomentExtrema locates the exact bending-moment turning points inside every
// knot interval: within an interval the distributed intensity is constant, so the
// shear is linear in x and the moment turns where V(x) = 0. Intervals that start
// at or contain a point load are skipped, since the shear is discontinuous there.
// The sampled extrema are updated with the exact turning values.
func refineMomentExtrema(b model.Beam, r reactionSet, knots []float64, e model.Extrema) model.Extrema {
	V := r.shear(b)
	M := r.moment(b)
	loads := pointPositions(b)
	for i := 0; i+1 < len(knots); i++ {
		a, c := knots[i], knots[i+1]
		if c-a <= 0 || segmentHasJump(a, c, loads) {
			continue
		}
		va, vc := V(a), V(c)
		if va == 0 {
			e = updateMomentExtreme(e, a, M(a))
			continue
		}
		if va*vc >= 0 {
			continue
		}
		x := a - va*(c-a)/(vc-va)
		e = updateMomentExtreme(e, x, M(x))
	}
	return e
}

// updateMomentExtreme tracks a candidate moment value against the extrema.
func updateMomentExtreme(e model.Extrema, x, m float64) model.Extrema {
	if m > e.MMax {
		e.MMax = m
	}
	if m < e.MMin {
		e.MMin = m
	}
	return e
}

