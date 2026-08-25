package beam

import "beam-vmd/internal/model"

func computeReactions(b model.Beam) (reactionSet, error) {
	sup, err := model.ResolveSupport(b.Support)
	if err != nil {
		return reactionSet{}, err
	}
	switch sup {
	case model.SupportSimply:
		return solveSimply(b), nil
	case model.SupportCantilever:
		return solveCantilever(b), nil
	case model.SupportFixed:
		return solveFixed(b), nil
	}
	return reactionSet{}, model.IllegalModel("unsupported support in solver")
}
