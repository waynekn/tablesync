package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToMatrix(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]string
		rowLen  int
		want    [][]string
		wantErr bool
	}{
		{
			name: "valid input",
			input: map[string]string{
				"0:0": "A",
				"0:1": "B",
				"1:0": "C",
				"1:1": "D",
			},
			rowLen: 2,
			want: [][]string{
				{"A", "B"},
				{"C", "D"},
			},
			wantErr: false,
		},
		{
			name: "invalid coordinates",
			input: map[string]string{
				"0:0": "A",
				"1:2": "B", // Invalid column index
			},
			rowLen: 2,
			// invalid column index will cause the second row to be skipped
			want: [][]string{
				{"A", ""},
			},
			wantErr: false,
		},
		{
			name: "invalid key format",
			input: map[string]string{
				"z:a": "A",
				"1:2": "B", // Invalid column index
			},
			want:    nil,
			rowLen:  2,
			wantErr: true, // Expecting an error due to invalid key format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapToMatrix(tt.input, tt.rowLen)

			if tt.wantErr {
				assert.Error(t, err, "mapToMatrix should return an error")
			} else {
				assert.NoError(t, err, "mapToMatrix should not return an error")
			}
			assert.Equal(t, tt.want, got, "mapToMatrix result mismatch")
		})
	}
}

func TestCoordsFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantRow int
		wantCol int
		wantErr bool
	}{
		{
			name:    "valid coordinates",
			input:   "1:2",
			wantRow: 1,
			wantCol: 2,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "1-2",
			wantRow: -1,
			wantCol: -1,
			wantErr: true,
		},
		{
			name:    "invalid row index",
			input:   "a:b",
			wantRow: -1,
			wantCol: -1,
			wantErr: true,
		},
		{
			name:    "negative indices",
			input:   "-5:-10",
			wantRow: -1,
			wantCol: -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRow, gotCol, err := coordsFromString(tt.input)

			if tt.wantErr {
				assert.Error(t, err, "coordsFromString should return an error")
			} else {
				assert.NoError(t, err, "coordsFromString should not return an error")
			}
			assert.Equal(t, tt.wantRow, gotRow, "coordsFromString row mismatch")
			assert.Equal(t, tt.wantCol, gotCol, "coordsFromString column mismatch")
		})
	}
}
