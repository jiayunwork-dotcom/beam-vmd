package model

import "fmt"

func ResolveSupport(s string) (SupportType, error) {
	switch SupportType(s) {
	case "":
		return DefaultSupport, nil
	case SupportSimply, SupportCantilever, SupportFixed:
		return SupportType(s), nil
	case "both":
		return SupportFixed, nil
	default:
		return "", fmt.Errorf("unknown support type %q (want simply|cantilever|fixed)", s)
	}
}

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

func (t SupportType) HasSupportAtRight() bool {
	return t == SupportSimply || t == SupportFixed
}
