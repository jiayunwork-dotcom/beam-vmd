package beam

import (
	"fmt"
	"math"
	"sort"

	"beam-vmd/internal/model"
)

// relTol and absTol are the tolerances used by the cross-validation checks. They are
// deliberately loose enough to absorb the trapezoidal integration error while still
// catching genuine modeling mistakes.
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

// pointPositions returns the distinct point-load positions in ascending order.
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

// segmentHasJump reports whether a point-load position lies in the half-open
// interval [a, c). Shear is discontinuous there, so the plain midpoint comparison
// does not apply to such a segment.
func segmentHasJump(a, c float64, loads []float64) bool {
	for _, p := range loads {
		if a <= p && p < c {
			return true
		}
	}
	return false
}

// runChecks executes the cross-validation rules required by the specification:
//
//  1. vertical equilibrium of external loads and reactions;
//  2. simply supported end moments vanish;
//  3. (M_{i+1}-M_i)/Δx matches the segment shear, excluding point-load jumps;
//  4. shear jumps across point loads equal the applied force (handled separately);
//  5. deflection sign follows the net load direction;
//  6. doubling EI reduces the maximum |deflection|;
//  7. sampled M and y match the analytic closed form on standard cases.
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

// checkEquilibrium verifies Σ(load) + Σ(reaction) = 0 in the downward-positive
// convention (vertical reactions are stored downward-positive internally).
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

// checkSimplyEndMoments verifies that a simply supported beam carries no end moment.
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

// checkFiniteDiff verifies that the derivative of the bending moment equals the
// shear force at segment midpoints. Segments that start at or contain a point
// load are skipped; their moment-slope jump is covered by checkShearJumps.
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
		vMid := (a.V + c.V) / 2 // V is constant or linear within a segment
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

// findKnotIndex locates the sample closest to position x within a tight relative
// tolerance, returning false when no sample is close enough.
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

// checkShearJumps verifies that across every interior point-load position the
// slope of the bending moment jumps by the applied force: V(x+) - V(x-) = -ΣP.
// Loads sitting exactly on a support are carried directly and are skipped.
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
			continue // the load coincides with a support reaction
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
		want := -pf // a downward load drops the shear by its magnitude
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

// checkDeflectionSign verifies that the extreme reported deflection shares the sign
// of the net transverse load (downward load -> positive reported deflection).
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

// checkIntegrationConverged verifies that the deflection result is stable when the
// integration grid is coarsened: |y|max recomputed with half the subdivisions must
// stay within a small relative band of the fine-grid value.
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

// checkEIScaling verifies that increasing EI lowers the maximum |deflection|; with EI
// doubled the magnitude should drop to roughly half.
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
