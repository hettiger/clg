package testdata

// GroupKeys returns group keys in ascending alphabetical order
func GroupKeys() []string {
	return []string{
		"back",
		"front",
	}
}

// Groups returns groups map
func Groups() map[string]string {
	return map[string]string{
		"front": "Frontend",
		"back":  "Backend",
	}
}
