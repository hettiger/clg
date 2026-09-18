package changelog_test

import "github.com/hettiger/clg/internal/support"

func testTypeKeys() []string {
	return support.SortedMapKeys(testTypes())
}

func testTypes() map[string]string {
	return map[string]string{
		"added":   "New Feature",
		"fixed":   "Bug Fix",
		"changed": "Feature Change",
		"ignore":  "No Changelog",
	}
}
