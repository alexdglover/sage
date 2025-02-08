package utils

import (
	"testing"
)

func TestPercentile(t *testing.T) {
	tests := []struct {
		input      []int
		percentile float64
		expected   int
	}{
		{[]int{1, 2, 3, 4, 5}, 0.01, 1},
		{[]int{1, 2, 3, 4, 5}, 0.15, 1},
		{[]int{1, 2, 3, 4, 5}, 0.2, 1},
		{[]int{1, 2, 3, 4, 5}, 0.21, 2},
		{[]int{1, 2, 3, 4, 5}, 0.25, 2},
		{[]int{1, 2, 3, 4, 5}, 0.40, 2},
		{[]int{1, 2, 3, 4, 5}, 0.41, 3},
		{[]int{1, 2, 3, 4, 5}, 0.5, 3},
		{[]int{1, 2, 3, 4, 5}, 0.75, 4},
		{[]int{1, 2, 3, 4, 5}, 0.9, 5},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 0.25, 3},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 0.75, 9},
		{[]int{}, 0.5, 0},
	}

	for _, test := range tests {
		result := Percentile(test.input, test.percentile)
		if result != test.expected {
			t.Errorf("Percentile(%v, %v) = %v; expected %v", test.input, test.percentile, result, test.expected)
		}
	}
}

func TestCentsToDollarStringHumanized(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"1", 1, "0.01"},
		{"10", 10, "0.10"},
		{"100", 100, "1.00"},
		{"1000", 1000, "10.00"},
		{"10000", 10000, "100.00"},
		{"100000", 100000, "1,000.00"},
		{"1000000", 1000000, "10,000.00"},
		{"90", 90, "0.90"},
		{"115", 115, "1.15"},
		{"107", 107, "1.07"},
		{"14.5", 145, "14.5"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := CentsToDollarStringHumanized(test.input)
			if result != test.expected {
				t.Errorf("CentsToDollarStringHumanized(%v) = %s; expected %s", test.input, result, test.expected)
			}
		})
	}

}

func TestTrimTrailingDigits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		inputString string
		inputDigits int
		expected    string
	}{
		{"1.00", 2, "1.00"},
		{"0.90", 2, "0.90"},
		{"1.115", 2, "1.11"},
		{"1.107", 2, "1.10"},
		{"2.1", 2, "2.10"},
		{"2", 2, "2.00"},
		{"2", 10, "2.0000000000"},
		{"1", 0, "1"},
		{"1", -1, "1"},
	}

	for _, test := range tests {
		t.Run(test.inputString, func(t *testing.T) {
			t.Parallel()
			result := trimTrailingDigits(test.inputString, test.inputDigits)
			if result != test.expected {
				t.Errorf("trimTrailingDigits(%s, %d) = %s; expected %s", test.inputString, test.inputDigits, result, test.expected)
			}
		})
	}

}
