package support_test

import (
	"testing"

	"github.com/hettiger/clg/internal/support"
	"github.com/stretchr/testify/require"
)

func TestSortedMapKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  []string
	}{
		{
			name:  "empty",
			input: map[string]string{},
			want:  []string{},
		},
		{
			name: "in order",
			input: map[string]string{
				"a": "a value",
				"b": "b value",
				"c": "c value",
			},
			want: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			name: "out of order",
			input: map[string]string{
				"c": "c value",
				"a": "a value",
				"b": "b value",
			},
			want: []string{
				"a",
				"b",
				"c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := support.SortedMapKeys(tt.input)

			require.Equal(t, tt.want, got)
		})
	}
}
