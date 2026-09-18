package config

func defaultTypes() Types {
	return Types{
		"added":       "New Feature",
		"fixed":       "Bug Fix",
		"hotfix":      "Hotfix",
		"changed":     "Feature Change",
		"deprecated":  "New Deprecation",
		"removed":     "Feature Removal",
		"security":    "Security Fix",
		"performance": "Performance Improvement",
		"other":       "Other",
		"ignore":      "No Changelog",
	}
}
