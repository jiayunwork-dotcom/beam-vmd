package beam

import (
	"fmt"
	"strings"

	"beam-vmd/internal/model"
)

// SampleRows returns at most max sample rows, evenly spaced along the list, for
// compact display. The first and last rows are always included.
func SampleRows(res model.Result, max int) []model.Sample {
	if max <= 0 || len(res.Samples) <= max {
		return res.Samples
	}
	step := float64(len(res.Samples)-1) / float64(max-1)
	out := make([]model.Sample, 0, max)
	for i := 0; i < max; i++ {
		out = append(out, res.Samples[int(float64(i)*step+0.5)])
	}
	return out
}

// Table renders a compact human-readable summary of a solved beam for the CLI.
func Table(res model.Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "support=%s  L=%.6g  EI=%.6g  samples=%d\n",
		res.Support, res.L, res.EI, len(res.Samples))
	sb.WriteString("reactions:\n")
	for _, r := range res.Reactions {
		fmt.Fprintf(&sb, "  x=%.6g %-8s % .6g\n", r.At, r.Kind, r.Force)
	}
	e := res.Extrema
	sb.WriteString("extrema:\n")
	fmt.Fprintf(&sb, "  V  [% .6g, % .6g]\n", e.VMin, e.VMax)
	fmt.Fprintf(&sb, "  M  [% .6g, % .6g]\n", e.MMin, e.MMax)
	fmt.Fprintf(&sb, "  y  [% .6g, % .6g]  |y|max=% .6g at x=%.6g\n",
		e.YMin, e.YMax, e.YAbsMax, e.YAbsMaxX)
	sb.WriteString("samples:\n")
	for _, s := range SampleRows(res, 15) {
		fmt.Fprintf(&sb, "  x=%.6g  V=% .6g  M=% .6g  y=% .6g  theta=% .6g\n",
			s.X, s.V, s.M, s.Y, s.Theta)
	}
	sb.WriteString("checks:\n")
	for _, c := range res.Checks.Items {
		mark := "ok "
		if !c.Passed {
			mark = "!! "
		}
		fmt.Fprintf(&sb, "  %s%s\n", mark, c.Message)
	}
	return sb.String()
}
