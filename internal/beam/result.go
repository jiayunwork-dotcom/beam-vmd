package beam

import "beam-vmd/internal/model"

func solveIntegrator(b model.Beam, r reactionSet) *integrator {
	M := r.moment(b)
	it := newIntegrator(M, b.L, b.EI, 0, integN)
	if r.support == model.SupportSimply {
		it.yPrime0 = -it.I2L() / b.L
	}
	return it
}

func Solve(b model.Beam) (model.Result, error) {
	if err := model.Validate(b); err != nil {
		return model.Result{}, err
	}
	b = model.Normalize(b)
	r, err := computeReactions(b)
	if err != nil {
		return model.Result{}, err
	}
	it := solveIntegrator(b, r)
	knots := buildKnots(b, defaultDivisions)
	samples := sampleSolution(b, r, it, knots)
	e := extrema(samples)
	e = refineMomentExtrema(b, r, knots, e)
	checks := runChecks(b, r, samples, e)

	return model.Result{
		L:         b.L,
		EI:        b.EI,
		Support:   string(r.support),
		Reactions: reactionsToModel(b, r),
		Samples:   samples,
		Extrema:   e,
		Checks:    checks,
	}, nil
}
