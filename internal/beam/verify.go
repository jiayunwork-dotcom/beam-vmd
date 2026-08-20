package beam

import (
	"fmt"
	"math"

	"beam-vmd/internal/model"
)

// sampleAt returns the sample at position x within a tight relative tolerance,
// so a load position that coincides with a grid knot after de-duplication is
// still matched.
func sampleAt(samples []model.Sample, x float64) (model.Sample, bool) {
	var best model.Sample
	bestD := math.Inf(1)
	for _, s := range samples {
		if d := math.Abs(s.X - x); d < bestD {
			best, bestD = s, d
		}
	}
	if len(samples) < 2 {
		return model.Sample{}, false
	}
	span := samples[len(samples)-1].X - samples[0].X
	if bestD > 1e-6*span {
		return model.Sample{}, false
	}
	return best, true
}

// fullSpanUDL reports whether the model is exactly one full-span uniform load.
func fullSpanUDL(b model.Beam) (model.DistributedLoad, bool) {
	if len(b.Points) != 0 || len(b.Distrib) != 1 {
		return model.DistributedLoad{}, false
	}
	d := b.Distrib[0]
	if math.Abs(d.From) > 1e-12 || math.Abs(d.To-b.L) > 1e-12*b.L {
		return model.DistributedLoad{}, false
	}
	return d, true
}

// fixedPointMAt extracts the sagging moment at the load position.
func fixedPointMAt(L, P, a float64) float64 {
	_, _, mAt := FixedPointMoments(L, P, a)
	return mAt
}

// checkClosedForm matches the model against standard cases with known analytic
// solutions and compares the sampled moment and deflection to the closed form.
// This pins down EI y'' = M both in sign and in magnitude: dropping EI in the
// integration would put the deflection envelope at the wrong scale.
func checkClosedForm(b model.Beam, samples []model.Sample) model.CheckItem {
	sup, err := model.ResolveSupport(b.Support)
	if err != nil {
		return model.CheckItem{Name: "closed-form", Passed: true, Message: "no reference"}
	}
	if len(b.Points) == 1 && len(b.Distrib) == 0 {
		p := b.Points[0]
		switch sup {
		case model.SupportSimply:
			return compareClosedForm("simply point load", samples, p.At, SimplyPointM(b.L, p.Force, p.At, p.At), p.At, SimplyPointY(b.L, b.EI, p.Force, p.At, p.At))
		case model.SupportCantilever:
			return compareClosedForm("cantilever point load", samples, 0, CantileverPointM(p.Force, p.At, 0), b.L, CantileverPointTipY(b.L, b.EI, p.Force, p.At))
		case model.SupportFixed:
			return compareClosedForm("fixed point load", samples, p.At, fixedPointMAt(b.L, p.Force, p.At), p.At, FixedPointYAt(b.L, b.EI, p.Force, p.At))
		}
	}
	if d, ok := fullSpanUDL(b); ok {
		switch sup {
		case model.SupportSimply:
			return compareClosedForm("simply uniform", samples, b.L/2, SimplyUDLM(b.L, d.Q, b.L/2), b.L/2, SimplyUDLY(b.L, b.EI, d.Q, b.L/2))
		case model.SupportCantilever:
			return compareClosedForm("cantilever uniform", samples, 0, CantileverUDLM(b.L, d.Q, 0), b.L, CantileverUDLTipY(b.L, b.EI, d.Q))
		case model.SupportFixed:
			_, mMid := FixedUDLMoments(b.L, d.Q)
			return compareClosedForm("fixed uniform", samples, b.L/2, mMid, b.L/2, FixedUDLYMid(b.L, b.EI, d.Q))
		}
	}
	return model.CheckItem{Name: "closed-form", Passed: true, Message: "no standard reference"}
}

// compareClosedForm compares the sampled moment at xM and deflection at xY with
// the analytic values for the matched standard case.
func compareClosedForm(label string, samples []model.Sample, xM, wantM, xY, wantY float64) model.CheckItem {
	sm, okM := sampleAt(samples, xM)
	sy, okY := sampleAt(samples, xY)
	if !okM || !okY {
		return model.CheckItem{
			Name:    "closed-form",
			Passed:  false,
			Message: fmt.Sprintf("%s: missing samples at x=%g / x=%g", label, xM, xY),
		}
	}
	scaleM := math.Max(math.Abs(wantM), math.Abs(sm.M)) + 1
	scaleY := math.Max(math.Abs(wantY), math.Abs(sy.Y)) + 1
	ok := closeEnough(sm.M, wantM, scaleM) && closeEnough(sy.Y, wantY, scaleY)
	msg := fmt.Sprintf("%s: M(%g)=%g want %g; y(%g)=%g want %g",
		label, xM, sm.M, wantM, xY, sy.Y, wantY)
	if !ok {
		msg = "FAIL: " + msg
	} else {
		msg = "pass: " + msg
	}
	return model.CheckItem{Name: "closed-form", Passed: ok, Message: msg}
}
