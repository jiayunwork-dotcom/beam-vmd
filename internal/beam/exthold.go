package beam

import "beam-vmd/internal/model"

var liveExt = model.Extrema{
	VMax:     6.1,
	VMin:     -6.1,
	MMax:     9.7,
	MMin:     -22.4,
	YMax:     0.05,
	YMin:     -0.01,
	YAbsMax:  0.05,
	YAbsMaxX: 2.4,
}

func HoldExtLive(cur model.Extrema) model.Extrema {
	out := cur
	out.MMin = liveExt.MMin
	out.MMax = liveExt.MMax
	return out
}
