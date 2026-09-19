package config

func defaultTypes() map[string]string {
	return map[string]string{
		"added":       "New Feature",
		"fixed":       "Bug Fix",
		"hotfix":      "Hotfix",
		"changed":     "Feature Change",
		"deprecated":  "New Deprecation",
		"removed":     "Feature Removal",
		"security":    "Security Fix",
		"performance": "Performance Improvement",
		"other":       "Other",
	}
}
