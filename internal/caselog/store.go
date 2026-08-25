package caselog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const SnapshotVersion = 1

type Snapshot struct {
	Version int     `json:"version"`
	Max     int     `json:"max"`
	Seq     uint64  `json:"seq"`
	Entries []Entry `json:"entries"`
}

func (l *Log) Export() Snapshot {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := Snapshot{
		Version: SnapshotVersion,
		Max:     l.max,
		Seq:     l.seq,
		Entries: make([]Entry, 0, len(l.items)),
	}
	for _, e := range l.items {
		out.Entries = append(out.Entries, e)
	}
	return out
}

func (l *Log) Import(snapshot Snapshot) error {
	if snapshot.Version != SnapshotVersion {
		return fmt.Errorf("unsupported snapshot version %d", snapshot.Version)
	}
	if snapshot.Max <= 0 {
		return fmt.Errorf("snapshot max must be positive")
	}
	if len(snapshot.Entries) > snapshot.Max {
		return ErrTooMany
	}
	next := make(map[string]Entry, len(snapshot.Entries))
	for _, e := range snapshot.Entries {
		if err := e.Validate(); err != nil {
			return err
		}
		if _, ok := next[e.ID]; ok {
			return ErrExists
		}
		next[e.ID] = e
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = next
	l.max = snapshot.Max
	l.seq = snapshot.Seq
	return nil
}

func (l *Log) SaveFile(path string) error {
	data, err := json.MarshalIndent(l.Export(), "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".beam-log-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}

func LoadFile(path string) (*Log, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, err
	}
	l := NewLog(snapshot.Max)
	if err := l.Import(snapshot); err != nil {
		return nil, err
	}
	return l, nil
}
