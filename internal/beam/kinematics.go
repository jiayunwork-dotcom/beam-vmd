// Package beam implements Euler-Bernoulli beam mechanics: piecewise-equilibrium
// internal forces, statically-determined and indeterminate support reactions,
// double integration of EI y'' = M, sampling, and cross-validation checks.
//
// Sign conventions (see the project README; they mirror internal/model):
//
//   - Loads use a downward-positive convention; reactions are first solved in the
//     same downward-positive convention, then reported upward-positive.
//   - Bending moment M(x) is sagging-positive (tension on the bottom fibre).
//   - Deflection is integrated with EI y'' = M and y measured upward-positive; the
//     value reported to users is the negative of that (so a downward load yields a
//     positive reported deflection).
package beam

import "beam-vmd/internal/model"

// reactionSet holds the solved support reactions in internal sign conventions.
// Vertical fields are signed downward-positive (identical to loads); end moments
// are sagging-positive. yPrime0 is the slope at x=0 used when integrating the
// deflection (upward-positive).
type reactionSet struct {
	support model.SupportType
	s0, sL float64 // vertical reaction at x=0 and x=L (downward-positive)
	xA, xB float64 // end moment at x=0 and x=L (sagging-positive)
	yPrime0 float64 // slope at x=0 for deflection integration (upward-positive)
}

// moment returns the full bending-moment closure M(x) including the loads and the
// solved reactions. Fixed-end moments enter as a constant shift so that the
// closure satisfies dM/dx = V(x) identically.
func (r reactionSet) moment(b model.Beam) func(float64) float64 {
	return func(x float64) float64 {
		m := model.MomentOfLoadsAt(b, x) - r.s0*x
		switch r.support {
		case model.SupportCantilever, model.SupportFixed:
			m += r.xA
		}
		return m
	}
}

// shear returns the full shear-force closure V(x) including loads and the left
// vertical reaction. The right vertical reaction lies to the right of every
// interior section and therefore does not appear.
func (r reactionSet) shear(b model.Beam) func(float64) float64 {
	return func(x float64) float64 {
		return model.ShearOfLoadsAt(b, x) - r.s0
	}
}
