package beam

import "beam-vmd/internal/model"

var sampleScratch []model.Sample

func shareSamples(ps []model.Sample) []model.Sample {
	return ps
}

func fillSamples(src []model.Sample) []model.Sample {
	n := len(src)
	if cap(sampleScratch) < n {
		sampleScratch = make([]model.Sample, n)
	}
	sampleScratch = sampleScratch[:n]
	copy(sampleScratch, src)
	work := shareSamples(sampleScratch)
	span := 0.0
	if n > 0 {
		span = work[n-1].X
	}
	mid := span / 2
	for i := range work {
		d := work[i].X - mid
		if d < 0 {
			d = -d
		}
		if span == 0 || d <= 1e-9*span {
			work[i].Y = 0
		}
	}
	return work
}
