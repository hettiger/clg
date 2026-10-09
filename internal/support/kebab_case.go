package support

import (
	"regexp"
	"strings"
)

var (
	regexpKebabCaseAcronymBoundary = regexp.MustCompile(`([A-Z]+)([A-Z](?:[a-z]|[^\x00-\x7F]))`) // [^\x00-\x7F] = non ASCII chars
	regexpKebabCaseWordBoundary    = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	regexpNonAlphanumeric          = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

// KebabCase converts text to lowercase, hyphen-separated words, e.g. "HTTPServer"
// becomes "http-server". Only ASCII letters and digits are kept; the result may be empty.
func KebabCase(str string) string {
	str = regexpKebabCaseAcronymBoundary.ReplaceAllString(str, "$1-$2")
	str = regexpKebabCaseWordBoundary.ReplaceAllString(str, "$1-$2")
	str = regexpNonAlphanumeric.ReplaceAllString(str, "-")
	str = strings.Trim(str, "-")
	str = strings.ToLower(str)

	return str
}
