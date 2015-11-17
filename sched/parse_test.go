package sched

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewParseError(t *testing.T) {
	err := newParseError("expression", "description")
	if err.Expression != "expression" || err.Description != "description" {
		t.Fail()
	}
}

func TestParseError_Error(t *testing.T) {
	err := newParseError("expression", "description")
	const want = `sched: could not parse "expression": description`
	if result := err.Error(); result != want {
		t.Errorf("err.Error() = %v WANT %v", result, want)
	}
}

func TestFieldsFunc(t *testing.T) {
	tests := []struct {
		value  rune
		result bool
	}{
		{' ', true},
		{'\t', true},
		{'e', false},
		{'0', false},
	}
	for _, test := range tests {
		if result := fieldSeparatorFunc(test.value); result != test.result {
			t.Errorf("fieldSeparatorFunc(%v) = %v WANT %v", test.value, result, test.result)
		}
	}
}

func TestFields(t *testing.T) {
	tests := []struct {
		exp    string
		result []string
	}{
		{"", []string{}},
		{"expression", []string{"expression"}},
		{"a b", []string{"a", "b"}},
		{"a  b", []string{"a", "b"}},
		{"a\tb", []string{"a", "b"}},
		{"a\t\tb", []string{"a", "b"}},
		{"a\t \tb", []string{"a", "b"}},
		{" a\t \tb\t", []string{"a", "b"}},
		{" a\t \tb\t\n    ", []string{"a", "b"}},
	}
	for _, test := range tests {
		if result := Fields(test.exp); !reflect.DeepEqual(result, test.result) {
			t.Errorf("Fields(%v) = %q WANT %q", test.exp, result, test.result)
		}
	}
}

func TestFieldParts(t *testing.T) {
	tests := []struct {
		field  string
		result []string
	}{
		{"", []string{}},
		{"field", []string{"field"}},
		{"a,b", []string{"a", "b"}},
		{"a,,b", []string{"a", "", "b"}},
		{"a,b,c", []string{"a", "b", "c"}},
		{"this is a string a,this is a string b", []string{"this is a string a", "this is a string b"}},
	}
	for _, test := range tests {
		if result := FieldParts(test.field); !reflect.DeepEqual(result, test.result) {
			t.Errorf("FieldParts(%v) = %q WANT %q", test.field, result, test.result)
		}
	}
}

func TestGetNormalizedFields(t *testing.T) {
	tests := []struct {
		expression string
		result     []string
		hasError   bool
	}{
		{"", nil, true},
		{"expression", nil, true},
		{Secondly, []string{Asterisk, Asterisk, Asterisk, Asterisk, Asterisk, Asterisk, Asterisk}, false},
		{"a b ", []string{"a", "b"}, false},
		{"a b c", nil, true},
		{"a b c d", nil, true},
		{"a b c d e", []string{"0", "a", "b", "c", "d", "e", Asterisk}, false},
		{"a b c d e f", []string{"a", "b", "c", "d", "e", "f", Asterisk}, false},
		{"a b c d e f g", []string{"a", "b", "c", "d", "e", "f", "g"}, false},
		{"a b c d e f g h", nil, true},
	}
	for _, test := range tests {
		result, err := getNormalizedFields(test.expression)
		if (err != nil) != test.hasError {
			t.Errorf("getNormalizedFields(%v) error = %v WANT to have an error %v", test.expression, err, test.hasError)
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("getNormalizedFields(%v) result = %v WANT %v", test.expression, result, test.result)
		}
		if test.result != nil && len(result) != int(fieldCount) && len(result) != 2 {
			t.Errorf("getNormalizedFields(%v) length must equal 2 or %v", test.expression, fieldCount)
		}
	}
}

func TestValidateNumberOfFields(t *testing.T) {
	errString := "number of fields must be 1, 2, 5, 6, or 7"
	tests := []struct {
		fields []string
		count  int
		err    string
	}{
		{nil, 0, errString},
		{[]string{}, 0, errString},
		{make([]string, 1), 1, ""},
		{make([]string, 2), 2, ""},
		{make([]string, 3), 0, errString},
		{make([]string, 4), 0, errString},
		{make([]string, 5), 5, ""},
		{make([]string, 6), 6, ""},
		{make([]string, 7), 7, ""},
		{make([]string, 8), 0, errString},
	}
	for _, test := range tests {
		count, err := validateNumberOfFields(test.fields)
		if (err != nil || test.err != "") && err.Error() != test.err {
			t.Errorf("validateNumberOfFields(%v) error = %v WANT %v", test.fields, err, test.err)
		}
		if count != test.count {
			t.Errorf("validateNumberOfFields(%v) count = %v WANT %v", test.fields, count, test.count)
		}
	}
}

