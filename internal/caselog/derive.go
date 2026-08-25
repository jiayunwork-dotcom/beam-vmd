package caselog

import (
	"math"
	"sort"

	"beam-vmd/internal/model"
)

func (l *Log) MaxDeflection() (float64, string) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	best := 0.0
	id := ""
	for _, e := range l.items {
		abs := math.Abs(e.Result.Extrema.YAbsMax)
		if abs > best {
			best = abs
			id = e.ID
		}
	}
	return best, id
}

func (l *Log) MaxMoment() (float64, string) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	best := 0.0
	id := ""
	for _, e := range l.items {
		abs := math.Max(math.Abs(e.Result.Extrema.MMax), math.Abs(e.Result.Extrema.MMin))
		if abs > best {
			best = abs
			id = e.ID
		}
	}
	return best, id
}

func (l *Log) AllChecksPass() (bool, int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	total := 0
	for _, e := range l.items {
		for _, c := range e.Result.Checks.Items {
			total++
			if !c.Passed {
				return false, total
			}
		}
	}
	return true, total
}

func (l *Log) Similar(target model.Beam, spanRel float64) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range l.items {
		if e.Beam.Support != target.Support {
			continue
		}
		dL := math.Abs(e.Beam.L-target.L) / target.L
		dEI := math.Abs(e.Beam.EI-target.EI) / target.EI
		if dL <= spanRel && dEI <= spanRel {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (l *Log) BySupport() map[string]int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]int)
	for _, e := range l.items {
		out[e.Result.Support]++
	}
	return out
}

func (l *Log) MeanStiffness(support string) (float64, int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range l.items {
		if support != "" && e.Result.Support != support {
			continue
		}
		sum += e.Beam.EI
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}
