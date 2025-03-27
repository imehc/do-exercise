package util

import (
	"testing"
	"time"
)

func TestParseDurationString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "valid days",
			input:   "2d",
			want:    2 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "valid hours",
			input:   "24h",
			want:    24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "valid minutes",
			input:   "30m",
			want:    30 * time.Minute,
			wantErr: false,
		},
		{
			name:    "valid seconds",
			input:   "45s",
			want:    45 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid number only (seconds)",
			input:   "60",
			want:    60 * time.Second,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid unit",
			input:   "24x",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDurationString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDurationString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDurationString() = %v, want %v", got, tt.want)
			}
		})
	}
}