func TestGetNormalizedDirectiveFields(t *testing.T) {
	tests := []struct {
		directive string
		result    []string
		err       string
	}{
		{Yearly, Fields(YearlyFormat), ""},
		{Annually, Fields(AnnuallyFormat), ""},
		{Monthly, Fields(MonthlyFormat), ""},
		{Weekly, Fields(WeeklyFormat), ""},
		{Daily, Fields(DailyFormat), ""},
		{strings.ToUpper(Daily), Fields(DailyFormat), ""},
		{Hourly, Fields(HourlyFormat), ""},
		{Minutely, Fields(MinutelyFormat), ""},
		{Secondly, Fields(SecondlyFormat), ""},
		{"@reboot", nil, `the directive "@reboot" is not recognized`},
	}
	for _, test := range tests {
		result, err := getNormalizedDirectiveFields(test.directive)
		if (err != nil || test.err != "") && err.Error() != test.err {
			t.Errorf("getNormalizedDirectiveFields(%v) error = %v WANT %v", test.directive, err, test.err)
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("getNormalizedDirectiveFields(%v) result = %v WANT %v", test.directive, result, test.result)
		}
	}
}

func TestParseDomDateField(t *testing.T) {
}

func TestParseDowDateField(t *testing.T) {
	tests := []struct {
		modifiers string
		fieldNexter
		resultLast bool
		resultHash int
		err        string
	}{
		{Hash + "6", nil, false, invalidValue, `invalid value for "#" modifier`},
		{Hash + "-1", nil, false, invalidValue, `invalid value for "#" modifier`},
		{Hash + "a", nil, false, invalidValue, `value after "#" ` + errParseInteger.Error()},
		{Last + Hash, nil, false, invalidValue, `cannot have "L" and "#" modifiers together`},
		{Hash + Last, nil, false, invalidValue, `cannot have "L" and "#" modifiers together`},
		{Last + Hash + "2", nil, false, invalidValue, `cannot have "L" and "#" modifiers together`},
		{Last, nil, true, invalidValue, ""},
		{Hash + "2", nil, false, 2, ""},
		{Last, nil, true, invalidValue, ""},
		{"", valueNexter(1), false, invalidValue, ""},
	}
	for _, test := range tests {
		dfn, err := parseDowDateField(test.fieldNexter, test.modifiers)
		if (err != nil || test.err != "") && err.Error() != test.err {
			t.Errorf("parseDowDateField(%v, %q) error = %v WANT %v", test.fieldNexter, test.modifiers, err, test.err)
		}
		if err != nil {
			continue
		}
		if dfn.fieldNexter != test.fieldNexter || dfn.isLast != test.resultLast || dfn.number != test.resultHash {
			t.Errorf("parseDowDateField(%v, %q) result = %v WANT %v, %v, %v",
				test.fieldNexter, test.modifiers, dfn, test.fieldNexter, test.resultLast, test.resultHash,
			)
		}
	}
}

func TestHasAndRemoveModifier(t *testing.T) {
	tests := []struct {
		value    string
		modifier string
		result   string
		has      bool
	}{
		{"", Last, "", false},
		{Last, "", Last, false},
		{Last + Asterisk, Last, Asterisk, true},
		{Asterisk + Last, Last, Asterisk, true},
		{Hash + "value", Hash, "value", true},
		{Hash + "value", Last, Hash + "value", false},
	}
	for _, test := range tests {
		result, has := hasAndRemoveModifier(test.value, test.modifier)
		if result != test.result || has != test.has {
			t.Errorf("hasAndRemoveModifier(%q, %q) = %q, %v WANT %q, %v", test.value, test.modifier, result, has, test.result, test.has)
		}
	}
}

func TestParseFieldNexterPart_errors(t *testing.T) {
	tests := []struct {
		value string
		fi    fieldIndex
	}{
		{"", second},
		{"-", second},
		{"JAN   ", month},
		{"JAN", second},
		{"/2", second},
		{"3/2", second},
		{"14-16/a", second},
		{"14-16/-1", second},
		{"15-34/", second},
		{"23-12", second},
		{Last, second},
		{Weekday, second},
		{Last + Weekday, second},
		{"2-34#5", second},
	}
	for _, test := range tests {
		result, err := parseFieldNexterPart(test.value, test.fi)
		if err == nil {
			t.Errorf("parseFieldNexterPart(%v, %v) = %v, %v WANT nil, non-nil", test.value, test.fi, result, err)
		}
	}
}

func TestParseFieldNexterPart_valueNexter(t *testing.T) {
	value := "2026"
	fi := year
	result, err := parseFieldNexterPart(value, fi)
	want := valueNexter(2026)
	if err != nil || result != want {
		t.Errorf("parseFieldNexterPart(%v, %v) = %v, %v WANT %v, %v", value, fi, result, err, want, nil)
	}
}

