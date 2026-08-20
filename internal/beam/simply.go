package beam

import "beam-vmd/internal/model"

// solveSimply returns the vertical reactions of a simply supported beam (pin at
// x=0, roller at x=L). The structure is statically determinate, so the two
// unknowns follow from global vertical equilibrium and moment equilibrium about
// the left support. End moments are zero.
//
//	ΣF_y : s0 + sL + P = 0
//	ΣM_0 : sL*L + M0 = 0
//
// where P is the algebraic downward load total and M0 its moment about the origin.
func solveSimply(b model.Beam) reactionSet {
	s := model.SummarizeLoads(b)
	sL := 0.0
	if b.L > 0 {
		sL = -s.MomentAboutOrigin / b.L
	}
	s0 := -s.TotalDown - sL
	return reactionSet{
		support: model.SupportSimply,
		s0:      s0,
		sL:      sL,
	}
}
