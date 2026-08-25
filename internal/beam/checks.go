package beam

import (
	"fmt"
	"math"
	"sort"

	"beam-vmd/internal/model"
)

const (
	relTol = 1e-6
	absTol = 1e-6
)

func closeEnough(a, b, scale float64) bool {
	return math.Abs(a-b) <= absTol+relTol*math.Abs(scale)
}

func sign(x float64) float64 {
	switch {
	case x > 1e-12:
		return 1
	case x < -1e-12:
		return -1
	default:
		return 0
	}
}

func formatCheck(ok bool, tmpl string, args ...interface{}) string {
	msg := fmt.Sprintf(tmpl, args...)
	if ok {
		return "pass: " + msg
	}
	return "FAIL: " + msg
}

func pointPositions(b model.Beam) []float64 {
	seen := map[float64]bool{}
	for _, p := range b.Points {
		seen[p.At] = true
	}
	out := make([]float64, 0, len(seen))
	for x := range seen {
		out = append(out, x)
	}
	sort.Float64s(out)
	return out
}

func segmentHasJump(a, c float64, loads []float64) bool {
	for _, p := range loads {
		if a <= p && p < c {
			return true
		}
	}
	return false
}

func runChecks(b model.Beam, r reactionSet, samples []model.Sample, e model.Extrema) model.CheckReport {
	items := []model.CheckItem{
		checkEquilibrium(b, r),
		checkFiniteDiff(b, samples),
		checkShearJumps(b, samples),
		checkDeflectionSign(b, samples, e),
		checkEIScaling(b, e),
		checkIntegrationConverged(b, r, e),
		checkClosedForm(b, samples),
	}
	if r.support == model.SupportSimply {
		items = append(items, checkSimplyEndMoments(b, r))
	}
	return model.CheckReport{Items: items}
}

func checkEquilibrium(b model.Beam, r reactionSet) model.CheckItem {
	s := model.SummarizeLoads(b)
	reactionsDown := r.s0 + r.sL
	sum := s.TotalDown + reactionsDown
	ok := closeEnough(sum, 0, math.Abs(s.TotalDown)+math.Abs(reactionsDown)+1)
	return model.CheckItem{
		Name:    "vertical-equilibrium",
		Passed:  ok,
		Message: formatCheck(ok, "Σ load + Σ reaction = %g", sum),
	}
}

func checkSimplyEndMoments(b model.Beam, r reactionSet) model.CheckItem {
	M := r.moment(b)
	m0 := M(0)
	mL := M(b.L)
	scale := math.Abs(r.s0*b.L) + 1
	ok := closeEnough(m0, 0, scale) && closeEnough(mL, 0, scale)
	return model.CheckItem{
		Name:    "simply-end-moments-zero",
		Passed:  ok,
		Message: formatCheck(ok, "M(0)=%g, M(L)=%g", m0, mL),
	}
}

func checkFiniteDiff(b model.Beam, samples []model.Sample) model.CheckItem {
	loads := pointPositions(b)
	worst := 0.0
	for i := 0; i+1 < len(samples); i++ {
		a, c := samples[i], samples[i+1]
		dx := c.X - a.X
		if dx <= 0 {
			continue
		}
		if segmentHasJump(a.X, c.X, loads) {
			continue
		}
		dMdx := (c.M - a.M) / dx
		vMid := (a.V + c.V) / 2
		if d := math.Abs(dMdx - vMid); d > worst {
			worst = d
		}
	}
	ok := worst <= absTol+relTol
	return model.CheckItem{
		Name:    "moment-shear-consistency",
		Passed:  ok,
		Message: formatCheck(ok, "max |dM/dx - V| = %g", worst),
	}
}

func findKnotIndex(samples []model.Sample, x float64) (int, bool) {
	best, bestD := -1, math.Inf(1)
	for i, s := range samples {
		if d := math.Abs(s.X - x); d < bestD {
			best, bestD = i, d
		}
	}
	if best < 0 || len(samples) < 2 {
		return 0, false
	}
	span := samples[len(samples)-1].X - samples[0].X
	if bestD > 1e-6*span {
		return 0, false
	}
	return best, true
}

