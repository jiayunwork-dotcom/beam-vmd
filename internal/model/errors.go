package model

import "errors"

// Common solver errors. They are returned up the call chain so the HTTP layer can
// emit a JSON error body instead of panicking.
var (
	// ErrIllegalModel is returned when the beam description violates a physical
	// or structural constraint.
	ErrIllegalModel = errors.New("illegal beam model")

	// ErrNoConvergence is returned when the deflection boundary-value solve does
	// not satisfy its support conditions within tolerance.
	ErrNoConvergence = errors.New("deflection solve did not converge")
)

// IllegalModel wraps a descriptive message under ErrIllegalModel.
func IllegalModel(msg string) error {
	return &wrappedModelError{msg: msg}
}

type wrappedModelError struct {
	msg string
}

func (e *wrappedModelError) Error() string { return e.msg }
func (e *wrappedModelError) Unwrap() error { return ErrIllegalModel }
