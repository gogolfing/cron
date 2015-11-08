package sched

import (
	"testing"
	"time"
)

func TestNewIntervalSchedule(t *testing.T) {
	s := NewIntervalSchedule(time.Duration(1))
	if time.Duration(s.(IntervalSchedule)) != time.Duration(1) {
		t.Errorf("NewIntervalSchedule() failed to cast duration correctly")
	}
}

func TestIntervalSchedule_NextTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		duration time.Duration
		from     time.Time
		after    bool
	}{
		{time.Hour, now, true},
		{time.Nanosecond, now, true},
		{0, now, false},
		{-1, now, false},
	}
	for _, test := range tests {
		s := NewIntervalSchedule(test.duration)
		wantNext := test.from.Add(test.duration)
		next, after := s.NextTime(test.from)
		if !next.Equal(wantNext) || after != test.after {
			t.Errorf("%v.NextTime(%v) = %v, %v WANT %v, %v", s, test.from, next, after, wantNext, test.after)
		}
	}
}

func TestIntervalSchedule_String(t *testing.T) {
	s := NewIntervalSchedule(time.Hour)
	want := "sched.IntervalSchedule(" + time.Hour.String() + ")"
	if result := s.(IntervalSchedule).String(); result != want {
		t.Errorf("%v.String() = %v WANT %v", s, result, want)
	}
}

func TestIntervalSchedule_Expression(t *testing.T) {
	s := NewIntervalSchedule(time.Minute)
	want := Every + " " + time.Minute.String()
	if result := s.Expression(); result != want {
		t.Errorf("%v.Expression() = %v WANT %v", s, result, want)
	}
}

func TestFieldIndex_String(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result string
	}{
		{-1, ""},
		{second, fieldNames[second]},
		{minute, fieldNames[minute]},
		{hour, fieldNames[hour]},
		{dom, fieldNames[dom]},
		{month, fieldNames[month]},
		{dow, fieldNames[dow]},
		{year, fieldNames[year]},
		{fieldCount, ""},
	}
	for _, test := range tests {
		if result := test.fi.String(); result != test.result {
			t.Errorf("fieldIndex(%v).String() = %v WANT %v", int(test.fi), result, test.result)
		}
	}
}

func TestFieldIndex_isInRange(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  int
		result bool
	}{
		{-1, 1, false},
		{second, MinSecond, true},
		{second, MaxSecond, true},
		{second, MinSecond - 1, false},
		{second, MaxSecond + 1, false},
		{minute, MinMinute, true},
		{minute, MaxMinute, true},
		{minute, MinMinute - 1, false},
		{minute, MaxMinute + 1, false},
		{hour, MinHour, true},
		{hour, MaxHour, true},
		{hour, MinHour - 1, false},
		{hour, MaxHour + 1, false},
		{dom, MinDom, true},
		{dom, MaxDom, true},
		{dom, MinDom - 1, false},
		{dom, MaxDom + 1, false},
		{month, MinMonth, true},
		{month, MaxMonth, true},
		{month, MinMonth - 1, false},
		{month, MaxMonth + 1, false},
		{dow, MinDow, true},
		{dow, MaxDow, true},
		{dow, MinDow - 1, false},
		{dow, MaxDow + 1, false},
		{year, MinYear, true},
		{year, MaxYear, true},
		{year, MinYear - 1, false},
		{year, MaxYear + 1, false},
		{fieldCount, 1, false},
	}
	for _, test := range tests {
		if result := test.fi.isInRange(test.value); result != test.result {
			t.Errorf("%v.isInRange(%v) = %v WANT %v", test.fi, test.value, result, test.result)
		}
	}
}

func TestFieldIndex_rangeString(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result string
	}{
		{-1, ""},
		{second, "0-59"},
		{minute, "0-59"},
		{hour, "0-23"},
		{dom, "1-31"},
		{month, "1-12"},
		{dow, "0-6"},
		{year, "0-2147483647"},
		{fieldCount, ""},
	}
	for _, test := range tests {
		if result := test.fi.rangeString(); result != test.result {
			t.Errorf("%v.rangeString() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldIndex_fieldRange(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result *fieldRange
	}{
		{-1, nil},
		{second, fieldRanges[second]},
		{minute, fieldRanges[minute]},
		{hour, fieldRanges[hour]},
		{dom, fieldRanges[dom]},
		{month, fieldRanges[month]},
		{dow, fieldRanges[dow]},
		{year, fieldRanges[year]},
		{fieldCount, nil},
	}
	for _, test := range tests {
		if result := test.fi.fieldRange(); result != test.result {
			t.Errorf("%v.fieldRange() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldIndex_canHaveHash(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result bool
	}{
		{-1, false},
		{second, false},
		{dow, true},
		{fieldCount, false},
	}
	for _, test := range tests {
		if result := test.fi.canHaveHash(); result != test.result {
			t.Errorf("%v.canHaveHash() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldIndex_canHaveWeekday(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result bool
	}{
		{-1, false},
		{second, false},
		{dom, true},
		{fieldCount, false},
	}
	for _, test := range tests {
		if result := test.fi.canHaveWeekday(); result != test.result {
			t.Errorf("%v.canHaveWeekday() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldIndex_isDateField(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result bool
	}{
		{-1, false},
		{second, false},
		{dom, true},
		{dow, true},
		{fieldCount, false},
	}
	for _, test := range tests {
		if result := test.fi.isDateField(); result != test.result {
			t.Errorf("%v.isDateField() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldIndex_modifiers(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		result string
	}{
		{-1, ""},
		{second, ""},
		{dom, Last + Weekday},
		{dow, Last + Hash},
		{fieldCount, ""},
	}
	for _, test := range tests {
		if result := test.fi.modifiers(); result != test.result {
			t.Errorf("%v.modifiers() = %v WANT %v", test.fi, result, test.result)
		}
	}
}

func TestFieldRange_isInRange(t *testing.T) {
	tests := []struct {
		min    int
		max    int
		value  int
		result bool
	}{
		{0, 10, -1, false},
		{0, 10, 0, true},
		{0, 10, 4, true},
		{0, 10, 10, true},
		{0, 10, 11, false},
		{8, 8, 7, false},
		{8, 8, 8, true},
		{8, 8, 9, false},
		{10, 0, -1, false},
		{10, 0, 0, false},
		{10, 0, 4, false},
		{10, 0, 10, false},
		{10, 0, 11, false},
	}
	for _, test := range tests {
		fr := &fieldRange{
			min: test.min,
			max: test.max,
		}
		result := fr.isInRange(test.value)
		if result != test.result {
			t.Errorf("*fieldRange{%v, %v}.isInRange(%v) = %v WANT %v", test.min, test.max, test.value, result, test.result)
		}
	}
}
