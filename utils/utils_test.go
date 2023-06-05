package utils_test

import (
	"testing"

	"github.com/bnema/gobatterywatcher/utils"
)

func TestConvertToWatts(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10 W", 10},
		{"100 mW", 0.1},
		{"500 uW", 0.0005},
		{"0.5 W", 0.5},
		{"0.05 mW", 0.00005},
		{"0.001 W", 0.001},
		{"invalid", 0},
	}

	for _, test := range tests {
		result := utils.ConvertToWatts(test.input)
		if result != test.expected {
			t.Errorf("unexpected result for input %q: got %v, want %v", test.input, result, test.expected)
		}
	}
}
func TestReadPowerTop(t *testing.T) {
	data, err := utils.ReadPowerTop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if there are any entries with 0.0% power usage
	for _, d := range data {
		if d.PowerUsage == 0 {
			t.Errorf("unexpected entry with 0.0%% power usage: %v", d)
		}
	}
}
