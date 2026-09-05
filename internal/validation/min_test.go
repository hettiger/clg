package validation_test

import (
	"testing"

	"github.com/hettiger/clg/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMin(t *testing.T) {
	tests := []struct {
		description string
		name        string
		value       string
		min         int
		wantError   bool
		wantMessage string
	}{
		{
			description: "exact size",
			name:        "Field",
			value:       "foo",
			min:         3,
			wantError:   false,
		},
		{
			description: "bigger size",
			name:        "Field",
			value:       "foo",
			min:         2,
			wantError:   false,
		},
		{
			description: "smaller size",
			name:        "Field",
			value:       "foo",
			min:         4,
			wantError:   true,
			wantMessage: `Field must have a minimum of 4 characters (got "foo")`,
		},
		{
			description: "empty input",
			name:        "Field",
			value:       "",
			min:         4,
			wantError:   true,
			wantMessage: `Field must have a minimum of 4 characters (got "")`,
		},
		{
			description: "negative min",
			name:        "Field",
			value:       "foo",
			min:         -1,
			wantError:   false,
		},
		{
			description: "zero min",
			name:        "Field",
			value:       "foo",
			min:         0,
			wantError:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateMin(tt.name, tt.value, tt.min)

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
