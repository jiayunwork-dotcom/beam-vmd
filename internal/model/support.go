package model

import "fmt"

// ResolveSupport normalizes the support string and returns the canonical
// SupportType. An empty string falls back to DefaultSupport.
func ResolveSupport(s string) (SupportType, error) {
	switch SupportType(s) {
	case "":
		return DefaultSupport, nil
	case SupportSimply, SupportCantilever, SupportFixed:
		return SupportType(s), nil
	case "both": // friendly alias for both-ends-fixed
		return SupportFixed, nil
	default:
		return "", fmt.Errorf("unknown support type %q (want simply|cantilever|fixed)", s)
	}
}

// IsFixedEnd reports whether the given support type restrains rotation (and
// translation) at the named end. end is 0 for the left end (x=0) and 1 for the
// right end (x=L).
func (t SupportType) IsFixedEnd(end int) bool {
	switch t {
	case SupportCantilever:
		return end == 0
	case SupportFixed:
		return end == 0 || end == 1
	default:
		return false
	}
}

// HasSupportAtRight reports whether the support occupies the right end x=L.
func (t SupportType) HasSupportAtRight() bool {
	return t == SupportSimply || t == SupportFixed
}
