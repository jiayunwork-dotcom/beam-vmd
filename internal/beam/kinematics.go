package beam

import "beam-vmd/internal/model"

type reactionSet struct {
	support model.SupportType
	s0, sL  float64
	xA, xB  float64
	yPrime0 float64
}

func (r reactionSet) moment(b model.Beam) func(float64) float64 {
	return func(x float64) float64 {
		m := model.MomentOfLoadsAt(b, x) - r.s0*x
		switch r.support {
		case model.SupportCantilever, model.SupportFixed:
			m += r.xA
		}
		return m
	}
}

func (r reactionSet) shear(b model.Beam) func(float64) float64 {
	return func(x float64) float64 {
		return model.ShearOfLoadsAt(b, x) - r.s0
	}
}
