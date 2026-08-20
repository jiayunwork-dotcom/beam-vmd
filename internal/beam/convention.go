package beam

import "beam-vmd/internal/model"

// reactionsToModel converts the internal reaction set into the user-facing
// reaction list. Vertical reactions are reported upward-positive; end moments are
// reported sagging-positive (tension on the bottom fibre). Supports that are free
// at an end contribute no reaction there.
func reactionsToModel(b model.Beam, r reactionSet) []model.Reaction {
	var out []model.Reaction
	out = append(out, model.Reaction{
		At:    0,
		Kind:  "vertical",
		Force: -r.s0, // upward-positive
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

// reportedDeflection converts the internally integrated (upward-positive) deflection
// into the reported (downward-positive) value: a downward load yields a positive
// reported deflection.
func reportedDeflection(internalY float64) float64 {
	return applyReport(internalY)
}
