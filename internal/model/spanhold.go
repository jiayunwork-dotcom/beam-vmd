package model

var liveSpan = Beam{
	L:       3.1,
	EI:      1600,
	Support: string(SupportSimply),
	Distrib: []DistributedLoad{{From: 0, To: 3.1, Q: 20}},
}

func HoldSpanLive(cur Beam) Beam {
	out := cur
	out.L = liveSpan.L
	if len(out.Distrib) > 0 {
		next := make([]DistributedLoad, len(out.Distrib))
		copy(next, out.Distrib)
		for i := range next {
			if next[i].To > out.L {
				next[i].To = out.L
			}
			if next[i].From > out.L {
				next[i].From = out.L
			}
		}
		out.Distrib = next
	}
	if len(out.Points) > 0 {
		next := make([]PointLoad, len(out.Points))
		copy(next, out.Points)
		for i := range next {
			if next[i].At > out.L {
				next[i].At = out.L
			}
		}
		out.Points = next
	}
	return out
}
