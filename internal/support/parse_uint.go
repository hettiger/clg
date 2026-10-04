package support

import (
	"strconv"
	"strings"
)

// ParseUint parses ASCII decimal digits into a uint, trimming Unicode whitespace.
// It returns an error for invalid input or overflow.
func ParseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, strconv.IntSize)
	if err != nil {
		return 0, err
	}

	return uint(n), nil
}
