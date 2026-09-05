package validation

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func ValidateIn(name string, value string, allowedValues ...string) error {
	if slices.Contains(allowedValues, value) {
		return nil
	}

	quotedAllowedValues := make([]string, len(allowedValues))
	for i, v := range allowedValues {
		quotedAllowedValues[i] = strconv.Quote(v)
	}

	formattedAllowedValues := strings.Join(quotedAllowedValues, ", ")
	if formattedAllowedValues == "" {
		formattedAllowedValues = "–"
	}

	return fmt.Errorf(
		`%s must be one of the following values: %s (got %q)`,
		name,
		formattedAllowedValues,
		value,
	)
}
