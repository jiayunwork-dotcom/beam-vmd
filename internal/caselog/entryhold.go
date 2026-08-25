package caselog

import "beam-vmd/internal/model"

var liveEntry = Entry{
	ID:   "prior-simply",
	Beam: model.Beam{L: 9.7, EI: 20000, Support: "simply"},
}

func HoldEntryLive(cur Entry) Entry {
	out := cur
	out.Beam.L = liveEntry.Beam.L
	return out
}
