package beam

// This file collects the analytic closed-form solutions used to cross-check the
// numerical solver. All deflection formulas are reported downward-positive (the
// same convention as the API result); all moments are sagging-positive.

// SimplyPointReactions returns the upward reactions of a simply supported beam
// carrying one point load P at position a (P downward-positive).
func SimplyPointReactions(L, P, a float64) (r0, rL float64) {
	b := L - a
	r0 = P * b / L
	rL = P * a / L
	return r0, rL
}

// SimplyPointM returns the sagging-positive bending moment at x of a simply
// supported beam with one point load P at a.
func SimplyPointM(L, P, a, x float64) float64 {
	if x <= a {
		return P * (L - a) / L * x
	}
	return P*(L-a)/L*x - P*(x-a)
}

// SimplyPointY returns the downward-positive deflection at x of a simply
// supported beam with one point load P at a. The two-branch cubic is the exact
// integration of EI y'' = M with y(0) = y(L) = 0.
func SimplyPointY(L, EI, P, a, x float64) float64 {
	b := L - a
	if x <= a {
		return P * b * x * (L*L - b*b - x*x) / (6 * EI * L)
	}
	return P * a * (L - x) * (2*L*x - x*x - a*a) / (6 * EI * L)
}

// SimplyUDLM returns the sagging-positive moment of a simply supported beam under
// a full-span uniform load q (downward-positive).
func SimplyUDLM(L, q, x float64) float64 {
	return q*L/2*x - q*x*x/2
}

// SimplyUDLY returns the downward-positive deflection of a simply supported beam
// under a full-span uniform load q.
func SimplyUDLY(L, EI, q, x float64) float64 {
	return q * x * (L*L*L - 2*L*x*x + x*x*x) / (24 * EI)
}

// CantileverPointM returns the sagging-positive moment of a cantilever fixed at
// x=0 under one point load P at distance a from the fixed end. Hogging regions
// come out negative.
func CantileverPointM(P, a, x float64) float64 {
	if x <= a {
		return -P * (a - x)
	}
	return 0
}

// CantileverPointTipY returns the downward-positive free-end deflection of a
// cantilever with one point load P at distance a from the fixed end.
func CantileverPointTipY(L, EI, P, a float64) float64 {
	return P * a * a * (3*L - a) / (6 * EI)
}

// CantileverUDLM returns the sagging-positive moment of a cantilever under a
// full-span uniform load q.
func CantileverUDLM(L, q, x float64) float64 {
	d := L - x
	return -q * d * d / 2
}

// CantileverUDLTipY returns the downward-positive free-end deflection of a
// cantilever under a full-span uniform load q.
func CantileverUDLTipY(L, EI, q float64) float64 {
	return q * L * L * L * L / (8 * EI)
}

// FixedPointMoments returns the fixed-end moments and the sagging moment at the
// load position for a both-ends-fixed beam with one point load P at a. The end
// moments are hogging and therefore negative in the sagging-positive convention.
func FixedPointMoments(L, P, a float64) (mA, mB, mAt float64) {
	b := L - a
	mA = -P * a * b * b / (L * L)
	mB = -P * a * a * b / (L * L)
	mAt = 2 * P * a * a * b * b / (L * L * L)
	return mA, mB, mAt
}

// FixedPointYAt returns the downward-positive deflection at the load position of
// a both-ends-fixed beam with one point load P at a.
func FixedPointYAt(L, EI, P, a float64) float64 {
	b := L - a
	return P * a * a * a * b * b * b / (3 * EI * L * L * L)
}

// FixedUDLMoments returns the fixed-end moments and the midspan moment of a
// both-ends-fixed beam under a full-span uniform load q.
func FixedUDLMoments(L, q float64) (mEnd, mMid float64) {
	mEnd = -q * L * L / 12
	mMid = q * L * L / 24
	return mEnd, mMid
}

// FixedUDLYMid returns the downward-positive midspan deflection of a
// both-ends-fixed beam under a full-span uniform load q.
func FixedUDLYMid(L, EI, q float64) float64 {
	return q * L * L * L * L / (384 * EI)
}
