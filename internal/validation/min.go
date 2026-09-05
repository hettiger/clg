package validation

import (
	"fmt"
	"unicode/utf8"
)

func ValidateMin(name, value string, min int) error {
	if utf8.RuneCountInString(value) >= min {
		return nil
	}

	return fmt.Errorf(
		`%s must have a minimum of %v characters (got %q)`,
		name,
		min,
		value,
	)
}
