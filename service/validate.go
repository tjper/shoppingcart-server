package service

import (
	"fmt"

	"github.com/pkg/errors"
)

// validateInt64 validates that the val is not empty. On failure, an error
// is appended to the err argument.
func validateInt64(err error, field string, val int64) {
	if val == 0 {
		var errMsg = fmt.Sprintf(
			"failed to ValidateInt64\tfield=%s\tval=%v",
			field,
			val)
		if err != nil {
			err = errors.Wrap(err, errMsg)
		} else {
			err = errors.New(errMsg)
		}
	}
}
