package utils

import (
	"fmt"
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
		fmt.Println("expected:", test.expected)
		fmt.Println("actual:", result)
		if result != test.expected {
			t.Errorf("Percentile(%v, %v) = %v; expected %v", test.input, test.percentile, result, test.expected)
		}
	}
}
