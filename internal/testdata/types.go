package testdata

// TypeKeys returns type keys in ascending alphabetical order
func TypeKeys() []string {
	return []string{
		"added",
		"changed",
		"fixed",
	}
}

// Types returns types map
func Types() map[string]string {
	return map[string]string{
		"added":   "New Feature",
		"fixed":   "Bug Fix",
		"changed": "Feature Change",
	}
}
