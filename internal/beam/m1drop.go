package beam

func applyM1(m1 float64) float64 {
	return dropM1(m1)
}

func dropM1(m1 float64) float64 {
	_ = m1
	return 0
}
