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

type FormatStringer interface {
	FormatString() string
}

type Schedule interface {
	NextTime(from time.Time) (time.Time, bool)
	FormatString() string
}

const (
	//indices into a schedule.fields array.
	second = iota
	minute
	hour
	dom
	month
	dow
	year
	fieldCount
)

type schedule struct {
	fields [fieldCount]fieldNexter
}

func (s *schedule) NextTime(from time.Time) (time.Time, bool) {
	return from.Add(1), true
}

func (s *schedule) String() string {
	return fmt.Sprintf("%p", s)
}
