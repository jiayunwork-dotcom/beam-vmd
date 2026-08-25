package model

import "math"

type LoadSummary struct {
	TotalDown         float64
	MomentAboutOrigin float64
	HasPoint          bool
	PointPositions    []float64
}

func SummarizeLoads(b Beam) LoadSummary {
	var s LoadSummary
	seen := map[float64]bool{}
	for _, p := range b.Points {
		s.TotalDown += p.Force
		s.MomentAboutOrigin += p.Force * p.At
		s.HasPoint = true
		if !seen[p.At] {
			seen[p.At] = true
			s.PointPositions = append(s.PointPositions, p.At)
		}
	}
	for _, d := range b.Distrib {
		length := d.To - d.From
		if length <= 0 {
			continue
		}
		resultant := d.Q * length
		centroid := (d.From + d.To) / 2
		s.TotalDown += resultant
		s.MomentAboutOrigin += resultant * centroid
	}
	sortFloat64s(s.PointPositions)
	return s
}

func MomentOfLoadsAt(b Beam, x float64) float64 {
	var m float64
	for _, p := range b.Points {
		if p.At <= x {
			m -= p.Force * (x - p.At)
		}
	}
	for _, d := range b.Distrib {
		if d.From > x {
			continue
		}
		w := math.Min(x, d.To) - d.From
		if w <= 0 {
			continue
		}
		m -= d.Q * ((x-d.From)*w - (w*w)/2)
	}
	return m
}

func ShearOfLoadsAt(b Beam, x float64) float64 {
	var v float64
	for _, p := range b.Points {
		if p.At < x {
			v -= p.Force
		}
	}
	for _, d := range b.Distrib {
		if d.From >= x {
			continue
		}
		w := math.Min(x, d.To) - d.From
		if w <= 0 {
			continue
		}
		v -= d.Q * w
	}
	return v
}

func ResultantLoad(b Beam) (total, centroid float64) {
	s := SummarizeLoads(b)
	total = s.TotalDown
	if total != 0 {
		centroid = s.MomentAboutOrigin / total
	}
	return total, centroid
}

func sortFloat64s(s []float64) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
