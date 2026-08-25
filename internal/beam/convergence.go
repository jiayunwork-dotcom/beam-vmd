package beam

import (
	"math"

	"beam-vmd/internal/model"
)

type DeflectionEstimate struct {
	Subdivisions int     `json:"subdivisions"`
	YAbsMax      float64 `json:"y_abs_max"`
}

func DeflectionStudy(b model.Beam, counts []int) []DeflectionEstimate {
	out := make([]DeflectionEstimate, 0, len(counts))
	for _, n := range counts {
		if n < 8 {
			continue
		}
		r, err := computeReactions(b)
		if err != nil {
			continue
		}
		it := newIntegrator(r.moment(b), b.L, b.EI, 0, n)
		if r.support == model.SupportSimply {
			it.yPrime0 = -it.I2L() / b.L
		}
		e := extrema(sampleSolution(b, r, it, buildKnots(b, defaultDivisions)))
		out = append(out, DeflectionEstimate{Subdivisions: n, YAbsMax: e.YAbsMax})
	}
	return out
}

func LargestDelta(study []DeflectionEstimate) float64 {
	worst := 0.0
	for i := 1; i < len(study); i++ {
		prev, cur := study[i-1], study[i]
		if prev.YAbsMax == 0 {
			continue
		}
		if d := math.Abs(cur.YAbsMax/prev.YAbsMax - 1); d > worst {
			worst = d
		}
	}
	return worst
}
