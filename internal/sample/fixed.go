package sample

import "beam-vmd/internal/model"

func ExampleFixed() model.Beam {
	return model.Beam{
		L:       8,
		EI:      1500,
		Support: string(model.SupportFixed),
		Points:  []model.PointLoad{{At: 4, Force: 120}},
	}
}
