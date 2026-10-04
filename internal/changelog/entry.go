package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hettiger/clg/internal/validation"
	"go.yaml.in/yaml/v3"
)

type ChangelogEntryFile struct {
	Entry ChangelogEntry
	Path  string
}

type ChangelogEntry struct {
	Group  string `yaml:"group,omitempty"`
	Type   string `yaml:"type,omitempty"`
	Title  string `yaml:"title,omitempty"`
	Author string `yaml:"author,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

func NewChangelogEntry(YAMLData []byte, groupKeys, typeKeys []string) (ChangelogEntry, error) {
	var entry ChangelogEntry
	if err := yaml.Unmarshal(YAMLData, &entry); err != nil {
		return entry, err
	}
	if err := entry.Validate(groupKeys, typeKeys); err != nil {
		return ChangelogEntry{}, err
	}
	return entry, nil
}

func (e ChangelogEntry) YAMLData() ([]byte, error) {
	return yaml.Marshal(e)
}

func (e ChangelogEntry) YAML() (string, error) {
	data, err := e.YAMLData()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e ChangelogEntry) Validate(groupKeys, typeKeys []string) error {
	if err := validation.ValidateMin("title", e.Title, 1); err != nil {
		return err
	}

	if len(groupKeys) == 0 {
		if e.Group != "" {
			return fmt.Errorf("Unsupported group (%s)", e.Group)
		}
	} else if err := validation.ValidateIn(
		"group",
		e.Group,
		groupKeys...,
	); err != nil {
		return err
	}

	if err := validation.ValidateIn(
		"type",
		e.Type,
		typeKeys...,
	); err != nil {
		return err
	}

	return nil
}

func (e ChangelogEntry) Filename(uuidV7 func() (string, error)) (string, error) {
	id, err := uuidV7()
	if err != nil {
		return "", err
	}

	var filename strings.Builder

	if e.Group != "" {
		fmt.Fprintf(&filename, "%s-", e.Group)
	}

	fmt.Fprintf(&filename, "%s-%s.yml", e.Type, id)

	return filename.String(), nil
}

var (
	ErrIssuePatternInvalid              = errors.New("the pattern is invalid and does not compile")
	ErrIssuePatternCaptureGroupsInvalid = errors.New("the pattern must have exactly one capturing group")
	ErrIssuePatternMismatch             = errors.New("the entry branch does not match the issue pattern")
)

func (e ChangelogEntry) Issue(pattern string) (string, error) {
	if pattern == "" {
		return "", nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("issue pattern %q: %w: %w", pattern, ErrIssuePatternInvalid, err)
	}

	if re.NumSubexp() != 1 {
		return "", fmt.Errorf("issue pattern %q: %w", pattern, ErrIssuePatternCaptureGroupsInvalid)
	}

	matches := re.FindStringSubmatch(e.Branch)

	if matches == nil {
		return "", fmt.Errorf("issue pattern %q for branch %q: %w", pattern, e.Branch, ErrIssuePatternMismatch)
	}

	return matches[1], nil
}
