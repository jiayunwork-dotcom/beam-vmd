package sample

import "beam-vmd/internal/model"

// ExampleSimplyUDL is a simply supported beam of span L=8 with EI=1600 carrying
// a full-span uniform load q=20.
//
//	midspan moment   M = q L^2 / 8      = 20*64/8   = 160
//	midspan deflection y = 5 q L^4 / (384 EI) = 5*20*4096/(384*1600) = 0.66667
func ExampleSimplyUDL() model.Beam {
	return model.Beam{
		L:       8,
		EI:      1600,
		Support: string(model.SupportSimply),
		Distrib: []model.DistributedLoad{{From: 0, To: 8, Q: 20}},
	}
}

// ExampleCantileverUDL is a cantilever fixed at x=0 with span L=5, EI=2500 and a
// full-span uniform load q=30.
//
//	fixed-end moment   M_fixed = -q L^2 / 2 = -30*25/2 = -375 (hogging)
//	free-end deflection y = q L^4 / (8 EI)  = 30*625/(8*2500) = 0.9375
func ExampleCantileverUDL() model.Beam {
	return model.Beam{
		L:       5,
		EI:      2500,
		Support: string(model.SupportCantilever),
		Distrib: []model.DistributedLoad{{From: 0, To: 5, Q: 30}},
	}
}

// ExampleFixedUDL is a both-ends-fixed beam of span L=6, EI=3000 under a
// full-span uniform load q=40.
//
//	fixed-end moment   M_end  = -q L^2 / 12 = -40*36/12 = -120
//	midspan moment     M_mid  =  q L^2 / 24 =  40*36/24 =  60
//	midspan deflection y = q L^4 / (384 EI) = 40*1296/(384*3000) = 0.045
func ExampleFixedUDL() model.Beam {
	return model.Beam{
		L:       6,
		EI:      3000,
		Support: string(model.SupportFixed),
		Distrib: []model.DistributedLoad{{From: 0, To: 6, Q: 40}},
	}
}

// ExampleSimplyOffset is a simply supported beam of span L=8, EI=2000 with a
// point load P=80 at x=2, offset from midspan.
//
//	reactions          r0 = P (L-a)/L = 60,  rL = P a/L = 20
//	moment at load     M  = r0 a = 120
//	deflection at load y  = P a^2 b^2 / (3 E I L) = 80*4*36/(3*2000*8) = 0.24
func ExampleSimplyOffset() model.Beam {
	return model.Beam{
		L:       8,
		EI:      2000,
		Support: string(model.SupportSimply),
		Points:  []model.PointLoad{{At: 2, Force: 80}},
	}
}

// ExampleSegmentedDistrib is a simply supported beam of span L=10, EI=10000 with
// two contiguous distributed segments: q=8 on [0,4] and q=12 on [4,10].
//
//	total load    = 8*4 + 12*6 = 104
//	moment about x=0 = 8*4*2 + 12*6*7 = 568
//	right reaction rL = 568/10 = 56.8 (upward), left reaction = 104 - 56.8 = 47.2
func ExampleSegmentedDistrib() model.Beam {
	return model.Beam{
		L:       10,
		EI:      10000,
		Support: string(model.SupportSimply),
		Distrib: []model.DistributedLoad{
			{From: 0, To: 4, Q: 8},
			{From: 4, To: 10, Q: 12},
		},
	}
}

// ExampleSimplyTwoPointLoads is a simply supported beam of span L=12, EI=12000
// carrying two point loads: P=60 at x=4 and P=90 at x=8.
//
//	reactions    r0 = (60*8 + 90*4)/12 = 70,  rL = 60 + 90 - 70 = 80
//	moments      M(4) = 70*4 = 280,  M(8) = 70*8 - 60*4 = 320
func ExampleSimplyTwoPointLoads() model.Beam {
	return model.Beam{
		L:       12,
		EI:      12000,
		Support: string(model.SupportSimply),
		Points: []model.PointLoad{
			{At: 4, Force: 60},
			{At: 8, Force: 90},
		},
	}
}

// ExampleFixedOffsetLoad is a both-ends-fixed beam of span L=8, EI=2000 with a
// point load P=100 at x=2, offset from midspan.
//
//	fixed-end moments mA = -P a b^2 / L^2 = -112.5, mB = -P a^2 b / L^2 = -37.5
//	moment at load   M  = 2 P a^2 b^2 / L^3 = 56.25
//	deflection at load y = P a^3 b^3 / (3 E I L^3) = 0.05625
func ExampleFixedOffsetLoad() model.Beam {
	return model.Beam{
		L:       8,
		EI:      2000,
		Support: string(model.SupportFixed),
		Points:  []model.PointLoad{{At: 2, Force: 100}},
	}
}

// ExampleCantileverOffsetLoad is a cantilever fixed at x=0 with span L=5, EI=1800
// and a point load P=40 at x=3, offset from the free end.
//
//	fixed-end moment   M_fixed = -P a = -120 (hogging)
//	free-end deflection y = P a^2 (3L - a) / (6 EI) = 40*9*12/10800 = 0.4
func ExampleCantileverOffsetLoad() model.Beam {
	return model.Beam{
		L:       5,
		EI:      1800,
		Support: string(model.SupportCantilever),
		Points:  []model.PointLoad{{At: 3, Force: 40}},
	}
}
