package sched

import (
	"fmt"
	"time"
)

const (
	MinSecond = 0
	MaxSecond = 59

	MinMinute = 0
	MaxMinute = 59

	MinHour = 0
	MaxHour = 23

	MinDom = 1
	MaxDom = 31

	MinMonth = 1
	MaxMonth = 12

	MinDow = 0
	MaxDow = 7

	MinYear = 0
	MaxYear = 1<<31 - 1 //grabbed from math.MaxInt32
)

const (
	Yearly       = "@yearly"
	YearlyFormat = "0 0 0 1 1 * *"

	Annually       = "@annually"
	AnnuallyFormat = YearlyFormat

	Monthly       = "@monthly"
	MonthlyFormat = "0 0 0 1 * * *"

	Weekly       = "@weekly"
	WeeklyFormat = "0 0 0 * * 0 *"

	Daily       = "@daily"
	DailyFormat = "0 0 0 * * * *"

	Hourly       = "@hourly"
	HourlyFormat = "0 0 * * * * *"

	Minutely       = "@minutely"
	MinutelyFormat = "0 * * * * * *"

	Secondly       = "@secondly"
	SecondlyFormat = "* * * * * * *"
)

const invalidValue = -1

type Schedule interface {
	NextTime(from time.Time) (time.Time, bool)
	Expression() string
}

type schedule struct {
	second fieldNexter
	minute fieldNexter
	hour   fieldNexter
	dom    dateFieldNexter
	month  fieldNexter
	dow    dateFieldNexter
	year   fieldNexter
}

func newSchedule() *schedule {
	return &schedule{}
}

func (s *schedule) setNexter(nexter interface{}, fi fieldIndex) {
	if nexter == nil {
		panic("nexter cannot be nil")
	}
}

func (s *schedule) NextTime(from time.Time) (time.Time, bool) {
	return from.Add(1), true
}

func (s *schedule) String() string {
	return s.Expression()
}

func (s *schedule) Expression() string {
	return ""
}

type fieldIndex int

const (
	second fieldIndex = iota
	minute
	hour
	dom
	month
	dow
	year
	fieldCount //this is not an actual field index value. just used as a count.
)

var fieldNames = [...]string{"second", "minute", "hour", "day of month", "month", "day of week", "year"}

func (fi fieldIndex) String() string {
	if fi >= 0 && fi < fieldCount {
		return fieldNames[fi]
	}
	return ""
}

func (fi fieldIndex) isInRange(value int) bool {
	fr := fi.fieldRange()
	if fr == nil {
		return false
	}
	return fr.isInRange(value)
}

func (fi fieldIndex) rangeString() string {
	fr := fi.fieldRange()
	if fr == nil {
		return ""
	}
	min, max := fr.min, fr.max
	if fi == dow {
		max = max - 1
	}
	return fmt.Sprintf("%v%v%v", min, Hyphen, max)
}

func (fi fieldIndex) fieldRange() *fieldRange {
	if fi >= 0 && fi < fieldCount {
		return fieldRanges[fi]
	}
	return nil
}

func (fi fieldIndex) canHaveHash() bool {
	return fi == dow
}

func (fi fieldIndex) canHaveWeekday() bool {
	return fi == dom
}

func (fi fieldIndex) isDateField() bool {
	return fi == dom || fi == dow
}

func (fi fieldIndex) modifiers() string {
	result := ""
	if fi.isDateField() {
		result += Last
	}
	if fi.canHaveHash() {
		result += Hash
	}
	if fi.canHaveWeekday() {
		result += Weekday
	}
	return result
}

type fieldRange struct {
	min int
	max int
}

func (fr *fieldRange) isInRange(value int) bool {
	return value >= fr.min && value <= fr.max
}

var fieldRanges = [...]*fieldRange{
	{MinSecond, MaxSecond},
	{MinMinute, MaxMinute},
	{MinHour, MaxHour},
	{MinDom, MaxDom},
	{MinMonth, MaxMonth},
	{MinDow, MaxDow},
	{MinYear, MaxYear},
}
