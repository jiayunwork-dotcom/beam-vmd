package model

type SupportType string

const (
	SupportSimply     SupportType = "simply"
	SupportCantilever SupportType = "cantilever"
	SupportFixed      SupportType = "fixed"
)

const DefaultSupport = SupportSimply

type PointLoad struct {
	At    float64 `json:"at"`
	Force float64 `json:"force"`
}

type DistributedLoad struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
	Q    float64 `json:"q"`
}

type Beam struct {
	L       float64           `json:"L"`
	EI      float64           `json:"EI"`
	Support string            `json:"support,omitempty"`
	Points  []PointLoad       `json:"points,omitempty"`
	Distrib []DistributedLoad `json:"distrib,omitempty"`
}

type Reaction struct {
	At    float64 `json:"at"`
	Kind  string  `json:"kind"`
	Force float64 `json:"force"`
}

type Extrema struct {
	VMax     float64 `json:"v_max"`
	VMin     float64 `json:"v_min"`
	MMax     float64 `json:"m_max"`
	MMin     float64 `json:"m_min"`
	YMax     float64 `json:"y_max"`
	YMin     float64 `json:"y_min"`
	YAbsMax  float64 `json:"y_abs_max"`
	YAbsMaxX float64 `json:"y_abs_max_x"`
}

type Sample struct {
	X     float64 `json:"x"`
	V     float64 `json:"V"`
	M     float64 `json:"M"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
}

type CheckItem struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

type CheckReport struct {
	Items []CheckItem `json:"items"`
}

type Result struct {
	L         float64     `json:"L"`
	EI        float64     `json:"EI"`
	Support   string      `json:"support"`
	Reactions []Reaction  `json:"reactions"`
	Samples   []Sample    `json:"samples"`
	Extrema   Extrema     `json:"extrema"`
	Checks    CheckReport `json:"checks"`
}
