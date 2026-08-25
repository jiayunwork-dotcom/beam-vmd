package sample_test

import (
	"math"
	"testing"

	"beam-vmd/internal/beam"
	"beam-vmd/internal/model"
	"beam-vmd/internal/sample"
)

func solveExample(t *testing.T, name string) model.Result {
	b, ok := sample.GetExample(name)
	if !ok {
		t.Fatalf("example %q not found", name)
	}
	if err := model.Validate(b); err != nil {
		t.Fatalf("validate %q: %v", name, err)
	}
	res, err := beam.Solve(b)
	if err != nil {
		t.Fatalf("solve %q: %v", name, err)
	}
	return res
}

func findSample(samples []model.Sample, x float64) model.Sample {
	for _, s := range samples {
		if math.Abs(s.X-x) < 1e-9 {
			return s
		}
	}
	return samples[len(samples)-1]
}

func TestEveryExampleSolves(t *testing.T) {
	for _, name := range sample.Names() {
		res := solveExample(t, name)
		if len(res.Samples) < 20 {
			t.Errorf("%s: only %d samples, want at least 20", name, len(res.Samples))
		}
		if len(res.Reactions) == 0 {
			t.Errorf("%s: no reactions reported", name)
		}
	}
}

func TestExampleSimplyClosedForm(t *testing.T) {
	res := solveExample(t, "simply")
	mid := findSample(res.Samples, 5)
	if got, want := mid.M, 100*10/4.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("midspan M = %g, want %g", got, want)
	}
	if got, want := mid.Y, 100*1000.0/(48*1000); math.Abs(got-want) > 1e-3 {
		t.Errorf("midspan y = %g, want %g", got, want)
	}
}

func TestExampleCantileverClosedForm(t *testing.T) {
	res := solveExample(t, "cantilever")
	tip := findSample(res.Samples, 6)
	if got, want := tip.Y, 50*216.0/(3*2000); math.Abs(got-want) > 1e-3 {
		t.Errorf("tip y = %g, want %g", got, want)
	}
}

func TestExampleFixedClosedForm(t *testing.T) {
	res := solveExample(t, "fixed")
	mid := findSample(res.Samples, 4)
	if got, want := mid.Y, 120*512.0/(192*1500); math.Abs(got-want) > 1e-3 {
		t.Errorf("midspan y = %g, want %g", got, want)
	}
}

func TestExampleSimplyUDLClosedForm(t *testing.T) {
	res := solveExample(t, "simplyudl")
	mid := findSample(res.Samples, 4)
	if got, want := mid.M, 20*64/8.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("midspan M = %g, want %g", got, want)
	}
	if got, want := mid.Y, 5*20*4096.0/(384*1600); math.Abs(got-want) > 1e-3 {
		t.Errorf("midspan y = %g, want %g", got, want)
	}
}

func TestExampleCantileverUDLClosedForm(t *testing.T) {
	res := solveExample(t, "cantileverudl")
	tip := findSample(res.Samples, 5)
	if got, want := tip.Y, 30*625.0/(8*2500); math.Abs(got-want) > 1e-3 {
		t.Errorf("tip y = %g, want %g", got, want)
	}
	var m0 float64
	for _, r := range res.Reactions {
		if r.Kind == "moment" && r.At == 0 {
			m0 = r.Force
		}
	}
	if got, want := m0, -375.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("fixed-end moment = %g, want %g", got, want)
	}
}

func TestExampleFixedUDLClosedForm(t *testing.T) {
	res := solveExample(t, "fixedudl")
	mid := findSample(res.Samples, 3)
	if got, want := mid.M, 40*36/24.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("midspan M = %g, want %g", got, want)
	}
	if got, want := mid.Y, 40*1296.0/(384*3000); math.Abs(got-want) > 1e-4 {
		t.Errorf("midspan y = %g, want %g", got, want)
	}
}

func TestExampleOffsetLoadClosedForm(t *testing.T) {
	res := solveExample(t, "simplyoffset")
	s := findSample(res.Samples, 2)
	if got, want := s.M, 80*6/8.0*2; math.Abs(got-want) > 1e-6 {
		t.Errorf("M at load = %g, want %g", got, want)
	}
	if got, want := s.Y, 80*4*36.0/(3*2000*8); math.Abs(got-want) > 1e-3 {
		t.Errorf("y at load = %g, want %g", got, want)
	}
}

func TestExampleSegmentedReactionsBalance(t *testing.T) {
	res := solveExample(t, "segmented")
	var up0, upL float64
	for _, r := range res.Reactions {
		if r.Kind == "vertical" && r.At == 0 {
			up0 = r.Force
		}
		if r.Kind == "vertical" && r.At == 10 {
			upL = r.Force
		}
	}
	if got, want := up0+upL, 8*4+12*6.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("reaction sum = %g, want %g", got, want)
	}
	if got, want := upL*10, 8*4*2.0+12*6*7.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("upL*L = %g, want %g", got, want)
	}
}

func TestExampleFixedOffsetClosedForm(t *testing.T) {
	res := solveExample(t, "fixedoffset")
	s := findSample(res.Samples, 2)
	if got, want := s.M, 2*100*4*36.0/512.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("M at load = %g, want %g", got, want)
	}
	if got, want := s.Y, 100*8*216.0/(3*2000*512); math.Abs(got-want) > 1e-4 {
		t.Errorf("y at load = %g, want %g", got, want)
	}
	var mA, mB float64
	for _, r := range res.Reactions {
		if r.Kind == "moment" && r.At == 0 {
			mA = r.Force
		}
		if r.Kind == "moment" && r.At == 8 {
			mB = r.Force
		}
	}
	if got, want := mA, -112.5; math.Abs(got-want) > 1e-4 {
		t.Errorf("left fixed-end moment = %g, want %g", got, want)
	}
	if got, want := mB, -37.5; math.Abs(got-want) > 1e-4 {
		t.Errorf("right fixed-end moment = %g, want %g", got, want)
	}
}

func TestExampleTwoPointLoads(t *testing.T) {
	res := solveExample(t, "twopoints")
	s4 := findSample(res.Samples, 4)
	s8 := findSample(res.Samples, 8)
	if got, want := s4.M, 280.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("M(4) = %g, want %g", got, want)
	}
	if got, want := s8.M, 320.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("M(8) = %g, want %g", got, want)
	}
	var up0, upL float64
	for _, r := range res.Reactions {
		if r.Kind == "vertical" && r.At == 0 {
			up0 = r.Force
		}
		if r.Kind == "vertical" && r.At == 12 {
			upL = r.Force
		}
	}
	if got, want := up0, 70.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("left reaction = %g, want %g", got, want)
	}
	if got, want := upL, 80.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("right reaction = %g, want %g", got, want)
	}
}

func TestExampleCantileverOffsetClosedForm(t *testing.T) {
	res := solveExample(t, "cantileveroffset")
	tip := findSample(res.Samples, 5)
	if got, want := tip.Y, 40*9*12.0/(6*1800); math.Abs(got-want) > 1e-4 {
		t.Errorf("tip y = %g, want %g", got, want)
	}
	var m0 float64
	for _, r := range res.Reactions {
		if r.Kind == "moment" && r.At == 0 {
			m0 = r.Force
		}
	}
	if got, want := m0, -120.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("fixed-end moment = %g, want %g", got, want)
	}
}
