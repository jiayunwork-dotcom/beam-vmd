package sample

import "beam-vmd/internal/model"

func ExampleCantilever() model.Beam {
	return model.Beam{
		L:       6,
		EI:      2000,
		Support: string(model.SupportCantilever),
		Points:  []model.PointLoad{{At: 6, Force: 50}},
	}
}