func TestParseFieldNexterPart_rangeNexter(t *testing.T) {
	value := "2-8"
	fi := month
	result, err := parseFieldNexterPart(value, fi)
	want := rangeNexter{2, 8}
	if err != nil || *(result.(*rangeNexter)) != want {
		t.Errorf("parseFieldNexterPart(%v, %v) = %v, %v WANT %v, %v", value, fi, result, err, want, nil)
	}
}

func TestParseFieldNexterPart_rangeDivNexter(t *testing.T) {
	value := "40-50/2"
	fi := minute
	result, err := parseFieldNexterPart(value, fi)
	want := newRangeDivNexter(newRangeNexter(40, 50), 2)
	if err != nil || !reflect.DeepEqual(result, want) {
		t.Errorf("parseFieldNexterPart(%v, %v) = %v, %v WANT %v, %v", value, fi, result, err, want, nil)
	}
}

func TestParseRangeOrConstantNexter_range(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  string
		result *rangeNexter
		err    string
	}{
		{second, Hyphen, nil, "left side of range " + errParseInteger.Error()},
		{second, "0-1", &rangeNexter{0, 1}, ""},
	}
	for _, test := range tests {
		result, err := parseRangeOrConstantNexter(test.value, test.fi)
		if (err != nil || test.err != "") && err.Error() != test.err {
			t.Errorf("parseRangeOrConstantNexter(%v, %v) error = %v WANT %v", test.value, test.fi, result, test.err)
		}
		if err != nil {
			continue
		}
		rn := result.(*rangeNexter)
		if rn.min != test.result.min || rn.max != test.result.max {
			t.Errorf("parseRangeOrConstantNexter(%v, %v) result = %v, %v WANT %v, %v",
				test.value, test.fi, rn.min, rn.max, test.result.min, test.result.max,
			)
		}
	}
}

func TestParseRangeOrConstantNexter_notRange(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  string
		result valueNexter
		err    string
	}{
		{second, "", invalidValue, errParseInteger.Error()},
		{second, "57", 57, ""},
	}
	for _, test := range tests {
		result, err := parseRangeOrConstantNexter(test.value, test.fi)
		if (err != nil || test.err != "") && err.Error() != test.err {
			t.Errorf("parseRangeOrConstantNexter(%v, %v) error = %v WANT %v", test.value, test.fi, err, test.err)
		}
		if err != nil {
			continue
		}
		vn := result.(valueNexter)
		if vn != test.result {
			t.Errorf("parseRangeOrConstantNexter(%v, %v) result = %v WANT %v", test.value, test.fi, vn, test.result)
		}
	}
}

func TestParseRangeNexter(t *testing.T) {
	tests := []struct {
		fi    fieldIndex
		value string
		min   int
		max   int
		err   string
	}{
		{second, "", invalidValue, invalidValue, errNoHyphen.Error()},
		{second, "a", invalidValue, invalidValue, errNoHyphen.Error()},
		{second, "-25", invalidValue, invalidValue, "left side of range " + errParseInteger.Error()},
		{second, "25-", invalidValue, invalidValue, "right side of range " + errParseInteger.Error()},
		{month, "a-b", invalidValue, invalidValue, "left side of range " + errParseIntegerAlias.Error()},
		{month, "-b", invalidValue, invalidValue, "left side of range " + errParseIntegerAlias.Error()},
		{month, "4-2", invalidValue, invalidValue, "left side value of range must be strictly less than right side value"},
		{dow, "SUN-SUN", invalidValue, invalidValue, "left side value of range must be strictly less than right side value"},
		{month, "feb-may", int(time.February), int(time.May), ""},
		{year, "1-3005", 1, 3005, ""},
	}
	for _, test := range tests {
		result, err := parseRangeNexter(test.value, test.fi)
		if err == nil && test.err != "" {
			t.Fatalf("parseRangeNexter(%v, %v) WANT ERROR got nil", test.value, test.fi)
		}
		if err != nil && err.Error() != test.err {
			t.Errorf("parseRangeNexter(%v, %v) error = %v WANT %v", test.value, test.fi, err, test.err)
		}
		if err == nil && (result == nil || result.min != test.min || result.max != test.max) {
			t.Errorf("parseRangeNexter(%v, %v) result = %v, %v WANT %v, %v", test.value, test.fi, result.min, result.max, test.min, test.max)
		}
	}
}

func TestConvertPossibleAnyToRange(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  string
		result string
	}{
		{second, "14", "14"},
		{second, Asterisk, second.rangeString()},
		{dom, Asterisk, dom.rangeString()},
		{dom, Question, dom.rangeString()},
		{dow, Asterisk, dow.rangeString()},
		{dow, Question, dow.rangeString()},
		{year, Question, Question},
	}
	for _, test := range tests {
		if result := convertPossibleAnyToRange(test.value, test.fi); result != test.result {
			t.Errorf("convertPossibleAnyToRange(%v, %v) = %v WANT %v", test.value, test.fi, result, test.result)
		}
	}
}

