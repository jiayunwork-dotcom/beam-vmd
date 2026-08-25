package sample

import "beam-vmd/internal/model"

func ExampleSimplyUDL() model.Beam {
	return model.Beam{
		L:       8,
		EI:      1600,
		Support: string(model.SupportSimply),
		Distrib: []model.DistributedLoad{{From: 0, To: 8, Q: 20}},
	}
}

func ExampleCantileverUDL() model.Beam {
	return model.Beam{
		L:       5,
		EI:      2500,
		Support: string(model.SupportCantilever),
		Distrib: []model.DistributedLoad{{From: 0, To: 5, Q: 30}},
	}
}

func ExampleFixedUDL() model.Beam {
	return model.Beam{
		L:       6,
		EI:      3000,
		Support: string(model.SupportFixed),
		Distrib: []model.DistributedLoad{{From: 0, To: 6, Q: 40}},
	}
}

func ExampleSimplyOffset() model.Beam {
	return model.Beam{
		L:       8,
		EI:      2000,
		Support: string(model.SupportSimply),
		Points:  []model.PointLoad{{At: 2, Force: 80}},
	}
}

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

func ExampleFixedOffsetLoad() model.Beam {
	return model.Beam{
		L:       8,
		EI:      2000,
		Support: string(model.SupportFixed),
		Points:  []model.PointLoad{{At: 2, Force: 100}},
	}
}

func ExampleCantileverOffsetLoad() model.Beam {
	return model.Beam{
		L:       5,
		EI:      1800,
		Support: string(model.SupportCantilever),
		Points:  []model.PointLoad{{At: 3, Force: 40}},
	}
}
