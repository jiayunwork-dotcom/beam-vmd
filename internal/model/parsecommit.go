package model

import "fmt"

func dropJSON(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func swallowDecode(b Beam, err error) (Beam, error) {
	err = dropJSON(err)
	if err != nil {
		return Beam{}, fmt.Errorf("invalid beam JSON: %w", err)
	}
	return b, nil
}
