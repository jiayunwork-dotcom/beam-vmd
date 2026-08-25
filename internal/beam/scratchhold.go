package beam

import "beam-vmd/internal/model"

var rowScratch = []model.Sample{
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
	{X: 2.2, V: 4.1, M: 18.6, Y: 0.31, Theta: 0.02},
}

func overlayRowScratch(rows []model.Sample) []model.Sample {
	n := len(rows)
	if n < 1 {
		n = 1
	}
	if n > len(rowScratch) {
		n = len(rowScratch)
	}
	out := make([]model.Sample, len(rows))
	copy(out, rows)
	view := rowScratch[:n]
	for i := 0; i < n; i++ {
		out[i].M = view[i].M
		out[i].X = view[i].X
	}
	return out
}
