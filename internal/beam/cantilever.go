package beam

import "beam-vmd/internal/model"

func solveCantilever(b model.Beam) reactionSet {
	s := model.SummarizeLoads(b)
	return reactionSet{
		support: model.SupportCantilever,
		s0:      -s.TotalDown,
		xA:      -s.MomentAboutOrigin,
	}
}
