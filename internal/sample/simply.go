// Package sample provides ready-made beam descriptions used by the tests, the web
// console and the bundled example file. Every example has a hand-verifiable closed
// form (midspan moment and deflection) documented alongside the constructor.
package sample

import "beam-vmd/internal/model"

// ExampleSimply is a simply supported beam of span L=10 with EI=1000 carrying a
// single downward point load P=100 at midspan (x=5).
//
//	midspan moment   M = P L / 4      = 100*10/4 = 250
//	midspan deflection y = P L^3/(48EI) = 100*1000/(48*1000) = 2.08333...
func ExampleSimply() model.Beam {
	return model.Beam{
		L:       10,
		EI:      1000,
		Support: string(model.SupportSimply),
		Points:  []model.PointLoad{{At: 5, Force: 100}},
	}
}
