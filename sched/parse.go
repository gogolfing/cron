package sched

import (
	"fmt"
	"strings"
)

const (
	Asterisk = "*"
	Question = "?"
	Hyphen   = "-"
	Slash    = "/"
	Comma    = ","
	Hash     = "#"
	Last     = "L"
	Weekday  = "W"

	FieldSeparators = " \t"
	TrimCutset      = FieldSeparators + "\n"
)

type ParseError struct {
	Expression  string
	Description string
}

func NewParseError(exp, desc string) *ParseError {
	return &ParseError{
		Expression:  exp,
		Description: desc,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("sched: could not parse %q: %v", p.Expression, p.Description)
}

var FieldsFunc = func(r rune) bool {
	return strings.ContainsRune(FieldSeparators, r)
}

func Fields(exp string) []string {
	return strings.FieldsFunc(strings.Trim(exp, TrimCutset), FieldsFunc)
}

func FieldParts(field string) []string {
	result := strings.Split(field, Comma)
	if len(result) == 1 && len(result[0]) == 0 {
		return []string{}
	}
	return result
}

func MustParse(value string) Schedule {
	s, err := Parse(value)
	if err != nil {
		panic(err)
	}
	return s
}

func Parse(value string) (Schedule, error) {
	return nil, nil
}
