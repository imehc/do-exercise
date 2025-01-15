package utils

import (
	"testing"
	"time"
)

func TestParseDurationString_ValidFormat_ReturnsCorrectDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"1d", 24 * time.Hour},
		{"2h", 2 * time.Hour},
		{"30m", 30 * time.Minute},
		{"5s", 5 * time.Second},
	}

	for _, test := range tests {
		duration, err := ParseDurationString(test.input)
		if err != nil {
			t.Errorf("ParseDurationString(%q) returned an error: %v", test.input, err)
		}
		if duration != test.expected {
			t.Errorf("ParseDurationString(%q) = %v, want %v", test.input, duration, test.expected)
		}
	}
}

func TestParseDurationString_InvalidUnit_ReturnsError(t *testing.T) {
	input := "1x"
	_, err := ParseDurationString(input)
	if err == nil {
		t.Errorf("ParseDurationString(%q) should have returned an error", input)
	}
}

func TestParseDurationString_PureNumericInput_ReturnsSeconds(t *testing.T) {
	input := "100"
	expected := 100 * time.Second
	duration, err := ParseDurationString(input)
	if err != nil {
		t.Errorf("ParseDurationString(%q) returned an error: %v", input, err)
	}
	if duration != expected {
		t.Errorf("ParseDurationString(%q) = %v, want %v", input, duration, expected)
	}
}

func TestParseDurationString_InvalidFormat_ReturnsError(t *testing.T) {
	input := "abc"
	_, err := ParseDurationString(input)
	if err == nil {
		t.Errorf("ParseDurationString(%q) should have returned an error", input)
	}
}
