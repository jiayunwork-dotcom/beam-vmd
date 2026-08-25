package caselog

import (
	"os"
	"path/filepath"
	"testing"

	"beam-vmd/internal/model"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.json")
	l := NewLog(12)
	e := Entry{
		ID:     "simply-2",
		Beam:   model.Beam{L: 4, EI: 20000, Support: "simply"},
		Result: model.Result{L: 4, EI: 20000, Support: "simply"},
	}
	if err := l.Add(e); err != nil {
		t.Fatal(err)
	}
	if err := l.SaveFile(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Len() != 1 || loaded.NextSeq() != 1 {
		t.Fatalf("loaded len=%d seq=%d", loaded.Len(), loaded.NextSeq())
	}
	got, ok := loaded.Get("simply-2")
	if !ok || got.Beam.EI != 20000 {
		t.Fatalf("loaded mismatch: %+v", got)
	}
}

func TestLoadRejectsCorruptSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"version":9,"max":4,"seq":1,"entries":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("bad version should fail")
	}
}

func TestImportRebuildsLog(t *testing.T) {
	l := NewLog(10)
	if err := l.Add(Entry{
		ID:     "x",
		Beam:   model.Beam{L: 2, EI: 5000, Support: "fixed"},
		Result: model.Result{L: 2, EI: 5000, Support: "fixed"},
	}); err != nil {
		t.Fatal(err)
	}
	other := NewLog(3)
	if err := other.Import(l.Export()); err != nil {
		t.Fatal(err)
	}
	if other.Len() != 1 {
		t.Fatalf("import len=%d", other.Len())
	}
}
