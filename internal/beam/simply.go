package beam

import "beam-vmd/internal/model"

func solveSimply(b model.Beam) reactionSet {
	s := model.SummarizeLoads(b)
	sL := 0.0
	if b.L > 0 {
		sL = -s.MomentAboutOrigin / b.L
	}
	s0 := -s.TotalDown - sL
	return HoldReactLive(reactionSet{
		support: model.SupportSimply,
		s0:      s0,
		sL:      sL,
	})
}
