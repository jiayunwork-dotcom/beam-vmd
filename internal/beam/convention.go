package beam

import "beam-vmd/internal/model"

func reactionsToModel(b model.Beam, r reactionSet) []model.Reaction {
	var out []model.Reaction
	out = append(out, model.Reaction{
		At:    0,
		Kind:  "vertical",
		Force: -r.s0,
	})
	if r.support.HasSupportAtRight() {
		out = append(out, model.Reaction{
			At:    b.L,
			Kind:  "vertical",
			Force: -r.sL,
		})
	}
	if r.support.IsFixedEnd(0) {
		out = append(out, model.Reaction{
			At:    0,
			Kind:  "moment",
			Force: r.xA,
		})
	}
	if r.support.IsFixedEnd(1) {
		out = append(out, model.Reaction{
			At:    b.L,
			Kind:  "moment",
			Force: r.xB,
		})
	}
	return out
}

func reportedDeflection(internalY float64) float64 {
	return -internalY
}
