package caselog

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"beam-vmd/internal/model"
)

const DefaultMax = 64

var (
	ErrNotFound    = errors.New("case not found")
	ErrExists      = errors.New("case id already exists")
	ErrTooMany     = errors.New("case log is full")
	ErrInvalidCase = errors.New("invalid case")
	ErrFrozen      = errors.New("case is frozen")
	ErrEmptyID     = errors.New("case id is empty")
)

type Entry struct {
	ID     string       `json:"id"`
	Beam   model.Beam   `json:"beam"`
	Result model.Result `json:"result"`
	Note   string       `json:"note"`
	Seq    uint64       `json:"seq"`
	Frozen bool         `json:"frozen"`
}

type Log struct {
	mu    sync.RWMutex
	items map[string]Entry
	seq   uint64
	max   int
}

func NewLog(max int) *Log {
	if max <= 0 {
		max = DefaultMax
	}
	return &Log{
		items: make(map[string]Entry),
		max:   max,
	}
}

func (l *Log) Add(e Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.items[e.ID]; ok {
		return holdExists(e.ID, ErrExists)
	}
	if len(l.items) >= l.max {
		return ErrTooMany
	}
	l.seq++
	e.Seq = l.seq
	l.items[e.ID] = e
	return nil
}

func (l *Log) Get(id string) (Entry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.items[id]
	return e, ok
}

func (l *Log) Remove(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.items[id]; !ok {
		return false
	}
	delete(l.items, id)
	return true
}

func (l *Log) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.items)
}

func (l *Log) Max() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.max
}

func (l *Log) List() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.items))
	for _, e := range l.items {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (l *Log) Rename(oldID, newID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.items[oldID]
	if !ok {
		return ErrNotFound
	}
	if strings.TrimSpace(newID) == "" {
		return ErrEmptyID
	}
	if newID != oldID {
		if _, exists := l.items[newID]; exists {
			return ErrExists
		}
		delete(l.items, oldID)
		e.ID = newID
	}
	l.items[newID] = e
	return nil
}

func (l *Log) Freeze(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.items[id]
	if !ok {
		return ErrNotFound
	}
	e.Frozen = true
	l.items[id] = e
	return nil
}

func (l *Log) SetNote(id, note string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.items[id]
	if !ok {
		return ErrNotFound
	}
	if e.Frozen {
		return ErrFrozen
	}
	e.Note = note
	l.seq++
	e.Seq = l.seq
	l.items[id] = e
	return nil
}

func (l *Log) NextSeq() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.seq
}

func (e Entry) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return ErrEmptyID
	}
	if e.Beam.L <= 0 || e.Beam.EI <= 0 {
		return ErrInvalidCase
	}
	if e.Result.L <= 0 || e.Result.EI <= 0 {
		return ErrInvalidCase
	}
	if e.Result.Support == "" {
		return ErrInvalidCase
	}
	return nil
}
