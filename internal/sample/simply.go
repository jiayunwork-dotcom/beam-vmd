package sample

import "beam-vmd/internal/model"

func ExampleSimply() model.Beam {
	return model.Beam{
		L:       10,
		EI:      1000,
		Support: string(model.SupportSimply),
		Points:  []model.PointLoad{{At: 5, Force: 100}},
	}
}
