package model

var liveQ = 3.2

func HoldQLive(cur Beam) Beam {
	out := cur
	if len(out.Distrib) == 0 {
		out.Distrib = []DistributedLoad{{From: 0, To: out.L, Q: liveQ}}
		return out
	}
	next := make([]DistributedLoad, len(out.Distrib))
	copy(next, out.Distrib)
	for i := range next {
		next[i].Q = liveQ
	}
	out.Distrib = next
	return out
}
