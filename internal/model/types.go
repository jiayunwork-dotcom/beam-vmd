// Package model defines the data types, JSON (de)serialization, validation
// and the shared sign conventions for the beam-vmd solver.
//
// Coordinate and sign conventions (documented here and in the project README):
//
//   - x runs from 0 (left end) to L (right end).
//   - Transverse loads use a "downward positive" convention: a point force or a
//     distributed intensity q with a positive value acts downward. Negative
//     values act upward. This is the single convention used for all external
//     loads and for the internal reaction vertical forces while solving.
//   - Bending moment M(x) uses the "sagging positive" convention: positive M
//     puts tension on the bottom fibre (sags the beam).
//   - Internal integration uses EI y'' = M with y measured upward positive; the
//     reported deflection is the negative of that (so a downward load produces a
//     positive reported deflection). See the beam package for details.
package model

// SupportType enumerates the supported boundary conditions.
type SupportType string

const (
	// SupportSimply is a simply supported beam: pin at x=0, roller at x=L.
	// Both ends carry a vertical reaction and allow free rotation.
	SupportSimply SupportType = "simply"
	// SupportCantilever is fixed at x=0 and free at x=L.
	SupportCantilever SupportType = "cantilever"
	// SupportFixed is fixed at both ends (both-ends-fixed).
	SupportFixed SupportType = "fixed"
)

// DefaultSupport is used when the JSON document omits the support field.
const DefaultSupport = SupportSimply

// PointLoad is a concentrated transverse force.
type PointLoad struct {
	At    float64 `json:"at"`    // position x in [0, L]
	Force float64 `json:"force"` // signed: +down, -up
}

// DistributedLoad is a (possibly partial) uniformly distributed transverse load.
type DistributedLoad struct {
	From float64 `json:"from"` // start position in [0, L]
	To   float64 `json:"to"`   // end position in [0, L], must be >= From
	Q    float64 `json:"q"`    // signed intensity: +down, -up
}

// Beam is the full problem description supplied by the user.
type Beam struct {
	L       float64          `json:"L"`
	EI      float64          `json:"EI"`
	Support string           `json:"support,omitempty"` // empty => DefaultSupport
	Points  []PointLoad      `json:"points,omitempty"`
	Distrib []DistributedLoad `json:"distrib,omitempty"`
}

// Reaction reports one support reaction. Vertical forces are reported with
// upward positive; moments are reported with sagging positive (tension on the
// bottom fibre).
type Reaction struct {
	At    float64 `json:"at"`
	Kind  string  `json:"kind"`  // "vertical" | "moment"
	Force float64 `json:"force"` // vertical: up-positive; moment: sagging-positive
}

// Extrema summarizes the extreme values found along the beam.
type Extrema struct {
	VMax     float64 `json:"v_max"`
	VMin     float64 `json:"v_min"`
	MMax     float64 `json:"m_max"`
	MMin     float64 `json:"m_min"`
	YMax     float64 `json:"y_max"`     // reported deflection (down positive)
	YMin     float64 `json:"y_min"`     // reported deflection (down positive)
	YAbsMax  float64 `json:"y_abs_max"` // max |reported deflection|
	YAbsMaxX float64 `json:"y_abs_max_x"`
}

// Sample is one sampled cross-section.
type Sample struct {
	X      float64 `json:"x"`
	V      float64 `json:"V"`      // internal shear force
	M      float64 `json:"M"`      // bending moment (sagging positive)
	Y      float64 `json:"y"`      // reported deflection (down positive)
	Theta  float64 `json:"theta"`  // slope (internal, up positive)
}

// CheckItem is the outcome of one cross-validation rule.
type CheckItem struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

// CheckReport aggregates all cross-validation rules.
type CheckReport struct {
	Items []CheckItem `json:"items"`
}

// Result is the full solver response returned by the API.
type Result struct {
	L        float64    `json:"L"`
	EI       float64    `json:"EI"`
	Support  string     `json:"support"`
	Reactions []Reaction `json:"reactions"`
	Samples  []Sample   `json:"samples"`
	Extrema  Extrema    `json:"extrema"`
	Checks   CheckReport `json:"checks"`
}