func TestParseValueNexter(t *testing.T) {
	tests := []struct {
		fi       fieldIndex
		value    string
		result   int
		hasError bool
	}{
		{second, "16", 16, false},
		{minute, "60", invalidValue, true},
		{hour, "0", 0, false},
		{month, "something", invalidValue, true},
		{dow, "tuesday", int(time.Tuesday), false},
		{year, "something", invalidValue, true},
	}
	for _, test := range tests {
		result, err := parseValueNexter(test.value, test.fi)
		if (err != nil) != test.hasError {
			t.Errorf("parseValueNexter(%v, %v) error = %v WANT ERROR %v", test.value, test.fi, err, test.hasError)
		}
		if int(result) != test.result {
			t.Errorf("parseValueNexter(%v, %v) result = %v WANT %v", test.value, test.fi, int(result), test.result)
		}
	}
}

func TestParseSingleValue(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  string
		result int
		err    string
	}{
		{second, "16", 16, ""},
		{minute, "60", invalidValue, errNotInRange.Error()},
		{hour, "0", 0, ""},
		{month, "something", invalidValue, "must be a decimal integer or valid string alias"},
		{dow, "tuesday", int(time.Tuesday), ""},
		{dow, "SUN", int(time.Sunday), ""},
		{year, "something", invalidValue, "must be a decimal integer"},
	}
	for _, test := range tests {
		result, err := parseSingleValue(test.value, test.fi)
		if err != nil && err.Error() != test.err {
			t.Errorf("parseSingleValue(%v, %v) error = %v WANT %v", test.value, test.fi, err, test.err)
		}
		if result != test.result {
			t.Errorf("parseSingleValue(%v, %v) result = %v WANT %v", test.value, test.fi, result, test.result)
		}
	}
}

func TestConvertPossibleMonthDowToInteger(t *testing.T) {
	tests := []struct {
		fi     fieldIndex
		value  string
		result string
	}{
		{second, "something", "something"},
		{month, "february", fmt.Sprint(int(time.February))},
		{dow, "sunday", fmt.Sprint(int(time.Sunday))},
		{month, "something", "something"},
	}
	for _, test := range tests {
		if result := convertPossibleMonthDowToInteger(test.value, test.fi); result != test.result {
			t.Errorf("convertPossibleMonthDowToInteger(%v) = %v WANT %v", test.value, result, test.result)
		}
	}
}

func TestConvertMonthToInteger(t *testing.T) {
	tests := []struct {
		value  string
		result string
	}{
		{"", ""},
		{"a", "a"},
		{"JA", "JA"},
		{"jan", fmt.Sprint(int(time.January))},
		{"MAY", fmt.Sprint(int(time.May))},
		{"Octo", fmt.Sprint(int(time.October))},
		{"july and some more stuff", "july and some more stuff"},
		{"4", "4"},
		{"20", "20"},
	}
	for _, test := range tests {
		if result := convertMonthToInteger(test.value); result != test.result {
			t.Errorf("convertMonthToInteger(%v) = %v WANT %v", test.value, result, test.result)
		}
	}
}

func TestConvertDowToInteger(t *testing.T) {
	tests := []struct {
		value  string
		result string
	}{
		{"", ""},
		{"a", "a"},
		{"JA", "JA"},
		{"sun", fmt.Sprint(int(time.Sunday))},
		{"Thurs", fmt.Sprint(int(time.Thursday))},
		{"FRIDAY", fmt.Sprint(int(time.Friday))},
		{"monday and some more stuff", "monday and some more stuff"},
		{"4", "4"},
		{"8", "8"},
	}
	for _, test := range tests {
		if result := convertDowToInteger(test.value); result != test.result {
			t.Errorf("convertDowToInteger(%v) = %v WANT %v", test.value, result, test.result)
		}
	}
}

func TestParseIncValue(t *testing.T) {
	const errString = "step value must be a positive decimal integer"
	tests := []struct {
		value  string
		result int
		err    string
	}{
		{"", invalidValue, "step value " + errEmpty.Error()},
		{"a", invalidValue, errString},
		{"-1", invalidValue, errString},
		{"0", invalidValue, errString},
		{"0xA", invalidValue, errString},
		{"12", 12, ""},
	}
	for _, test := range tests {
		result, err := parseIncValue(test.value)
		if (test.err != "" || err != nil) && err.Error() != test.err {
			t.Errorf("parseIncValue(%v) error = %v WANT %v", test.value, err, test.err)
		}
		if result != test.result {
			t.Errorf("parseIncValue(%v) = %v WANT %v", test.value, result, test.result)
		}
	}
}
