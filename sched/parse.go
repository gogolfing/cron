package sched

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	fieldStrings, err := getNormalizedFields(expression)
	if err != nil {
		return nil, newParseError(expression, err.Error())
	}
	s := newSchedule()
	for i, fieldString := range fieldStrings {
		fi := fieldIndex(i)
		nexter, err := parseField(fieldString, fi)
		if err != nil {
			return nil, newParseError(expression, err.Error())
		}
		s.setNexter(nexter, fi)
	}
	return s, nil
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

func parseField(field string, fi fieldIndex) (nexter interface{}, err error) {
	parts := FieldParts(field)
	if fi.isDateField() {
		nexter, err = parseDateFieldNexterParts(parts, fi)
	} else {
		nexter, err = parseFieldNexterParts(parts, fi)
	}
	if err != nil {
		err = fmt.Errorf("%v field: %v", fi, err.Error())
		nexter = nil
	}
	return
}

func parseDateFieldNexterParts(parts []string, fi fieldIndex) (dateFieldNexter, error) {
	if len(parts) == 1 {
		return parseDateFieldNexterPart(parts[0], fi)
	}
	result := multiDateFieldNexter(make([]dateFieldNexter, 0, len(parts)))
	for i, part := range parts {
		nexter, err := parseDateFieldNexterPart(part, fi)
		if err != nil {
			return nil, newPartError(i, err)
		}
		result = append(result, nexter)
	}
	return result, nil
}

func parseFieldNexterParts(parts []string, fi fieldIndex) (fieldNexter, error) {
	if len(parts) == 1 {
		return parseFieldNexterPart(parts[0], fi)
	}
	result := multiNexter(make([]fieldNexter, 0, len(parts)))
	for i, part := range parts {
		nexter, err := parseFieldNexterPart(part, fi)
		if err != nil {
			return nil, newPartError(i, err)
		}
		result = append(result, nexter)
	}
	return result, nil
}

func newPartError(index int, old error) error {
	return fmt.Errorf("part %v: %v", index+1, old.Error())
}

func parseFieldNexterPart(part string, fi fieldIndex) (fieldNexter, error) {
	if len(part) == 0 {
		return nil, fmt.Errorf("cannot be empty")
	}
	slashIndex := strings.Index(part, Slash)
	if slashIndex < 0 {
		return parseRangeOrConstantNexter(part[:slashIndex], fi)
	}
	if slashIndex == 0 {
		return nil, fmt.Errorf("value before step cannot be empty")
	}
	rn, err := parseRangeNexter(part[:slashIndex], fi)
	if err != nil {
		return nil, err
	}
	inc, err := parseIncValue(part[slashIndex+1:])
	if err != nil {
		return nil, err
	}
	return newRangeDivNexter(rn, inc), nil
}

func parseDateFieldNexterPart(part string, fi fieldIndex) (dateFieldNexter, error) {
	if len(part) == 0 {
		return nil, fmt.Errorf("cannot be empty")
	}
	return nil, nil
}

func parseRangeOrConstantNexter(part string, fi fieldIndex) (fieldNexter, error) {
	part = convertPossibleAnyToRange(part, fi)
	if strings.Contains(part, Hyphen) {
		return parseRangeNexter(part, fi)
	}
	return parseValueNexter(part, fi)
}

func convertPossibleAnyToRange(part string, fi fieldIndex) string {
	if fi.isDateField() {
		part = strings.Replace(part, Question, Asterisk, -1)
	}
	return strings.Replace(part, Asterisk, fi.rangeString(), -1)
}

//parseRangeNexter will panic if part does not conatin Hyphen.
func parseRangeNexter(part string, fi fieldIndex) (*rangeNexter, error) {
	hyphenIndex := strings.Index(part, Hyphen)
	min, err := parseSingleValue(part[:hyphenIndex], fi)
	if err != nil {
		return nil, fmt.Errorf("left side of range: %v", err.Error())
	}
	max, err := parseSingleValue(part[hyphenIndex+1:], fi)
	if err != nil {
		return nil, fmt.Errorf("right side of range: %v", err.Error())
	}
	if min >= max {
		return nil, fmt.Errorf("left side of range must be strictly less than right side")
	}
	return newRangeNexter(min, max), nil
}

func parseValueNexter(part string, fi fieldIndex) (valueNexter, error) {
	value, err := parseSingleValue(part, fi)
	if err != nil {
		return valueNexter(invalidValue), err
	}
	return valueNexter(value), nil
}

func parseSingleValue(value string, fi fieldIndex) (int, error) {
	value = convertPossibleMonthDowToInteger(value, fi)
	result, err := strconv.Atoi(value)
	if err != nil {
		if fi.isDateField() {
			return invalidValue, fmt.Errorf("must be a decimal integer or valid string alias")
		}
		return invalidValue, fmt.Errorf("must be a decimal integer")
	}
	if !fi.isInRange(result) {
		return invalidValue, fmt.Errorf("not in range")
	}
	return result, nil
}

func convertPossibleMonthDowToInteger(value string, fi fieldIndex) string {
	if fi == month {
		return convertMonthToInteger(value)
	}
	if fi == dow {
		return converDowToInteger(value)
	}
	return value
}

func convertMonthToInteger(value string) string {
	if len(value) < 3 {
		return value
	}
	value = strings.ToUpper(value)
	for m := time.January; m <= time.December; m++ {
		if strings.HasPrefix(strings.ToUpper(fmt.Sprint(m)), value) {
			return fmt.Sprint(int(m))
		}
	}
	return value
}

func converDowToInteger(value string) string {
	if len(value) < 3 {
		return value
	}
	value = strings.ToUpper(value)
	for w := time.Sunday; w <= time.Saturday; w++ {
		if strings.HasPrefix(strings.ToUpper(fmt.Sprint(w)), value) {
			return fmt.Sprint(int(w))
		}
	}
	return value
}

func parseIncValue(value string) (int, error) {
	if len(value) == 0 {
		return invalidValue, fmt.Errorf("step value cannot be empty")
	}
	inc, err := strconv.Atoi(value)
	if err != nil || inc <= 0 {
		return invalidValue, fmt.Errorf("step value must be a positive decimal integer")
	}
	return inc, nil
}
