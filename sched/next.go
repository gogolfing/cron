package sched

import (
	"fmt"
	"strings"
)

type fieldNexter interface {
	next(int) (int, bool)
	formatStringer
}

type valueNexter int

func (vn valueNexter) next(now int) (int, bool) {
	return now, true
}

func (vn valueNexter) formatString() string {
	return fmt.Sprint(vn)
}

type anyNexter struct {
	*rangeNexter
}

func (an *anyNexter) formatString() string {
	return Asterisk
}

type rangeDivNexter struct {
	*rangeNexter
	inc int
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

func (rdn *rangeDivNexter) formatString() string {
	return fmt.Sprintf("%v%v%v", rdn.rangeNexter.formatString(), Slash, rdn.inc)
}

type rangeNexter struct {
	min int
	max int
}

func (rn *rangeNexter) next(now int) (int, bool) {
	result := now + 1
	if result > rn.max {
		return rn.min, true
	}
	return result, false
}

func (rn *rangeNexter) formatString() string {
	return fmt.Sprintf("%v%v%v", rn.min, Hyphen, rn.max)
}

type multiNexter []fieldNexter

func (mn multiNexter) next(now int) (int, bool) {
	return now, true
}

func (mn multiNexter) formatString() string {
	values := make([]string, len(mn))
	for i, field := range mn {
		values[i] = field.formatString()
	}
	return strings.Join(values, Comma)
}
