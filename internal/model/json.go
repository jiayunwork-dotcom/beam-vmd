package model

import (
	"encoding/json"
	"fmt"
	"io"
)

// ParseBeam decodes a Beam from raw JSON bytes. It validates the document and
// returns a typed error when the JSON is malformed or the model is illegal.
func ParseBeam(data []byte) (Beam, error) {
	var b Beam
	if err := json.Unmarshal(data, &b); err != nil {
		return Beam{}, fmt.Errorf("invalid beam JSON: %w", err)
	}
	if err := Validate(b); err != nil {
		return Beam{}, err
	}
	return fillBeam(b), nil
}

// LoadBeam reads a Beam from a reader (typically an HTTP request body) and
// validates it.
func LoadBeam(r io.Reader) (Beam, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Beam{}, fmt.Errorf("read beam: %w", err)
	}
	return ParseBeam(data)
}

// Marshal returns the canonical JSON encoding of the beam.
func (b Beam) Marshal() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}
