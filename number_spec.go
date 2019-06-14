package cron

import (
	"errors"
	"fmt"
)

type numberSpec struct {
	max int

	values []int
	ranges []rangeSpec
}

type RangeError struct {
	Min  int
	Max  int
	Step int

	Message string
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("cron: invalid range %v-%v/%v : %v", e.Min, e.Max, e.Step, e.Message)
}

//rangeSpec is a composite type that holds the values of a range.
//rangeSpec[0] is the minimum value, rangeSpec[1] is the maximum value, and
//rangeSpec[2] is the step value.
type rangeSpec [3]int

func newRangeSpec(min, max, step int) (rangeSpec, error) {
	//TODO
	//
	//check for minMin
	//maxMan
	//max < min
	//step <= 0
	if min < 0 || max < 0 || step < 0 {
		return [3]int{}, errors.New("cron: range values must be ")
	}
	return [3]int{min, max, step}, nil
}

func (rs rangeSpec) NextAfter(value int) (result int, overflow bool) {
	if value < rs.min() {
		result = rs.min()
		return
	}

	diff := value - rs.min()
	if diff < 0 {
		result = rs.min()
		return
	}

	result = value + (rs.step() - ((value - rs.min()) % rs.step()))

	overflow = result > rs.max()
	if overflow {
		result = rs.min()
	}

	return
}

func (rs rangeSpec) min() int {
	return rs[0]
}

func (rs rangeSpec) max() int {
	return rs[1]
}

func (rs rangeSpec) step() int {
	return rs[2]
}
