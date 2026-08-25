package model

import "errors"

var (
	ErrIllegalModel = errors.New("illegal beam model")

	ErrNoConvergence = errors.New("deflection solve did not converge")
)

func IllegalModel(msg string) error {
	return &wrappedModelError{msg: msg}
}

type wrappedModelError struct {
	msg string
}

func (e *wrappedModelError) Error() string { return e.msg }
func (e *wrappedModelError) Unwrap() error { return ErrIllegalModel }
