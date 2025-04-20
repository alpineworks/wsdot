package ferries

import (
	"testing"
)

func TestParseOffsetMilli(t *testing.T) {
	tests := []struct {
		name      string
		offset    string
		want      int
		expectErr bool
	}{
		{
			name:      "Valid positive offset",
			offset:    "+0700",
			want:      25200000, // 7 hours in milliseconds
			expectErr: false,
		},
		{
			name:      "Valid negative offset",
			offset:    "-0530",
			want:      -19800000, // 5 hours 30 minutes in milliseconds
			expectErr: false,
		},
		{
			name:      "Invalid offset length",
			offset:    "+070",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid sign",
			offset:    "*0700",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid hour value",
			offset:    "+2500",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid minute value",
			offset:    "+0760",
			want:      0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOffsetMilli(tt.offset)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseOffsetMilli() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseOffsetMilli() = %v, want %v", got, tt.want)
			}
		})
	}
}
