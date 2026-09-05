package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateIn(t *testing.T) {
	tests := []struct {
		description   string
		name          string
		value         string
		allowedValues []string
		wantError     bool
		wantMessage   string
	}{
		{
			description:   "missing",
			name:          "Field",
			value:         "foo",
			allowedValues: []string{"bar", "baz"},
			wantError:     true,
			wantMessage:   `Field must be one of the following values: "bar", "baz" (got "foo")`,
		},
		{
			description:   "present",
			name:          "Field",
			value:         "foo",
			allowedValues: []string{"bar", "foo", "baz"},
			wantError:     false,
		},
		{
			description:   "empty input",
			name:          "Field",
			value:         "",
			allowedValues: []string{"foo", "bar"},
			wantError:     true,
			wantMessage:   `Field must be one of the following values: "foo", "bar" (got "")`,
		},
		{
			description:   "empty allowed values",
			name:          "Field",
			value:         "foo",
			allowedValues: []string{},
			wantError:     true,
			wantMessage:   `Field must be one of the following values: – (got "foo")`,
		},
		{
			description:   "values requiring escaping",
			name:          "Field",
			value:         `a"b`,
			allowedValues: []string{"foo", `x\y`},
			wantError:     true,
			wantMessage:   `Field must be one of the following values: "foo", "x\\y" (got "a\"b")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateIn(tt.name, tt.value, tt.allowedValues...)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantMessage != "" {
				assert.EqualError(t, err, tt.wantMessage)
			}
		})
	}
}
