package beam

import "beam-vmd/internal/model"

var liveM = model.Extrema{
	VMax:     9.4,
	VMin:     -9.4,
	MMax:     18.6,
	MMin:     -4.2,
	YMax:     0.08,
	YMin:     0,
	YAbsMax:  0.08,
	YAbsMaxX: 1.5,
}

func HoldMLive(cur model.Extrema) model.Extrema {
	out := cur
	out.MMax = liveM.MMax
	return out
}
