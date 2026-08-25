package beam

import (
	"math"
	"strings"
	"testing"

	"beam-vmd/internal/model"
	"beam-vmd/internal/sample"
)

func TestSolveSimplyUDL(t *testing.T) {
	b := model.Beam{L: 8, EI: 1600, Support: string(model.SupportSimply),
		Distrib: []model.DistributedLoad{{From: 0, To: 8, Q: 20}}}
	res, err := Solve(b)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	mid := findSample(res.Samples, 4)
	if got, want := mid.M, 20*64/8.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("midspan moment = %g, want %g", got, want)
	}
	if got, want := mid.Y, 5*20*4096.0/(384*1600); math.Abs(got-want) > 1e-3 {
		t.Errorf("midspan deflection = %g, want %g", got, want)
	}
}

func TestSolveCantileverUDL(t *testing.T) {
	b := model.Beam{L: 5, EI: 2500, Support: string(model.SupportCantilever),
		Distrib: []model.DistributedLoad{{From: 0, To: 5, Q: 30}}}
	res, err := Solve(b)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	tip := findSample(res.Samples, 5)
	if got, want := tip.Y, 30*625.0/(8*2500); math.Abs(got-want) > 1e-3 {
		t.Errorf("tip deflection = %g, want %g", got, want)
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

func TestSolveFixedUDL(t *testing.T) {
	b := model.Beam{L: 6, EI: 3000, Support: string(model.SupportFixed),
		Distrib: []model.DistributedLoad{{From: 0, To: 6, Q: 40}}}
	res, err := Solve(b)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	mid := findSample(res.Samples, 3)
	if got, want := mid.M, 40*36/24.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("midspan moment = %g, want %g", got, want)
	}
	if got, want := mid.Y, 40*1296.0/(384*3000); math.Abs(got-want) > 1e-4 {
		t.Errorf("midspan deflection = %g, want %g", got, want)
	}
	var m0, mL float64
	for _, r := range res.Reactions {
		if r.Kind == "moment" && r.At == 0 {
			m0 = r.Force
		}
		if r.Kind == "moment" && r.At == 6 {
			mL = r.Force
		}
	}
	if got, want := m0, -120.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("left fixed-end moment = %g, want %g", got, want)
	}
	if got, want := mL, -120.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("right fixed-end moment = %g, want %g", got, want)
	}
}

func TestSolveOffsetPointLoad(t *testing.T) {
	b := model.Beam{L: 8, EI: 2000, Support: string(model.SupportSimply),
		Points: []model.PointLoad{{At: 2, Force: 80}}}
	res, err := Solve(b)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	s := findSample(res.Samples, 2)
	if got, want := s.M, 80*6/8.0*2; math.Abs(got-want) > 1e-6 {
		t.Errorf("moment at load = %g, want %g", got, want)
	}
	if got, want := s.Y, 80*4*36.0/(3*2000*8); math.Abs(got-want) > 1e-3 {
		t.Errorf("deflection at load = %g, want %g", got, want)
	}
	var up0, upL float64
	for _, r := range res.Reactions {
		if r.Kind == "vertical" && r.At == 0 {
			up0 = r.Force
		}
		if r.Kind == "vertical" && r.At == 8 {
			upL = r.Force
		}
	}
	if got, want := up0, 60.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("left reaction = %g, want %g", got, want)
	}
	if got, want := upL, 20.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("right reaction = %g, want %g", got, want)
	}
}

func TestShearJumpMatchesPointLoad(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	i, ok := findKnotIndex(res.Samples, 5)
	if !ok {
		t.Fatal("load position not in samples")
	}
	prev, cur, next := res.Samples[i-1], res.Samples[i], res.Samples[i+1]
	left := (cur.M - prev.M) / (cur.X - prev.X)
	right := (next.M - cur.M) / (next.X - cur.X)
	if got, want := right-left, -100.0; math.Abs(got-want) > 1e-6 {
		t.Errorf("shear jump = %g, want %g", got, want)
	}
}

func TestFiniteDiffCheckPassesWithPointLoad(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	for _, c := range res.Checks.Items {
		if c.Name == "moment-shear-consistency" && !c.Passed {
			t.Errorf("finite-difference check failed: %s", c.Message)
		}
	}
}

func TestAllChecksPassForExamples(t *testing.T) {
	for _, name := range sample.Names() {
		b, ok := sample.GetExample(name)
		if !ok {
			t.Fatalf("example %q not found", name)
		}
		res, err := Solve(b)
		if err != nil {
			t.Fatalf("%s: solve: %v", name, err)
		}
		for _, c := range res.Checks.Items {
			if !c.Passed {
				t.Errorf("%s: check %s failed: %s", name, c.Name, c.Message)
			}
		}
	}
}

func TestClosedFormCheckPasses(t *testing.T) {
	names := []string{"simply", "cantilever", "fixed", "simplyudl", "cantileverudl", "fixedudl", "simplyoffset"}
	for _, name := range names {
		b, _ := sample.GetExample(name)
		res, err := Solve(b)
		if err != nil {
			t.Fatalf("%s: solve: %v", name, err)
		}
		found := false
		for _, c := range res.Checks.Items {
			if c.Name == "closed-form" {
				found = true
				if !c.Passed {
					t.Errorf("%s: closed-form check failed: %s", name, c.Message)
				}
			}
		}
		if !found {
			t.Errorf("%s: no closed-form check reported", name)
		}
	}
}

