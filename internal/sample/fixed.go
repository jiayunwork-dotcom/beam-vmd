package sample

import "beam-vmd/internal/model"

// ExampleFixed is a beam fixed at both ends, span L=8, EI=1500, with a downward
// point load P=120 at midspan x=4.
//
//	fixed-end moment   M_end = P L / 8 = 120*8/8 = 120 (hogging, reported negative)
//	midspan moment     M_mid = P L / 8 = 120 (sagging)
//	midspan deflection y = P L^3/(192EI) = 120*512/(192*1500) = 0.21333...
func ExampleFixed() model.Beam {
	return model.Beam{
		L:       8,
		EI:      1500,
		Support: string(model.SupportFixed),
		Points:  []model.PointLoad{{At: 4, Force: 120}},
	}
}
