package beam

import (
	"math"
	"sort"

	"beam-vmd/internal/model"
)

// defaultDivisions is the number of equal span subdivisions used when building the
// default sampling grid. It is well above the 20-division minimum required by the
// specification; load positions are always injected in addition.
const defaultDivisions = 24

// buildKnots returns the sorted, de-duplicated cross-section positions at which the
// solution is reported. It always includes both supports (0 and L), at least
// defaultDivisions equal subdivisions, and every point-load position plus every
// distributed-segment boundary so that moment/shear jumps are captured exactly.
func buildKnots(b model.Beam, divisions int) []float64 {
	if divisions < 20 {
		divisions = defaultDivisions
	}
	set := map[float64]bool{}
	set[0] = true
	set[b.L] = true
	for i := 1; i < divisions; i++ {
		set[b.L*float64(i)/float64(divisions)] = true
	}
	for _, p := range b.Points {
		set[p.At] = true
	}
	for _, d := range b.Distrib {
		set[d.From] = true
		set[d.To] = true
	}
	knots := make([]float64, 0, len(set))
	for x := range set {
		knots = append(knots, x)
	}
	sort.Float64s(knots)
	// collapse near-duplicates caused by floating point (e.g. L*i/n == a load pos)
	clean := knots[:0]
	for _, x := range knots {
		if len(clean) == 0 || math.Abs(x-clean[len(clean)-1]) > 1e-9*b.L {
			clean = append(clean, x)
		}
	}
	return clean
}

// sample evaluates V, M, deflection and slope at the requested knots. Shear is
// reported as the left-limit value at each knot (loads at exactly x are excluded),
// which is the conventional value on a shear diagram. Reported deflection is
// downward-positive.
func sampleSolution(b model.Beam, r reactionSet, it *integrator, knots []float64) []model.Sample {
	M := r.moment(b)
	V := r.shear(b)
	out := make([]model.Sample, 0, len(knots))
	for _, x := range knots {
		out = append(out, model.Sample{
			X:      x,
			V:      V(x),
			M:      M(x),
			Y:      reportedDeflection(it.y(x)),
			Theta:  it.theta(x),
		})
	}
	return out
}

// extrema scans the samples (which include the support and load positions) for the
// extreme shear, moment and reported deflection values.
func extrema(samples []model.Sample) model.Extrema {
	e := model.Extrema{
		VMax:    samples[0].V,
		VMin:    samples[0].V,
		MMax:    samples[0].M,
		MMin:    samples[0].M,
		YMax:    samples[0].Y,
		YMin:    samples[0].Y,
		YAbsMax: math.Abs(samples[0].Y),
		YAbsMaxX: samples[0].X,
	}
	for _, s := range samples {
		if s.V > e.VMax {
			e.VMax = s.V
		}
		if s.V < e.VMin {
			e.VMin = s.V
		}
		if s.M > e.MMax {
			e.MMax = s.M
		}
		if s.M < e.MMin {
			e.MMin = s.M
		}
		if s.Y > e.YMax {
			e.YMax = s.Y
		}
		if s.Y < e.YMin {
			e.YMin = s.Y
		}
		if a := math.Abs(s.Y); a > e.YAbsMax {
			e.YAbsMax = a
			e.YAbsMaxX = s.X
		}
	}
	return e
}
