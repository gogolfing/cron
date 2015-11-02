package sched

import "strings"

const (
	AnyString    = "*"
	AnyDayString = "?"

	RangeSeparator = "-"
	RangeDividor   = "/"

	FieldSeparator     = " \t"
	FieldPartSeparator = ","

	TrimCutset = FieldSeparator + "\n"
)

var FieldsFunc = func(r rune) bool {
	return strings.ContainsRune(FieldSeparator, r)
}

func Fields(value string) []string {
	return strings.FieldsFunc(strings.Trim(value, TrimCutset), FieldsFunc)
}

func FieldParts(field string) []string {
	return strings.Split(field, FieldPartSeparator)
}
