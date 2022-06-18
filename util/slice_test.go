package util_test

import (
	"testing"

	"github.com/yzx9/otodo/util"
)

func TestMap(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		in       []int
		expected []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{1, 2, 3}, []int{2, 4, 6}},
	}

	double := func(in int) int {
		return in * 2
	}

	for _, tt := range tests {
		actual := util.Map(double, tt.in)
		if !LoopCompare(actual, tt.expected) {
			t.Errorf("Map(double, %v) = %v; expected: %v", tt.in, actual, tt.expected)
		}
	}
}

func LoopCompare[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
