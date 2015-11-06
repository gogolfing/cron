package sched

import (
	"reflect"
	"strings"
	"testing"
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
		{"a b ", nil, true},
		{"a b c", nil, true},
		{"a b c d", nil, true},
		{"a b c d e", []string{Asterisk, "a", "b", "c", "d", "e", Asterisk}, false},
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
	}
}

func TestValidateNumberOfFields(t *testing.T) {
	errString := "number of fields must be 1, 5, 6, or 7"
	tests := []struct {
		fields []string
		count  int
		err    string
	}{
		{nil, 0, errString},
		{[]string{}, 0, errString},
		{make([]string, 1), 1, ""},
		{make([]string, 2), 0, errString},
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
