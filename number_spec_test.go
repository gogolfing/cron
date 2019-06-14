package cron

import "testing"

func TestNewRangeSpec_ErrorsWithIncorrectMinMaxStep(t *testing.T) {

}

func TestRangeSpec_NextValue(t *testing.T) {
	cases := []struct {
		min      int
		max      int
		step     int
		value    int
		result   int
		overflow bool
	}{
		{0, 59, 1, -1234, 0, false},
		{0, 59, 1, -1, 0, false},
		{0, 59, 1, 0, 1, false},
		{0, 59, 1, 32, 33, false},
		{0, 59, 1, 58, 59, false},
		{0, 59, 1, 59, 0, true},
		{0, 59, 1, 1234, 0, true},

		{1, 100, 5, 0, 1, false},
		{1, 100, 5, 1, 6, false},
		{1, 100, 5, 2, 6, false},
		{1, 100, 5, 3, 6, false},
		{1, 100, 5, 4, 6, false},
		{1, 100, 5, 5, 6, false},
		{1, 100, 5, 6, 11, false},
		{1, 100, 5, 96, 1, true},

		{50, 50, 10, 49, 50, false},
		{50, 50, 10, 50, 50, true},
		{50, 50, 10, 51, 50, true},

		{17, 30, 3, 16, 17, false},
		{17, 30, 3, 17, 20, false},
		{17, 30, 3, 21, 23, false},
		{17, 30, 3, 28, 29, false},
		{17, 30, 3, 29, 17, true},
	}

	for i, tc := range cases {
		rs, _ := newRangeSpec(tc.min, tc.max, tc.step)

		result, overflow := rs.NextAfter(tc.value)

		if result != tc.result || overflow != tc.overflow {
			t.Errorf("%d: %v = %v %v WANT %v %v", i, tc, result, overflow, tc.result, tc.overflow)
		}
	}
}
