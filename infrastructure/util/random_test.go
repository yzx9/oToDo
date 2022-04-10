package util_test

import (
	"math/rand"
	"testing"

	"github.com/yzx9/otodo/infrastructure/util"
)

func TestRandomString(t *testing.T) {
	t.Parallel()

	rand.Seed(0)

	tests := []struct {
		in       int
		excepted int
	}{
		{-10, 0},
		{-1, 0},
		{0, 0},
		{1, 1},
		{10, 10},
		{100, 100},
		{1000, 1000},
		{10000, 10000},
	}

	for _, i := range tests {
		str := util.RandomString(i.in)
		if length := len(str); length != i.excepted {
			t.Errorf("len(RandomString(%v)) = %v; excepted: %v", i.in, length, i.excepted)
		}
	}
}
