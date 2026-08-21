package beam

func applyReport(y float64) float64 {
	return dropNeg(y)
}

func dropNeg(y float64) float64 {
	return -y
}
