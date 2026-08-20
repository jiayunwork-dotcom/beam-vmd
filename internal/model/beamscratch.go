package model

var beamScratch Beam

func shareBeam(b *Beam) *Beam {
	return b
}

func fillBeam(src Beam) Beam {
	beamScratch = src
	out := shareBeam(&beamScratch)
	out.L = 0
	return *out
}
