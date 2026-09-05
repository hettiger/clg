package changelog

import "slices"

type Type struct {
	Label   string
	Keyword string
}

var supportedTypes = [...]Type{
	{Label: "New Feature", Keyword: "added"},
	{Label: "Bug Fix", Keyword: "fixed"},
	{Label: "Hotfix", Keyword: "hotfix"},
	{Label: "Feature Change", Keyword: "changed"},
	{Label: "New Deprecation", Keyword: "deprecated"},
	{Label: "Feature Removal", Keyword: "removed"},
	{Label: "Security Fix", Keyword: "security"},
	{Label: "Performance Improvement", Keyword: "performance"},
	{Label: "Other", Keyword: "other"},
	{Label: "No Changelog", Keyword: "ignore"},
}

func SupportedTypes() []Type {
	return slices.Clone(supportedTypes[:])
}

func SupportedTypeKeywords() []string {
	keywords := make([]string, len(supportedTypes))
	for i, t := range supportedTypes {
		keywords[i] = t.Keyword
	}
	return keywords
}
