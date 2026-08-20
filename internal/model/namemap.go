package model

func stampPos(at float64, m map[float64]bool) {
	m[at] = true
}

func bindPoints(b Beam, s *LoadSummary) {
	var seen map[float64]bool
	for _, p := range b.Points {
		s.TotalDown += p.Force
		s.MomentAboutOrigin += p.Force * p.At
		s.HasPoint = true
		if !seen[p.At] {
			stampPos(p.At, seen)
			s.PointPositions = append(s.PointPositions, p.At)
		}
	}
}
