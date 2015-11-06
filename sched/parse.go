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

func newParseError(exp, desc string) *ParseError {
	return &ParseError{
		Expression:  exp,
		Description: desc,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("sched: could not parse %q: %v", p.Expression, p.Description)
}

var fieldSeparatorFunc = func(r rune) bool {
	return strings.ContainsRune(FieldSeparators, r)
}

func Fields(expression string) []string {
	return strings.FieldsFunc(strings.Trim(expression, TrimCutset), fieldSeparatorFunc)
}

func FieldParts(field string) []string {
	result := strings.Split(field, Comma)
	if len(result) == 1 && len(result[0]) == 0 {
		return []string{}
	}
	return result
}

func MustParse(expression string) Schedule {
	s, err := Parse(expression)
	if err != nil {
		panic(err)
	}
	return s
}

func Parse(expression string) (Schedule, error) {
	//ParseErrors should be returned from this function and no others.
	fields, err := getNormalizedFields(expression)
	if err != nil {
		return nil, newParseError(expression, err.Error())
	}
	fmt.Println(fields, err)
	return nil, nil
}

func getNormalizedFields(expression string) ([]string, error) {
	fields := Fields(expression)
	var err error = nil
	count, err := validateNumberOfFields(fields)
	if err != nil {
		return nil, err
	}
	if count == 1 {
		fields, err = getNormalizedDirectiveFields(fields[0])
		if err != nil {
			return nil, err
		}
	}
	if count == 5 {
		fields = append([]string{Asterisk}, fields...)
		count++
	}
	if count == 6 {
		fields = append(fields, Asterisk)
		count++
	}
	return fields, nil
}

func validateNumberOfFields(fields []string) (int, error) {
	count := len(fields)
	if count != 1 && count != 5 && count != 6 && count != 7 {
		return 0, fmt.Errorf("number of fields must be 1, 5, 6, or 7")
	}
	return count, nil
}

func getNormalizedDirectiveFields(directive string) ([]string, error) {
	format := ""
	directive = strings.ToLower(directive)
	switch directive {
	case Yearly:
		format = YearlyFormat
	case Annually:
		format = AnnuallyFormat
	case Monthly:
		format = MonthlyFormat
	case Weekly:
		format = WeeklyFormat
	case Daily:
		format = DailyFormat
	case Hourly:
		format = HourlyFormat
	case Minutely:
		format = MinutelyFormat
	case Secondly:
		format = SecondlyFormat
	}
	if format == "" {
		return nil, fmt.Errorf("the directive %q is not recognized", directive)
	}
	return Fields(format), nil
}
