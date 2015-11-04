package sched

import (
	"reflect"
	"testing"
)

func TestNewParseError(t *testing.T) {
	err := NewParseError("expression", "description")
	if err.Expression != "expression" || err.Description != "description" {
		t.Fail()
	}
}

func TestParseError_Error(t *testing.T) {
	err := NewParseError("expression", "description")
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
		if result := FieldsFunc(test.value); result != test.result {
			t.Errorf("FieldsFunc(%v) = %v WANT %v", test.value, result, test.result)
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