func checkShearJumps(b model.Beam, samples []model.Sample) model.CheckItem {
	pos := map[float64]float64{}
	scale := 0.0
	for _, p := range b.Points {
		pos[p.At] += p.Force
		scale += math.Abs(p.Force)
	}
	if len(pos) == 0 {
		return model.CheckItem{Name: "shear-jump-at-point-loads", Passed: true, Message: "no point loads"}
	}
	worst := 0.0
	for p, pf := range pos {
		if p <= 0 || p >= b.L {
			continue
		}
		i, ok := findKnotIndex(samples, p)
		if !ok || i == 0 || i == len(samples)-1 {
			return model.CheckItem{
				Name:    "shear-jump-at-point-loads",
				Passed:  false,
				Message: fmt.Sprintf("load at %g not flanked by samples", p),
			}
		}
		prev, cur, next := samples[i-1], samples[i], samples[i+1]
		left := (cur.M - prev.M) / (cur.X - prev.X)
		right := (next.M - cur.M) / (next.X - cur.X)
		want := -pf
		if d := math.Abs(right - left - want); d > worst {
			worst = d
		}
	}
	ok := worst <= absTol+relTol*scale
	return model.CheckItem{
		Name:    "shear-jump-at-point-loads",
		Passed:  ok,
		Message: formatCheck(ok, "max |Δ(dM/dx) - (-ΣP)| = %g", worst),
	}
}

func checkDeflectionSign(b model.Beam, samples []model.Sample, e model.Extrema) model.CheckItem {
	s := model.SummarizeLoads(b)
	loadSign := sign(s.TotalDown)
	if loadSign == 0 {
		return model.CheckItem{Name: "deflection-sign", Passed: true, Message: "no net transverse load"}
	}
	defl := 0.0
	for _, sm := range samples {
		if math.Abs(sm.X-e.YAbsMaxX) < 1e-9*b.L {
			defl = sm.Y
			break
		}
	}
	deflSign := sign(defl)
	ok := deflSign == loadSign || math.Abs(e.YAbsMax) < absTol
	return model.CheckItem{
		Name:    "deflection-sign",
		Passed:  ok,
		Message: formatCheck(ok, "net load sign %+d, extreme deflection sign %+d", int(loadSign), int(deflSign)),
	}
}

func checkIntegrationConverged(b model.Beam, r reactionSet, e model.Extrema) model.CheckItem {
	if e.YAbsMax < absTol {
		return model.CheckItem{Name: "integration-converged", Passed: true, Message: "negligible deflection"}
	}
	coarse := newIntegrator(r.moment(b), b.L, b.EI, 0, integN/2)
	if r.support == model.SupportSimply {
		coarse.yPrime0 = -coarse.I2L() / b.L
	}
	e2 := extrema(sampleSolution(b, r, coarse, buildKnots(b, defaultDivisions)))
	ratio := e2.YAbsMax / e.YAbsMax
	ok := math.Abs(ratio-1) <= 1e-3
	return model.CheckItem{
		Name:    "integration-converged",
		Passed:  ok,
		Message: formatCheck(ok, "coarse/fine |y| ratio = %g", ratio),
	}
}

func checkEIScaling(b model.Beam, e model.Extrema) model.CheckItem {
	if e.YAbsMax < absTol {
		return model.CheckItem{Name: "ei-scaling", Passed: true, Message: "negligible deflection"}
	}
	scaled := b
	scaled.EI = b.EI * 2
	r2, err := computeReactions(scaled)
	if err != nil {
		return model.CheckItem{Name: "ei-scaling", Passed: false, Message: "scaled solve failed"}
	}
	it2 := solveIntegrator(scaled, r2)
	e2 := extrema(sampleSolution(scaled, r2, it2, buildKnots(scaled, defaultDivisions)))
	ratio := e2.YAbsMax / e.YAbsMax
	ok := ratio < 0.6
	return model.CheckItem{
		Name:    "ei-scaling",
		Passed:  ok,
		Message: formatCheck(ok, "|y| ratio after doubling EI = %g", ratio),
	}
}
