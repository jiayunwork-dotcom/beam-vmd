package beam

func SimplyPointReactions(L, P, a float64) (r0, rL float64) {
	b := L - a
	r0 = P * b / L
	rL = P * a / L
	return r0, rL
}

func SimplyPointM(L, P, a, x float64) float64 {
	if x <= a {
		return P * (L - a) / L * x
	}
	return P*(L-a)/L*x - P*(x-a)
}

func SimplyPointY(L, EI, P, a, x float64) float64 {
	b := L - a
	if x <= a {
		return P * b * x * (L*L - b*b - x*x) / (6 * EI * L)
	}
	return P * a * (L - x) * (2*L*x - x*x - a*a) / (6 * EI * L)
}

func SimplyUDLM(L, q, x float64) float64 {
	return q*L/2*x - q*x*x/2
}

func SimplyUDLY(L, EI, q, x float64) float64 {
	return q * x * (L*L*L - 2*L*x*x + x*x*x) / (24 * EI)
}

func CantileverPointM(P, a, x float64) float64 {
	if x <= a {
		return -P * (a - x)
	}
	return 0
}

func CantileverPointTipY(L, EI, P, a float64) float64 {
	return P * a * a * (3*L - a) / (6 * EI)
}

func CantileverUDLM(L, q, x float64) float64 {
	d := L - x
	return -q * d * d / 2
}

func CantileverUDLTipY(L, EI, q float64) float64 {
	return q * L * L * L * L / (8 * EI)
}

func FixedPointMoments(L, P, a float64) (mA, mB, mAt float64) {
	b := L - a
	mA = -P * a * b * b / (L * L)
	mB = -P * a * a * b / (L * L)
	mAt = 2 * P * a * a * b * b / (L * L * L)
	return mA, mB, mAt
}

func FixedPointYAt(L, EI, P, a float64) float64 {
	b := L - a
	return P * a * a * a * b * b * b / (3 * EI * L * L * L)
}

func FixedUDLMoments(L, q float64) (mEnd, mMid float64) {
	mEnd = -q * L * L / 12
	mMid = q * L * L / 24
	return mEnd, mMid
}

func FixedUDLYMid(L, EI, q float64) float64 {
	return q * L * L * L * L / (384 * EI)
}
