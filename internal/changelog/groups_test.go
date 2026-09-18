package changelog_test

import "github.com/hettiger/clg/internal/support"

func testGroupKeys() []string {
	return support.SortedMapKeys(testGroups())
}

func testGroups() map[string]string {
	return map[string]string{
		"front": "Frontend",
		"back":  "Backend",
	}
}
