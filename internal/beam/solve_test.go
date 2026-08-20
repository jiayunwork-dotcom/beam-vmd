package beam

import (
	"math"
	"testing"

	"beam-vmd/internal/model"
	"beam-vmd/internal/sample"
)

func findSample(samples []model.Sample, x float64) model.Sample {
	for _, s := range samples {
		if math.Abs(s.X-x) < 1e-9 {
			return s
		}
	}
	return samples[len(samples)-1]
}

func TestSolveSimplyClosedForm(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	mid := findSample(res.Samples, 5)
	if got, want := mid.M, 250.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("midspan moment = %g, want %g", got, want)
	}
	// reported deflection is downward positive
	if got, want := mid.Y, 100*1000/(48.0*1000); math.Abs(got-want) > 1e-3 {
		t.Errorf("midspan deflection = %g, want %g", got, want)
	}
	// reactions: each carries half the load upward
	var up0, upL float64
	for _, rct := range res.Reactions {
		if rct.Kind == "vertical" && rct.At == 0 {
			up0 = rct.Force
		}
		if rct.Kind == "vertical" && rct.At == 10 {
			upL = rct.Force
		}
	}
	if math.Abs(up0-50) > 1e-6 || math.Abs(upL-50) > 1e-6 {
		t.Errorf("reactions up = %g, %g; want 50, 50", up0, upL)
	}
}

func TestSolveCantileverClosedForm(t *testing.T) {
	res, err := Solve(sample.ExampleCantilever())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	tip := findSample(res.Samples, 6)
	want := 50 * 216.0 / (3.0 * 2000.0) // P L^3 / (3EI) = 1.8
	if math.Abs(tip.Y-want) > 1e-3 {
		t.Errorf("free-end deflection = %g, want %g", tip.Y, want)
	}
	var fixedMoment float64
	for _, rct := range res.Reactions {
		if rct.Kind == "moment" && rct.At == 0 {
			fixedMoment = rct.Force
		}
	}
	if math.Abs(fixedMoment+300) > 1e-6 {
		t.Errorf("fixed-end moment = %g, want -300", fixedMoment)
	}
}

func TestSolveFixedClosedForm(t *testing.T) {
	res, err := Solve(sample.ExampleFixed())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	mid := findSample(res.Samples, 4)
	if math.Abs(mid.M-120) > 1e-6 {
		t.Errorf("midspan moment = %g, want 120", mid.M)
	}
	want := 120 * 512.0 / (192.0 * 1500.0) // 0.21333...
	if math.Abs(mid.Y-want) > 1e-3 {
		t.Errorf("midspan deflection = %g, want %g", mid.Y, want)
	}
	var m0, mL float64
	for _, rct := range res.Reactions {
		if rct.Kind == "moment" && rct.At == 0 {
			m0 = rct.Force
		}
		if rct.Kind == "moment" && rct.At == 8 {
			mL = rct.Force
		}
	}
	if math.Abs(m0+120) > 1e-6 || math.Abs(mL+120) > 1e-6 {
		t.Errorf("fixed-end moments = %g, %g; want -120, -120", m0, mL)
	}
}
