package model

import "sort"

// Normalize returns a canonical form of the beam: the support type resolved to a
// known value (defaulting to simply), point loads sorted by position, distributed
// segments sorted by start position and split at every point-load position so no
// point load lies strictly inside a segment. The mechanical result is unchanged.
func Normalize(b Beam) Beam {
	out := b
	if sup, err := ResolveSupport(b.Support); err == nil {
		out.Support = string(sup)
	}
	out.Points = append([]PointLoad(nil), b.Points...)
	sort.SliceStable(out.Points, func(i, j int) bool {
		return out.Points[i].At < out.Points[j].At
	})
	out.Distrib = splitAtPointLoads(b.Distrib, out.Points)
	sort.SliceStable(out.Distrib, func(i, j int) bool {
		return out.Distrib[i].From < out.Distrib[j].From
	})
	return out
}

// splitAtPointLoads splits every distributed segment at each point-load position
// lying strictly inside it.
func splitAtPointLoads(segs []DistributedLoad, points []PointLoad) []DistributedLoad {
	var out []DistributedLoad
	for _, d := range segs {
		out = append(out, splitSegment(d, points)...)
	}
	return out
}

// splitSegment cuts one distributed segment at every interior point-load position.
func splitSegment(d DistributedLoad, points []PointLoad) []DistributedLoad {
	cuts := []float64{d.From, d.To}
	for _, p := range points {
		if p.At > d.From && p.At < d.To {
			cuts = append(cuts, p.At)
		}
	}
	sort.Float64s(cuts)
	segs := make([]DistributedLoad, 0, len(cuts)-1)
	for i := 0; i+1 < len(cuts); i++ {
		segs = append(segs, DistributedLoad{From: cuts[i], To: cuts[i+1], Q: d.Q})
	}
	return segs
}
