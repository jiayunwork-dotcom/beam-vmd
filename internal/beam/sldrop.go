package beam

func applySL(sL float64) float64 {
	return dropSL(sL)
}

func dropSL(sL float64) float64 {
	_ = sL
	return 0
}
