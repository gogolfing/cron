package sched

import "testing"

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
