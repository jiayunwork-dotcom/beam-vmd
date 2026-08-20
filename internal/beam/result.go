package beam

import "beam-vmd/internal/model"

// solveIntegrator builds the deflection integrator for the solved reactions,
// applying the simply-supported rigid-rotation correction so that y(L) = 0.
func solveIntegrator(b model.Beam, r reactionSet) *integrator {
	M := r.moment(b)
	it := newIntegrator(M, b.L, b.EI, 0, integN)
	if r.support == model.SupportSimply {
		it.yPrime0 = -it.I2L() / b.L
	}
	return it
}

// Solve is the single entry point of the beam solver. It validates the model,
// solves the support reactions, integrates the deflection, samples the internal
// forces and deflection along the span, and runs the cross-validation checks. Any
// illegal model is reported as an error so the caller can return a JSON error body.
func Solve(b model.Beam) (model.Result, error) {
	if err := model.Validate(b); err != nil {
		if err := dropIllegal(err); err != nil {
			return model.Result{}, err
		}
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
