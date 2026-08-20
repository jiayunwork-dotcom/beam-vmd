package sample

import "beam-vmd/internal/model"

// ExampleCantilever is a cantilever fixed at x=0 with span L=6, EI=2000 and a
// downward tip load P=50 at the free end x=6.
//
//	fixed-end moment   M_fixed = P L      = 50*6   = 300 (hogging, reported negative)
//	free-end deflection y = P L^3/(3EI)   = 50*216/(3*2000) = 1.8
func ExampleCantilever() model.Beam {
	return model.Beam{
		L:       6,
		EI:      2000,
		Support: string(model.SupportCantilever),
		Points:  []model.PointLoad{{At: 6, Force: 50}},
	}
}
