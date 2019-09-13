package service

import (
	"github.com/pkg/errors"
)

type validate struct {
	Err error
}

// validate runs a set of checks in order to ensure the data is as expected.
func (v *validate) check(field string, checks ...func() error) {
	for _, check := range checks {
		if err := check(); err != nil {
			v.Err = err
			return
		}
	}
}

func intNotEmpty(val int) func() error {
	return func() error {
		if val == 0 {
			return errors.Errorf("failed to intNotEmpty\tval=%v", val)
		}
		return nil
	}
}

func intGreaterThan(val int, min int) func() error {
	return func() error {
		if val <= min {
			return errors.Errorf("failed to validateGreaterThan\tval=%v\tmin=%v", val, min)
		}
		return nil
	}
}
