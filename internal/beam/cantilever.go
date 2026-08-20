package beam

import "beam-vmd/internal/model"

// solveCantilever returns the vertical reaction and fixing moment of a cantilever
// fixed at x=0 and free at x=L. The free end carries no reaction. With the
// downward-positive load summary (P total, M0 about the origin):
//
//	ΣF_y : s0 + P = 0   -> s0 = -P
//	ΣM_0 : xA + M0 = 0  -> xA = -M0
//
// so an upward reaction and a hogging (negative, sagging-positive) end moment
// balance the downward loading.
func solveCantilever(b model.Beam) reactionSet {
	s := model.SummarizeLoads(b)
	return reactionSet{
		support: model.SupportCantilever,
		s0:      -s.TotalDown,
		xA:      -s.MomentAboutOrigin,
	}
}
