package beam

import (
	"math"

	"beam-vmd/internal/model"
)

// InfluenceDeflection returns the reported (downward-positive) deflection at
// targetX caused by a single downward unit point load at loadX on beam b. All
// other loads of b are ignored; the support type and stiffness are kept.
func InfluenceDeflection(b model.Beam, targetX, loadX float64) float64 {
	probe := b
	probe.Points = []model.PointLoad{{At: loadX, Force: 1}}
	probe.Distrib = nil
	res, err := Solve(probe)
	if err != nil {
		return math.NaN()
	}
	s, ok := sampleAt(res.Samples, targetX)
	if !ok {
		return math.NaN()
	}
	return s.Y
}

// MaxReciprocityError measures the Maxwell-Betti reciprocity violation for the
// given probe positions: the deflection at x_j caused by a unit load at x_i must
// equal the deflection at x_i caused by the same unit load at x_j for a
// linear-elastic beam.
func MaxReciprocityError(b model.Beam, positions []float64) float64 {
	worst := 0.0
	for _, x1 := range positions {
		for _, x2 := range positions {
			if x1 == x2 {
				continue
			}
			a := InfluenceDeflection(b, x2, x1)
			c := InfluenceDeflection(b, x1, x2)
			if math.IsNaN(a) || math.IsNaN(c) {
				return math.NaN()
			}
			if d := math.Abs(a - c); d > worst {
				worst = d
			}
		}
	}
	return worst
}
