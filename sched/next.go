package sched

import "time"

type fieldNexter interface {
	next(int) (int, bool)
}

type valueNexter int

func newValueNexter(value int) valueNexter {
	return valueNexter(value)
}

func (vn valueNexter) next(now int) (int, bool) {
	return now, true
}

type anyNexter struct {
	*rangeNexter
}

func newAnyNexter(min, max int) *anyNexter {
	return &anyNexter{
		rangeNexter: newRangeNexter(min, max),
	}
}

type rangeDivNexter struct {
	*rangeNexter
	inc int
}

func newRangeDivNexter(rn *rangeNexter, inc int) *rangeDivNexter {
	return &rangeDivNexter{
		rangeNexter: rn,
		inc:         inc,
	}
}

func (rdn *rangeDivNexter) next(now int) (int, bool) {
	if now < rdn.min {
		return rdn.min, false
	}
	value := now - rdn.min
	result := rdn.min + value + (rdn.inc - (value % rdn.inc))
	if result > rdn.max {
		return rdn.min, true
	}
	return result, false
}

type rangeNexter struct {
	min int
	max int
}

func newRangeNexter(min, max int) *rangeNexter {
	return &rangeNexter{
		min: min,
		max: max,
	}
}

func (rn *rangeNexter) next(now int) (int, bool) {
	if now < rn.min {
		return rn.min, false
	}
	result := now + 1
	if result > rn.max {
		return rn.min, true
	}
	return result, false
}

type multiNexter []fieldNexter

func newMultiNexter(fns ...fieldNexter) multiNexter {
	return multiNexter(fns)
}

func (mn multiNexter) next(now int) (int, bool) {
	return now, true
}

type dateFieldNexter interface {
	next(now int, time time.Time) (int, bool)
}

type multiDateFieldNexter []dateFieldNexter

func (mdn multiDateFieldNexter) next(now int, time time.Time) (int, bool) {
	return now, true
}

type domFieldNexter struct {
	fieldNexter
	isLast    bool
	isWeekday bool
}

type dowFieldNexter struct {
	fieldNexter
	isLast bool
	number int
}
