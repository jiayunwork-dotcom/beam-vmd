package sample

import "beam-vmd/internal/model"

func Names() []string {
	return []string{
		"simply", "cantilever", "fixed",
		"simplyudl", "cantileverudl", "fixedudl",
		"simplyoffset", "segmented", "twopoints", "fixedoffset", "cantileveroffset",
	}
}

func GetExample(name string) (model.Beam, bool) {
	switch name {
	case "", "simply":
		return ExampleSimply(), true
	case "cantilever":
		return ExampleCantilever(), true
	case "fixed":
		return ExampleFixed(), true
	case "simplyudl":
		return ExampleSimplyUDL(), true
	case "cantileverudl":
		return ExampleCantileverUDL(), true
	case "fixedudl":
		return ExampleFixedUDL(), true
	case "simplyoffset":
		return ExampleSimplyOffset(), true
	case "segmented":
		return ExampleSegmentedDistrib(), true
	case "twopoints":
		return ExampleSimplyTwoPointLoads(), true
	case "fixedoffset":
		return ExampleFixedOffsetLoad(), true
	case "cantileveroffset":
		return ExampleCantileverOffsetLoad(), true
	default:
		return model.Beam{}, false
	}
}
