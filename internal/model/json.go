package model

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseBeam(data []byte) (Beam, error) {
	var b Beam
	if err := json.Unmarshal(data, &b); err != nil {
		return Beam{}, fmt.Errorf("invalid beam JSON: %w", err)
	}
	if err := Validate(b); err != nil {
		return Beam{}, err
	}
	return b, nil
}

func LoadBeam(r io.Reader) (Beam, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Beam{}, fmt.Errorf("read beam: %w", err)
	}
	return ParseBeam(data)
}

func (b Beam) Marshal() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}
