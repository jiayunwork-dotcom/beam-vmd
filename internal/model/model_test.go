package model

import "testing"

func validBeam() Beam {
	return Beam{
		L:       10,
		EI:      1000,
		Support: string(SupportSimply),
		Points:  []PointLoad{{At: 5, Force: 100}},
	}
}

func TestValidateAcceptsValidBeam(t *testing.T) {
	if err := Validate(validBeam()); err != nil {
		t.Errorf("Validate: %v", err)
	}
}

func TestValidateRejectsNonPositiveSpan(t *testing.T) {
	for _, L := range []float64{0, -3} {
		b := validBeam()
		b.L = L
		if err := Validate(b); err == nil {
			t.Errorf("L=%g: expected error", L)
		}
	}
}

func TestValidateRejectsNonPositiveEI(t *testing.T) {
	for _, EI := range []float64{0, -5} {
		b := validBeam()
		b.EI = EI
		if err := Validate(b); err == nil {
			t.Errorf("EI=%g: expected error", EI)
		}
	}
}

func TestValidateRejectsPointLoadOutsideSpan(t *testing.T) {
	for _, at := range []float64{-1, 11} {
		b := validBeam()
		b.Points = []PointLoad{{At: at, Force: 10}}
		if err := Validate(b); err == nil {
			t.Errorf("at=%g: expected error", at)
		}
	}
}

func TestValidateRejectsBadDistributedSegment(t *testing.T) {
	bad := []DistributedLoad{
		{From: -1, To: 5, Q: 2},
		{From: 0, To: 12, Q: 2},
		{From: 7, To: 3, Q: 2},
	}
	for _, d := range bad {
		b := validBeam()
		b.Distrib = []DistributedLoad{d}
		if err := Validate(b); err == nil {
			t.Errorf("segment %+v: expected error", d)
		}
	}
}

func TestValidateRejectsUnknownSupport(t *testing.T) {
	b := validBeam()
	b.Support = "pinned"
	if err := Validate(b); err == nil {
		t.Error("unknown support: expected error")
	}
}

func TestResolveSupportDefaultsAndAlias(t *testing.T) {
	if s, err := ResolveSupport(""); err != nil || s != DefaultSupport {
		t.Errorf("empty support: got %q, %v", s, err)
	}
	if s, err := ResolveSupport("both"); err != nil || s != SupportFixed {
		t.Errorf("alias 'both': got %q, %v", s, err)
	}
	if s, err := ResolveSupport("roller"); err == nil {
		t.Errorf("unknown support: got %q, expected error", s)
	}
}

func TestParseBeamRejectsMalformedJSON(t *testing.T) {
	if _, err := ParseBeam([]byte("{not json")); err == nil {
		t.Error("malformed JSON: expected error")
	}
}

func TestParseBeamValidatesDocument(t *testing.T) {
	doc := `{"L":10,"EI":1000,"support":"cantilever","points":[{"at":6,"force":50}]}`
	b, err := ParseBeam([]byte(doc))
	if err != nil {
		t.Fatalf("ParseBeam: %v", err)
	}
	if b.L != 10 || b.EI != 1000 || b.Support != "cantilever" {
		t.Errorf("parsed fields: %+v", b)
	}
	if len(b.Points) != 1 || b.Points[0].At != 6 || b.Points[0].Force != 50 {
		t.Errorf("parsed points: %+v", b.Points)
	}
}

func TestSummarizeLoadsNetResultant(t *testing.T) {
	b := Beam{
		L:       10,
		EI:      1000,
		Points:  []PointLoad{{At: 2, Force: 30}, {At: 6, Force: -10}},
		Distrib: []DistributedLoad{{From: 0, To: 4, Q: 5}},
	}
	s := SummarizeLoads(b)
	if got, want := s.TotalDown, 30.0-10+5*4; got != want {
		t.Errorf("TotalDown = %g, want %g", got, want)
	}
	if got, want := s.MomentAboutOrigin, 30*2.0-10*6+5*4*2.0; got != want {
		t.Errorf("MomentAboutOrigin = %g, want %g", got, want)
	}
}

func TestNormalizeSortsAndSplits(t *testing.T) {
	b := Beam{
		L:       10,
		EI:      1000,
		Support: "both",
		Points:  []PointLoad{{At: 3, Force: 5}, {At: 1, Force: 2}},
		Distrib: []DistributedLoad{{From: 0, To: 10, Q: 4}},
	}
	n := Normalize(b)
	if n.Support != string(SupportFixed) {
		t.Errorf("support = %q, want %q", n.Support, SupportFixed)
	}
	if n.Points[0].At != 1 || n.Points[1].At != 3 {
		t.Errorf("points not sorted: %+v", n.Points)
	}
	if len(n.Distrib) != 3 {
		t.Fatalf("distrib segments = %d, want 3", len(n.Distrib))
	}
	if got, want := SummarizeLoads(n).TotalDown, 4*10+5+2.0; got != want {
		t.Errorf("total after split = %g, want %g", got, want)
	}
}

func TestLoadPiecesCoverSpan(t *testing.T) {
	b := Beam{
		L:       10,
		EI:      1000,
		Points:  []PointLoad{{At: 5, Force: 100}},
		Distrib: []DistributedLoad{{From: 2, To: 8, Q: 4}},
	}
	pieces := LoadPieces(b)
	wantQ := []float64{0, 4, 4, 0}
	if len(pieces) != len(wantQ) {
		t.Fatalf("pieces = %d, want %d", len(pieces), len(wantQ))
	}
	for i, p := range pieces {
		if p.Q != wantQ[i] {
			t.Errorf("piece %d Q = %g, want %g", i, p.Q, wantQ[i])
		}
	}
	total := 0.0
	for _, p := range pieces {
		r, _ := p.Resultant()
		total += r
	}
	if got, want := total, 4*6.0; got != want {
		t.Errorf("piece resultant total = %g, want %g", got, want)
	}
}

func TestIntensityAtOutsideSegments(t *testing.T) {
	b := Beam{L: 10, EI: 1000,
		Distrib: []DistributedLoad{{From: 2, To: 8, Q: 4}}}
	for _, x := range []float64{0, 1, 9, 10} {
		if got := IntensityAt(b, x); got != 0 {
			t.Errorf("IntensityAt(%g) = %g, want 0", x, got)
		}
	}
	if got := IntensityAt(b, 5); got != 4 {
		t.Errorf("IntensityAt(5) = %g, want 4", got)
	}
}

func TestResultantLoad(t *testing.T) {
	b := Beam{
		L:       10,
		EI:      1000,
		Points:  []PointLoad{{At: 2, Force: 30}, {At: 6, Force: -10}},
		Distrib: []DistributedLoad{{From: 0, To: 4, Q: 5}},
	}
	total, centroid := ResultantLoad(b)
	if got, want := total, 30.0-10+20; got != want {
		t.Errorf("total = %g, want %g", got, want)
	}
	if got, want := centroid, (30*2.0-10*6+5*4*2.0)/total; got != want {
		t.Errorf("centroid = %g, want %g", got, want)
	}
}