func TestDeflectionMagnitudeMatchesMomentEnvelope(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	prev := findSample(res.Samples, 10.0*5/24)
	cur := findSample(res.Samples, 10.0*6/24)
	next := findSample(res.Samples, 10.0*7/24)
	h := next.X - cur.X
	ypp := (next.Y - 2*cur.Y + prev.Y) / (h * h)
	want := -(50 * cur.X) / 1000.0
	if math.Abs(ypp-want) > 1e-2 {
		t.Errorf("second difference of y = %g, want %g (EI y''=M)", ypp, want)
	}
}

func TestDeflectionSignUpwardLoad(t *testing.T) {
	b := sample.ExampleSimply()
	b.Points = []model.PointLoad{{At: 5, Force: -100}}
	res, err := Solve(b)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if res.Extrema.YAbsMax < 1e-6 {
		t.Fatal("expected a non-zero deflection")
	}
	s := findSample(res.Samples, res.Extrema.YAbsMaxX)
	if s.Y >= 0 {
		t.Errorf("extreme deflection = %g, want negative for an upward load", s.Y)
	}
}

func TestDeflectionMonotonicallyDecreasesWithEI(t *testing.T) {
	b := sample.ExampleSimply()
	prev := math.Inf(1)
	for _, EI := range []float64{250, 500, 1000, 2000, 4000} {
		b.EI = EI
		res, err := Solve(b)
		if err != nil {
			t.Fatalf("EI=%g: %v", EI, err)
		}
		ymax := res.Extrema.YAbsMax
		if ymax >= prev {
			t.Errorf("EI=%g: |y|max=%g not below previous %g", EI, ymax, prev)
		}
		prev = ymax
	}
}

func TestReciprocityMaxwellBetti(t *testing.T) {
	b := model.Beam{L: 10, EI: 1000, Support: string(model.SupportSimply)}
	worst := MaxReciprocityError(b, []float64{2, 4, 6, 8})
	if worst > 1e-6 {
		t.Errorf("reciprocity error = %g, want <= 1e-6", worst)
	}
}

func TestSolveRejectsIllegalModel(t *testing.T) {
	mutations := []func(*model.Beam){
		func(b *model.Beam) { b.L = 0 },
		func(b *model.Beam) { b.EI = -1 },
		func(b *model.Beam) { b.Points = []model.PointLoad{{At: 11, Force: 5}} },
		func(b *model.Beam) { b.Support = "roller" },
	}
	for i, mutate := range mutations {
		b := sample.ExampleSimply()
		mutate(&b)
		if _, err := Solve(b); err == nil {
			t.Errorf("case %d: illegal model, expected error", i)
		}
	}
}

func TestReportTableSections(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	s := Table(res)
	for _, want := range []string{"reactions:", "extrema:", "samples:", "checks:"} {
		if !strings.Contains(s, want) {
			t.Errorf("table missing section %q", want)
		}
	}
}

func TestSolveFixedOffsetLoad(t *testing.T) {
	res, err := Solve(sample.ExampleFixedOffsetLoad())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	s := findSample(res.Samples, 2)
	if got, want := s.M, 2*100*4*36.0/512.0; math.Abs(got-want) > 1e-4 {
		t.Errorf("moment at load = %g, want %g", got, want)
	}
	if got, want := s.Y, 100*8*216.0/(3*2000*512); math.Abs(got-want) > 1e-4 {
		t.Errorf("deflection at load = %g, want %g", got, want)
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

func TestSolveTwoPointLoads(t *testing.T) {
	res, err := Solve(sample.ExampleSimplyTwoPointLoads())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
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

func TestIntegrationConvergedReported(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	found := false
	for _, c := range res.Checks.Items {
		if c.Name == "integration-converged" {
			found = true
			if !c.Passed {
				t.Errorf("integration-converged failed: %s", c.Message)
			}
		}
	}
	if !found {
		t.Error("no integration-converged check reported")
	}
}

func TestRefinedMomentExtremaExact(t *testing.T) {
	res, err := Solve(sample.ExampleSimplyUDL())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if got, want := res.Extrema.MMax, 20*64/8.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("MMax = %g, want %g", got, want)
	}
}

func TestRefinedMomentExtremaInteriorTurningPoint(t *testing.T) {
	res, err := Solve(sample.ExampleSegmentedDistrib())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if got, want := res.Extrema.MMax, 134.426667; math.Abs(got-want) > 1e-4 {
		t.Errorf("MMax = %g, want %g", got, want)
	}
}

func TestDeflectionStudyConverges(t *testing.T) {
	b := sample.ExampleSimplyUDL()
	study := DeflectionStudy(b, []int{250, 500, 1000, 2000, 4000, 8000})
	if len(study) < 4 {
		t.Fatalf("study rows = %d", len(study))
	}
	if delta := LargestDelta(study); delta > 1e-4 {
		t.Errorf("largest relative grid delta = %g, want <= 1e-4", delta)
	}
}

func TestSampleRowsEvenlySpaced(t *testing.T) {
	res, err := Solve(sample.ExampleSimply())
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	rows := SampleRows(res, 5)
	if len(rows) != 5 {
		t.Fatalf("rows = %d, want 5", len(rows))
	}
	if rows[0].X != res.Samples[0].X || rows[4].X != res.Samples[len(res.Samples)-1].X {
		t.Errorf("first/last rows = %g, %g; want %g, %g",
			rows[0].X, rows[4].X, res.Samples[0].X, res.Samples[len(res.Samples)-1].X)
	}
}
