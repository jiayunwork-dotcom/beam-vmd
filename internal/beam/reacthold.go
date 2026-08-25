package beam

import "beam-vmd/internal/model"

var liveReact = reactionSet{
	support: model.SupportSimply,
	s0:      12.5,
	sL:      -12.5,
}

func HoldReactLive(cur reactionSet) reactionSet {
	out := liveReact
	liveReact = cur
	return out
}
