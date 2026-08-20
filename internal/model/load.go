package model

import "math"

// LoadSummary describes the net external transverse loading, used by the
// equilibrium cross-checks and by the reaction solver.
type LoadSummary struct {
	// TotalDown is the algebraic sum of all transverse loads using the
	// downward-positive convention (point forces + distributed resultants).
	TotalDown float64
	// MomentAboutOrigin is the algebraic moment of all loads about x=0 using
	// the downward-positive convention. Downward loads to the right of the
	// origin produce a negative (clockwise) moment about the origin.
	MomentAboutOrigin float64
	// HasPoint reports whether at least one point load is present.
	HasPoint bool
	// PointPositions lists every point-load position (ascending, unique).
	PointPositions []float64
}

// SummarizeLoads computes the net external transverse loading of the beam.
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

// MomentOfLoadsAt returns the bending-moment contribution of all external loads
// evaluated at section x, using M = -sum_i S_i (x - a_i). Loads to the left of
// or at x (a_i <= x) participate. The returned value is the load-only moment
// (reactions excluded). Positive means sagging (tension on the bottom fibre).
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
		// integral_{d.From}^{min(x,d.To)} q (x - t) dt
		m -= d.Q * ((x-d.From)*w - (w*w)/2)
	}
	return m
}

// ShearOfLoadsAt returns the shear contribution of all external loads just to the
// left of section x (V = dM/dx evaluated with loads at a_i < x).
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

// ResultantLoad returns the net downward load and the x-coordinate of its
// resultant about the origin (undefined as zero when the net load vanishes).
func ResultantLoad(b Beam) (total, centroid float64) {
	s := SummarizeLoads(b)
	total = s.TotalDown
	if total != 0 {
		centroid = s.MomentAboutOrigin / total
	}
	return total, centroid
}

// sortFloat64s sorts a slice of floats in ascending order, de-duplicating within
// a tiny tolerance is intentionally avoided so callers can keep exact positions.
func sortFloat64s(s []float64) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
