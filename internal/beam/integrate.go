package beam

import "math"

const integN = 8000

type integrator struct {
	L       float64
	EI      float64
	yPrime0 float64
	xs      []float64
	i1      []float64
	i2      []float64
}

func newIntegrator(M func(float64) float64, L, EI, yPrime0 float64, n int) *integrator {
	xs := make([]float64, n+1)
	me := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		xs[i] = L * float64(i) / float64(n)
		me[i] = M(xs[i]) / EI
	}
	i1 := cumulativeTrap(xs, me)
	i2 := cumulativeTrap(xs, i1)
	return &integrator{L: L, EI: EI, yPrime0: yPrime0, xs: xs, i1: i1, i2: i2}
}

func cumulativeTrap(xs, y []float64) []float64 {
	n := len(xs) - 1
	c := make([]float64, n+1)
	for i := 1; i <= n; i++ {
		h := xs[i] - xs[i-1]
		c[i] = c[i-1] + 0.5*h*(y[i]+y[i-1])
	}
	return c
}

func (it *integrator) theta(x float64) float64 {
	return interp(it.xs, it.i1, x) + it.yPrime0
}

func (it *integrator) y(x float64) float64 {
	return interp(it.xs, it.i2, x) + it.yPrime0*x
}

func (it *integrator) I1L() float64 { return it.i1[len(it.i1)-1] }

func (it *integrator) I2L() float64 { return it.i2[len(it.i2)-1] }

func interp(xs, y []float64, x float64) float64 {
	n := len(xs) - 1
	if x <= xs[0] {
		return y[0]
	}
	if x >= xs[n] {
		return y[n]
	}
	lo, hi := 0, n
	for hi-lo > 1 {
		mid := (lo + hi) / 2
		if xs[mid] <= x {
			lo = mid
		} else {
			hi = mid
		}
	}
	t := (x - xs[lo]) / (xs[hi] - xs[lo])
	return y[lo] + t*(y[hi]-y[lo])
}

func solve2x2(J [2][2]float64, rhs [2]float64) [2]float64 {
	det := J[0][0]*J[1][1] - J[0][1]*J[1][0]
	if math.Abs(det) < 1e-12 {
		return [2]float64{0, 0}
	}
	return [2]float64{
		(J[1][1]*rhs[0] - J[0][1]*rhs[1]) / det,
		(J[0][0]*rhs[1] - J[1][0]*rhs[0]) / det,
	}
}
