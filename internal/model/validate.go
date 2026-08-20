package model

import (
	"fmt"
)

// Validate checks the beam for physical and structural consistency. It returns a
// descriptive error for any illegal input:
//
//   - L must be strictly positive.
//   - EI must be strictly positive.
//   - every point load must sit inside [0, L].
//   - every distributed segment must satisfy 0 <= From <= To <= L.
//   - the support type must be one of the known values.
func Validate(b Beam) error {
	if b.L <= 0 {
		return fmt.Errorf("span L must be > 0, got %g", b.L)
	}
	if b.EI <= 0 {
		return fmt.Errorf("flexural rigidity EI must be > 0, got %g", b.EI)
	}
	if _, err := ResolveSupport(b.Support); err != nil {
		return err
	}
	for i, p := range b.Points {
		if p.At < 0 || p.At > b.L {
			return fmt.Errorf("point load %d at x=%g is outside the span [0, %g]", i, p.At, b.L)
		}
	}
	for i, d := range b.Distrib {
		if d.From < 0 || d.To > b.L {
			return fmt.Errorf("distributed load %d segment [%g, %g] is outside the span [0, %g]", i, d.From, d.To, b.L)
		}
		if d.To < d.From {
			return fmt.Errorf("distributed load %d has From=%g > To=%g", i, d.From, d.To)
		}
	}
	return nil
}
