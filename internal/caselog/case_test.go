package caselog

import (
	"errors"
	"testing"

	"beam-vmd/internal/model"
)

func TestLogAddGetRemove(t *testing.T) {
	l := NewLog(4)
	e := Entry{
		ID:     "simply-1",
		Beam:   model.Beam{L: 4, EI: 20000, Support: "simply", Points: []model.PointLoad{{At: 2, Force: 10}}},
		Result: model.Result{L: 4, EI: 20000, Support: "simply", Extrema: model.Extrema{YAbsMax: 0.001}},
	}
	if err := l.Add(e); err != nil {
		t.Fatal(err)
	}
	if l.Len() != 1 || l.NextSeq() != 1 {
		t.Fatalf("len=%d seq=%d", l.Len(), l.NextSeq())
	}
	got, ok := l.Get("simply-1")
	if !ok || got.Beam.L != 4 {
		t.Fatalf("get failed: %+v", got)
	}
	if err := l.Add(e); !errors.Is(err, ErrExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if !l.Remove("simply-1") {
		t.Fatal("remove failed")
	}
}

func TestLogRenameFreezeSetNote(t *testing.T) {
	l := NewLog(8)
	e := Entry{
		ID:     "a",
		Beam:   model.Beam{L: 3, EI: 10000, Support: "cantilever"},
		Result: model.Result{L: 3, EI: 10000, Support: "cantilever"},
	}
	if err := l.Add(e); err != nil {
		t.Fatal(err)
	}
	if err := l.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := l.Freeze("b"); err != nil {
		t.Fatal(err)
	}
	if err := l.SetNote("b", "changed"); !errors.Is(err, ErrFrozen) {
		t.Fatalf("frozen set note err=%v", err)
	}
}

func TestValidateRejectsBadEntries(t *testing.T) {
	bad := []Entry{
		{ID: "", Beam: model.Beam{L: 1, EI: 1}, Result: model.Result{L: 1, EI: 1, Support: "simply"}},
		{ID: "x", Beam: model.Beam{L: 0, EI: 1}, Result: model.Result{L: 1, EI: 1, Support: "simply"}},
		{ID: "x", Beam: model.Beam{L: 1, EI: 0}, Result: model.Result{L: 1, EI: 1, Support: "simply"}},
		{ID: "x", Beam: model.Beam{L: 1, EI: 1}, Result: model.Result{L: 1, EI: 1, Support: ""}},
	}
	for i, e := range bad {
		if err := e.Validate(); err == nil {
			t.Fatalf("entry %d should fail", i)
		}
	}
}

func TestSimilarAndAggregates(t *testing.T) {
	l := NewLog(16)
	entries := []Entry{
		{ID: "a", Beam: model.Beam{L: 4, EI: 20000, Support: "simply"}, Result: model.Result{L: 4, EI: 20000, Support: "simply", Extrema: model.Extrema{YAbsMax: 0.001, MMax: 10}}},
		{ID: "b", Beam: model.Beam{L: 4.2, EI: 21000, Support: "simply"}, Result: model.Result{L: 4.2, EI: 21000, Support: "simply", Extrema: model.Extrema{YAbsMax: 0.0011, MMax: 11}}},
		{ID: "c", Beam: model.Beam{L: 5, EI: 30000, Support: "fixed"}, Result: model.Result{L: 5, EI: 30000, Support: "fixed", Extrema: model.Extrema{YAbsMax: 0.0005, MMax: 8}}},
	}
	for _, e := range entries {
		if err := l.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	sim := l.Similar(model.Beam{L: 4, EI: 20000, Support: "simply"}, 0.03)
	if len(sim) != 1 || sim[0].ID != "a" {
		t.Fatalf("similar=%+v", sim)
	}
	bySupport := l.BySupport()
	if bySupport["simply"] != 2 || bySupport["fixed"] != 1 {
		t.Fatalf("by support=%v", bySupport)
	}
	mean, n := l.MeanStiffness("simply")
	if n != 2 || mean != 20500 {
		t.Fatalf("mean=%v n=%d", mean, n)
	}
	maxY, id := l.MaxDeflection()
	if id != "b" || maxY != 0.0011 {
		t.Fatalf("maxY=%v id=%s", maxY, id)
	}
	ok, checks := l.AllChecksPass()
	if !ok || checks != 0 {
		t.Fatalf("checks ok=%v total=%d", ok, checks)
	}
}
