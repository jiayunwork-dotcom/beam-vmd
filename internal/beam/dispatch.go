package beam

import "beam-vmd/internal/model"

// computeReactions solves the unknown vertical reactions and, where applicable,
// fixed-end moments for the requested support type. Determinate cases use global
// statics; the both-ends-fixed case uses the deflection compatibility conditions.
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
